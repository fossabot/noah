/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 */

package sql

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgerror"
	nodes "github.com/Ready-Stock/pg_query_go/nodes"
)

func (ex *connExecutor) execPrepare(parseCmd PrepareStmt) error {
	if parseCmd.Name != "" {
		if _, ok := ex.prepStmtsNamespace.prepStmts[parseCmd.Name]; ok {
			err := pgerror.NewErrorf(
				pgerror.CodeDuplicatePreparedStatementError,
				"prepared statement %q already exists", parseCmd.Name,
			)
			return err
		}
	} else {
		// Deallocate the unnamed statement, if it exists.
		ex.deletePreparedStmt("")
	}

	_, err := ex.addPreparedStmt(parseCmd.Name, parseCmd.PGQuery)
	if err != nil {
		return err
	}
	return nil
}

func (ex *connExecutor) execBind(bindCmd BindStmt) error {
	portalName := bindCmd.PortalName
	// The unnamed portal can be freely overwritten.
	if portalName != "" {
		if _, ok := ex.prepStmtsNamespace.portals[portalName]; ok {
			return pgerror.NewErrorf(pgerror.CodeDuplicateCursorError, "portal %q already exists", portalName)
		}
	} else {
		// Deallocate the unnamed portal, if it exists.
		ex.deletePortal("")
	}

	ps, ok := ex.prepStmtsNamespace.prepStmts[bindCmd.PreparedStatementName]
	if !ok {
		return pgerror.NewErrorf(pgerror.CodeInvalidSQLStatementNameError, "unknown prepared statement %q", bindCmd.PreparedStatementName)
	}

	// Create the new PreparedPortal.
	if err := ex.addPortal(portalName, bindCmd.PreparedStatementName, ps.PreparedStatement); err != nil {
		return err
	}
	return nil
}

func (ex *connExecutor) addPortal(portalName string, psName string, stmt *PreparedStatement) error {
	if _, ok := ex.prepStmtsNamespace.portals[portalName]; ok {
		panic(fmt.Sprintf("portal already exists: %q", portalName))
	}
	portal := ex.newPreparedPortal(stmt)
	ex.prepStmtsNamespace.portals[portalName] = portalEntry{
		PreparedPortal: &portal,
		psName:         psName,
	}
	ex.prepStmtsNamespace.prepStmts[psName].portals[portalName] = struct{}{}
	return nil
}

func (ex *connExecutor) addPreparedStmt(name string, stmt nodes.Stmt) (*PreparedStatement, error) {
	if _, ok := ex.prepStmtsNamespace.prepStmts[name]; ok {
		panic(fmt.Sprintf("prepared statement already exists: %q", name))
	}
	// Prepare the query. This completes the typing of placeholders.
	prepared, err := ex.prepare(stmt)
	if err != nil {
		return nil, err
	}
	ex.prepStmtsNamespace.prepStmts[name] = prepStmtEntry{
		PreparedStatement: prepared,
		portals:           make(map[string]struct{}),
	}
	return prepared, nil
}

func (ex *connExecutor) prepare(stmt nodes.Stmt) (*PreparedStatement, error) {
	prepared := &PreparedStatement{
		Statement: &stmt,
	}
	return prepared, nil
}

func (ex *connExecutor) deletePreparedStmt(name string) {
	psEntry, ok := ex.prepStmtsNamespace.prepStmts[name]
	if !ok {
		return
	}
	for portalName := range psEntry.portals {
		ex.deletePortal(portalName)
	}
	delete(ex.prepStmtsNamespace.prepStmts, name)
}

func (ex *connExecutor) deletePortal(name string) {
	portalEntry, ok := ex.prepStmtsNamespace.portals[name]
	if !ok {
		return
	}
	delete(ex.prepStmtsNamespace.portals, name)
	delete(ex.prepStmtsNamespace.prepStmts[portalEntry.psName].portals, name)
}
