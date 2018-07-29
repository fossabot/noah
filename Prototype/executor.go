package Prototype

import (
	//"github.com/Ready-Stock/pg_query_go/nodes"
	pgq "github.com/Ready-Stock/pg_query_go"
	"errors"
	query "github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/Ready-Stock/Noah/Prototype/queries"
	"fmt"
)

func InjestQuery(query string) error {
	if parsed, err := pgq.Parse(query); err != nil {
		return err
	} else {
		fmt.Printf("Injesting Query: %s \n", query)
		return handleParseTree(*parsed)
	}
}

func handleParseTree(tree pgq.ParsetreeList) error {
	if len(tree.Statements) == 0 {
		return errors.New("no statements provided")
	} else {
		raw := tree.Statements[0].(query.RawStmt)
		switch stmt := raw.Stmt.(type) {
		case query.AlterCollationStmt:
		case query.AlterDatabaseSetStmt:
		case query.AlterDatabaseStmt:
		case query.AlterDefaultPrivilegesStmt:
		case query.AlterDomainStmt:
		case query.AlterEnumStmt:
		case query.AlterEventTrigStmt:
		case query.AlterExtensionContentsStmt:
		case query.AlterExtensionStmt:
		case query.AlterFdwStmt:
		case query.AlterForeignServerStmt:
		case query.AlterFunctionStmt:
		case query.AlterObjectDependsStmt:
		case query.AlterObjectSchemaStmt:
		case query.AlterOperatorStmt:
		case query.AlterOpFamilyStmt:
		case query.AlterOwnerStmt:
		case query.AlterPolicyStmt:
		case query.AlterPublicationStmt:
		case query.AlterRoleSetStmt:
		case query.AlterRoleStmt:
		case query.AlterSeqStmt:
		case query.AlterSubscriptionStmt:
		case query.AlterSystemStmt:
		case query.AlterTableMoveAllStmt:
		case query.AlterTableSpaceOptionsStmt:
		case query.AlterTableStmt:
		case query.AlterTSConfigurationStmt:
		case query.AlterTSDictionaryStmt:
		case query.AlterUserMappingStmt:
		case query.CheckPointStmt:
		case query.ClosePortalStmt:
		case query.ClusterStmt:
		case query.CommentStmt:
		case query.CompositeTypeStmt:
		case query.ConstraintsSetStmt:
		case query.CopyStmt:
		case query.CreateAmStmt:
		case query.CreateCastStmt:
		case query.CreateConversionStmt:
		case query.CreateDomainStmt:
		case query.CreateEnumStmt:
		case query.CreateEventTrigStmt:
		case query.CreateExtensionStmt:
		case query.CreateFdwStmt:
		case query.CreateForeignServerStmt:
		case query.CreateForeignTableStmt:
		case query.CreateFunctionStmt:
		case query.CreatePLangStmt:
		case query.CreatePolicyStmt:
		case query.CreatePublicationStmt:
		case query.CreateRangeStmt:
		case query.CreateRoleStmt:
		case query.CreateSchemaStmt:
		case query.CreateSeqStmt:
		case query.CreateStatsStmt:
		case query.CreateStmt:
		case query.CreateSubscriptionStmt:
		case query.CreateTableAsStmt:
		case query.CreateTableSpaceStmt:
		case query.CreateTransformStmt:
		case query.CreateTrigStmt:
		case query.CreateUserMappingStmt:
		case query.CreatedbStmt:
		case query.DeallocateStmt:
		case query.DeclareCursorStmt:
		case query.DefineStmt:
		case query.DeleteStmt:
		case query.DiscardStmt:
		case query.DoStmt:
		case query.DropOwnedStmt:
		case query.DropRoleStmt:
		case query.DropStmt:
		case query.DropSubscriptionStmt:
		case query.DropTableSpaceStmt:
		case query.DropUserMappingStmt:
		case query.DropdbStmt:
		case query.ExecuteStmt:
		case query.ExplainStmt:
		case query.FetchStmt:
		case query.GrantRoleStmt:
		case query.ImportForeignSchemaStmt:
		case query.IndexStmt:
		case query.InsertStmt:
			return queries.HandleInsert(stmt)
		case query.ListenStmt:
		case query.LoadStmt:
		case query.LockStmt:
		case query.NotifyStmt:
		case query.PrepareStmt:
		case query.ReassignOwnedStmt:
		case query.RefreshMatViewStmt:
		case query.ReindexStmt:
		case query.RenameStmt:
		case query.ReplicaIdentityStmt:
		case query.RuleStmt:
		case query.SecLabelStmt:
		case query.SelectStmt:
			return queries.HandleSelect(stmt)
		case query.SetOperationStmt:
		case query.TransactionStmt:
		case query.TruncateStmt:
		case query.UnlistenStmt:
		case query.UpdateStmt:
		case query.VacuumStmt:
		case query.VariableSetStmt:
		case query.VariableShowStmt:
		case query.ViewStmt:
		default:
			return errors.New("invalid or unsupported query type")
		}
		return nil
	}
	return nil
}
