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
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package sql

import (
    "errors"
    nodes "github.com/readystock/pg_query_go/nodes"
)

func (ex *connExecutor) execStmt(
    tree nodes.Stmt,
    res RestrictedCommandResult,
    pos CmdPos,
) error {
    switch tree.StatementType() {
    case nodes.DDL, nodes.RowsAffected:

    }

    handler, err := ex.getStatementHandler(tree)
    if err != nil {
        return err
    }

    return handler.Execute(ex, res)
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
    // case nodes.DeleteStmt:
    // case nodes.DiscardStmt:
    // case nodes.DoStmt:
    // case nodes.DropOwnedStmt:
    // case nodes.DropRoleStmt:
    // case nodes.DropStmt:
    //     // return nil, _drop.CreateDropStatment(stmt, tree).HandleComment(ctx)
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
