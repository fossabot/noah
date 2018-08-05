package sql

import (
	"github.com/lib/pq/oid"
	"strconv"
	"fmt"
	"github.com/Ready-Stock/Noah/Database/sql/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/Database/sql/pgwire/pgwirebase"
	"github.com/Ready-Stock/Noah/Database/util/fsm"
	"github.com/Ready-Stock/Noah/Database/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/Ready-Stock/pg_query_go"
)

func (ex *connExecutor) execPrepare(parseCmd PrepareStmt) (fsm.Event, fsm.EventPayload) {

	retErr := func(err error) (fsm.Event, fsm.EventPayload) {
		return eventNonRetriableErr{IsCommit: fsm.False}, eventNonRetriableErrPayload{err: err}
	}

	// The anonymous statement can be overwritter.
	if parseCmd.Name != "" {
		if _, ok := ex.prepStmtsNamespace.prepStmts[parseCmd.Name]; ok {
			err := pgerror.NewErrorf(
				pgerror.CodeDuplicatePreparedStatementError,
				"prepared statement %q already exists", parseCmd.Name,
			)
			return retErr(err)
		}
	} else {
		// Deallocate the unnamed statement, if it exists.
		ex.deletePreparedStmt("")
	}

	ps, err := ex.addPreparedStmt(parseCmd.Name, parseCmd.PGQuery, parseCmd.TypeHints)
	if err != nil {
		return retErr(err)
	}

	// Convert the inferred SQL types back to an array of pgwire Oids.
	inTypes := make([]oid.Oid, 0, len(ps.TypeHints))
	if len(ps.TypeHints) > pgwirebase.MaxPreparedStatementArgs {
		return retErr(
			pgwirebase.NewProtocolViolationErrorf(
				"more than %d arguments to prepared statement: %d",
				pgwirebase.MaxPreparedStatementArgs, len(ps.TypeHints)))
	}
	for k, t := range ps.TypeHints {
		i, err := strconv.Atoi(k)
		if err != nil || i < 1 {
			return retErr(pgerror.NewErrorf(
				pgerror.CodeUndefinedParameterError, "invalid placeholder name: $%s", k))
		}
		// Placeholder names are 1-indexed; the arrays in the protocol are
		// 0-indexed.
		i--
		// Grow inTypes to be at least as large as i. Prepopulate all
		// slots with the hints provided, if any.
		for j := len(inTypes); j <= i; j++ {
			inTypes = append(inTypes, 0)
			if j < len(parseCmd.RawTypeHints) {
				inTypes[j] = parseCmd.RawTypeHints[j]
			}
		}
		// OID to Datum is not a 1-1 mapping (for example, int4 and int8
		// both map to TypeInt), so we need to maintain the types sent by
		// the client.
		if inTypes[i] != 0 {
			continue
		}
		inTypes[i] = t.Oid()
	}
	for i, t := range inTypes {
		if t == 0 {
			return retErr(pgerror.NewErrorf(
				pgerror.CodeIndeterminateDatatypeError,
				"could not determine data type of placeholder $%d", i+1))
		}
	}
	// Remember the inferred placeholder types so they can be reported on
	// Describe.
	ps.InTypes = inTypes
	return nil, nil
}

func (ex *connExecutor) addPreparedStmt(name string, stmt pg_query.ParsetreeList, placeholderHints tree.PlaceholderTypes) (*PreparedStatement, error) {
	if _, ok := ex.prepStmtsNamespace.prepStmts[name]; ok {
		panic(fmt.Sprintf("prepared statement already exists: %q", name))
	}

	// Prepare the query. This completes the typing of placeholders.
	prepared, err := ex.prepare(stmt, placeholderHints)
	if err != nil {
		return nil, err
	}

	ex.prepStmtsNamespace.prepStmts[name] = prepStmtEntry{
		PreparedStatement: prepared,
		portals:           make(map[string]struct{}),
	}
	return prepared, nil
}

