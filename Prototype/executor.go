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
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package Prototype

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/Ready-Stock/noah/Prototype/context"
    _comment "github.com/Ready-Stock/noah/Prototype/queries/comment"
    _create "github.com/Ready-Stock/noah/Prototype/queries/create"
    _drop "github.com/Ready-Stock/noah/Prototype/queries/drop"
    _insert "github.com/Ready-Stock/noah/Prototype/queries/insert"
    "github.com/Ready-Stock/noah/Prototype/queries/select"
    _transaction "github.com/Ready-Stock/noah/Prototype/queries/transaction"
    _update "github.com/Ready-Stock/noah/Prototype/queries/update"
    pgq "github.com/Ready-Stock/pg_query_go"
    query "github.com/Ready-Stock/pg_query_go/nodes"
)

func Start() context.SessionContext {
	return context.SessionContext{
		TransactionState: context.StateNoTxn,
		Nodes:            map[int]context.NodeContext{},
	}
}

func InjestQuery(ctx *context.SessionContext, query string) (*sql.Rows, error) {
	if parsed, err := pgq.Parse(query); err != nil {
		return nil, err
	} else {
		fmt.Printf("Injesting Query: %s \n", query)
		return handleParseTree(ctx, *parsed)
	}
}

func handleParseTree(ctx *context.SessionContext, tree pgq.ParsetreeList) (*sql.Rows, error) {
	if len(tree.Statements) == 0 {
		return nil, errors.New("no statements provided")
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
			return nil, _comment.CreateCommentStatment(stmt, tree).HandleComment(ctx)
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
			return nil, _create.CreateCreateStatment(stmt, tree).HandleCreate(ctx)
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
			return nil, _drop.CreateDropStatment(stmt, tree).HandleComment(ctx)
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
			return nil, _insert.CreateInsertStatment(stmt, tree).HandleInsert(ctx)
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
			return _select.CreateSelectStatement(stmt, tree).HandleSelect(ctx)
		case query.SetOperationStmt:
		case query.TransactionStmt:
			return nil, _transaction.HandleTransaction(ctx, stmt)
		case query.TruncateStmt:
		case query.UnlistenStmt:
		case query.UpdateStmt:
			return nil, _update.HandleUpdate(ctx, stmt)
		case query.VacuumStmt:
		case query.VariableSetStmt:
		case query.VariableShowStmt:
		case query.ViewStmt:
		default:
			return nil, errors.New("invalid or unsupported query type")
		}
		return nil, nil
	}
	return nil, nil
}
