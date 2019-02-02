/*
 * Copyright (c) 2019 Ready Stock
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
 */

package sql

import (
	"context"
	"errors"
	"github.com/readystock/noah/db/sql/plan"
	nodes "github.com/readystock/pg_query_go/nodes"
)

func (ex *connExecutor) execStmt(
	tree nodes.Stmt,
	res RestrictedCommandResult,
	pos CmdPos,
	pinfo *plan.PlaceholderInfo,
) error {
	switch tree.StatementType() {
	case nodes.DDL, nodes.RowsAffected:

	}

	handler, err := ex.getStatementHandler(tree)
	if err != nil {
		return err
	}

	return handler.Execute(context.Background(), ex, res, pinfo)
}

func (ex *connExecutor) getStatementHandler(tree nodes.Stmt) (IQueryStatement, error) {
	switch stmt := tree.(type) {
	// case nodes.AlterCollationStmt:
	// case nodes.AlterDatabaseSetStmt:
	// case nodes.AlterDatabaseStmt:
	// case nodes.AlterDefaultPrivilegesStmt:
	// case nodes.AlterDomainStmt:
	// case nodes.AlterEnumStmt:
	// case nodes.AlterEventTrigStmt:
	// case nodes.AlterExtensionContentsStmt:
	// case nodes.AlterExtensionStmt:
	// case nodes.AlterFdwStmt:
	// case nodes.AlterForeignServerStmt:
	// case nodes.AlterFunctionStmt:
	// case nodes.AlterObjectDependsStmt:
	// case nodes.AlterObjectSchemaStmt:
	// case nodes.AlterOperatorStmt:
	// case nodes.AlterOpFamilyStmt:
	// case nodes.AlterOwnerStmt:
	// case nodes.AlterPolicyStmt:
	// case nodes.AlterPublicationStmt:
	// case nodes.AlterRoleSetStmt:
	// case nodes.AlterRoleStmt:
	// case nodes.AlterSeqStmt:
	// case nodes.AlterSubscriptionStmt:
	// case nodes.AlterSystemStmt:
	// case nodes.AlterTableMoveAllStmt:
	// case nodes.AlterTableSpaceOptionsStmt:
	// case nodes.AlterTableStmt:
	// case nodes.AlterTSConfigurationStmt:
	// case nodes.AlterTSDictionaryStmt:
	// case nodes.AlterUserMappingStmt:
	// case nodes.CheckPointStmt:
	// case nodes.ClosePortalStmt:
	// case nodes.ClusterStmt:
	// case nodes.CommentStmt:
	//     // return nil, _comment.CreateCommentStatment(stmt, tree).HandleComment(ctx)
	// case nodes.CompositeTypeStmt:
	// case nodes.ConstraintsSetStmt:
	// case nodes.CopyStmt:
	// case nodes.CreateAmStmt:
	// case nodes.CreateCastStmt:
	// case nodes.CreateConversionStmt:
	// case nodes.CreateDomainStmt:
	// case nodes.CreateEnumStmt:
	// case nodes.CreateEventTrigStmt:
	// case nodes.CreateExtensionStmt:
	// case nodes.CreateFdwStmt:
	// case nodes.CreateForeignServerStmt:
	// case nodes.CreateForeignTableStmt:
	// case nodes.CreateFunctionStmt:
	// case nodes.CreatePLangStmt:
	// case nodes.CreatePolicyStmt:
	// case nodes.CreatePublicationStmt:
	// case nodes.CreateRangeStmt:
	// case nodes.CreateRoleStmt:
	// case nodes.CreateSchemaStmt:
	// case nodes.CreateSeqStmt:
	// case nodes.CreateStatsStmt:
	case nodes.CreateStmt:
		return CreateCreateStatement(stmt), nil
	// case nodes.CreateSubscriptionStmt:
	// case nodes.CreateTableAsStmt:
	// case nodes.CreateTableSpaceStmt:
	// case nodes.CreateTransformStmt:
	// case nodes.CreateTrigStmt:
	// case nodes.CreateUserMappingStmt:
	// case nodes.CreatedbStmt:
	// case nodes.DeallocateStmt:
	// case nodes.DeclareCursorStmt:
	// case nodes.DefineStmt:
	case nodes.DeleteStmt:
		return CreateDeleteStatement(stmt), nil
	// case nodes.DiscardStmt:
	// case nodes.DoStmt:
	// case nodes.DropOwnedStmt:
	// case nodes.DropRoleStmt:
	case nodes.DropStmt:
		return CreateDropStatement(stmt), nil
	// case nodes.DropSubscriptionStmt:
	// case nodes.DropTableSpaceStmt:
	// case nodes.DropUserMappingStmt:
	// case nodes.DropdbStmt:
	// case nodes.ExecuteStmt:
	// case nodes.ExplainStmt:
	// case nodes.FetchStmt:
	// case nodes.GrantRoleStmt:
	// case nodes.ImportForeignSchemaStmt:
	// case nodes.IndexStmt:
	case nodes.InsertStmt:
		return CreateInsertStatement(stmt), nil
	// case nodes.ListenStmt:
	// case nodes.LoadStmt:
	// case nodes.LockStmt:
	// case nodes.NotifyStmt:
	// case nodes.PrepareStmt:
	// case nodes.ReassignOwnedStmt:
	// case nodes.RefreshMatViewStmt:
	// case nodes.ReindexStmt:
	// case nodes.RenameStmt:
	// case nodes.ReplicaIdentityStmt:
	// case nodes.RuleStmt:
	// case nodes.SecLabelStmt:
	case nodes.SelectStmt:
		return CreateSelectStatement(stmt), nil
	// case nodes.SetOperationStmt:
	case nodes.TransactionStmt:
		return CreateTransactionStatement(stmt), nil
	// case nodes.TruncateStmt:
	// case nodes.UnlistenStmt:
	// case nodes.UpdateStmt:
	//     // return nil, _update.HandleUpdate(ctx, stmt)
	// case nodes.VacuumStmt:
	case nodes.VariableSetStmt:
		return CreateVariableSetStatement(stmt), nil
	case nodes.VariableShowStmt:
		return CreateVariableShowStatement(stmt), nil
	// case nodes.ViewStmt:
	default:
		return nil, errors.New("invalid or unsupported nodes type")
	}
}