func (ex *connExecutor) prepare(stmt pg_query.ParsetreeList, placeholderHints tree.PlaceholderTypes) (*PreparedStatement, error) {
	if placeholderHints == nil {
		placeholderHints = make(tree.PlaceholderTypes)
	}

	prepared := &PreparedStatement{
		TypeHints: placeholderHints,
	}

	if stmt.AST == nil {
		return prepared, nil
	}
	prepared.Str = stmt.String()

	prepared.Statement = stmt.AST
	prepared.AnonymizedStr = anonymizeStmt(stmt)

	if err := placeholderHints.ProcessPlaceholderAnnotations(stmt.AST); err != nil {
		return nil, err
	}
	// Preparing needs a transaction because it needs to retrieve db/table
	// descriptors for type checking.
	// TODO(andrei): Needing a transaction for preparing seems non-sensical, as
	// the prepared statement outlives the txn. I hope that it's not used for
	// anything other than getting a timestamp.
	txn := client.NewTxn(ex.server.cfg.DB, ex.server.cfg.NodeID.Get(), client.RootTxn)

	// Create a plan for the statement to figure out the typing, then close the
	// plan.
	if err := func() error {
		p := &ex.planner
		ex.resetPlanner(ctx, p, txn, ex.server.cfg.Clock.PhysicalTime() /* stmtTimestamp */)
		p.semaCtx.Placeholders.SetTypeHints(placeholderHints)
		p.extendedEvalCtx.PrepareOnly = true
		p.extendedEvalCtx.ActiveMemAcc = &prepared.memAcc
		// constantMemAcc accounts for all constant folded values that are computed
		// prior to any rows being computed.
		constantMemAcc := p.extendedEvalCtx.Mon.MakeBoundAccount()
		p.extendedEvalCtx.ActiveMemAcc = &constantMemAcc
		defer constantMemAcc.Close(ctx)

		protoTS, err := p.isAsOf(stmt.AST, ex.server.cfg.Clock.Now() /* max */)
		if err != nil {
			return err
		}
		if protoTS != nil {
			p.asOfSystemTime = true
			// We can't use cached descriptors anywhere in this query, because
			// we want the descriptors at the timestamp given, not the latest
			// known to the cache.
			p.avoidCachedDescriptors = true
			txn.SetFixedTimestamp(ctx, *protoTS)
		}

		// PREPARE has a limited subset of statements it can be run with. Postgres
		// only allows SELECT, INSERT, UPDATE, DELETE and VALUES statements to be
		// prepared.
		// See: https://www.postgresql.org/docs/current/static/sql-prepare.html
		// However, we allow a large number of additional statements.
		// As of right now, the optimizer only works on SELECT statements and will
		// fallback for all others, so this should be safe for the foreseeable
		// future.
		if optimizerPlanned, err := p.optionallyUseOptimizer(ctx, ex.sessionData, stmt); err != nil {
			return err
		} else if !optimizerPlanned {
			isCorrelated := p.curPlan.isCorrelated
			log.VEventf(ctx, 1, "query is correlated: %v", isCorrelated)
			// Fallback if the optimizer was not enabled or used.
			if err := p.prepare(ctx, stmt.AST); err != nil {
				enhanceErrWithCorrelation(err, isCorrelated)
				return err
			}
		}

		if p.curPlan.plan == nil {
			// The statement cannot be prepared. Nothing to do.
			return nil
		}
		defer p.curPlan.close(ctx)

		prepared.Columns = p.curPlan.columns()
		for _, c := range prepared.Columns {
			if err := checkResultType(c.Typ); err != nil {
				return err
			}
		}
		prepared.Types = p.semaCtx.Placeholders.Types
		return nil
	}(); err != nil {
		txn.CleanupOnError(ctx, err)
		return nil, err
	}
	if err := txn.CommitOrCleanup(ctx); err != nil {
		return nil, err
	}

	return prepared, nil
}

func (ex *connExecutor) deletePreparedStmt(name string) {
	psEntry, ok := ex.prepStmtsNamespace.prepStmts[name]
	if !ok {
		return
	}
	// If the prepared statement only exists in prepStmtsNamespace, it's up to us
	// to close it.
	baseP, inBase := ex.extraTxnState.prepStmtsNamespaceAtTxnRewindPos.prepStmts[name]
	if !inBase || (baseP.PreparedStatement != psEntry.PreparedStatement) {
		psEntry.close(ctx)
	}
	for portalName := range psEntry.portals {
		ex.deletePortal(ctx, portalName)
	}
	delete(ex.prepStmtsNamespace.prepStmts, name)
}