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
	"context"
	"github.com/Ready-Stock/noah/db/sql/pgwire/pgerror"
)

func (ex *connExecutor) execDescribe(
	ctx context.Context, descCmd DescribeStmt, res DescribeResult,
) error {
	return pgerror.NewErrorf(
		pgerror.CodeDataExceptionError,
		"unknown prepared statement %q", descCmd.Name)
	// retErr := func(err error) (fsm.Event, fsm.EventPayload) {
	// 	return eventNonRetriableErr{IsCommit: fsm.False}, eventNonRetriableErrPayload{err: err}
	// }

	// switch descCmd.Type {
	// case pgwirebase.PrepareStatement:
	// 	ps, ok := ex.prepStmtsNamespace.prepStmts[descCmd.Name]
	// 	if !ok {
	// 		return pgerror.NewErrorf(
	// 			pgerror.CodeInvalidSQLStatementNameError,
	// 			"unknown prepared statement %q", descCmd.Name)
	// 	}
	//
	// 	res.SetInTypes(ps.InTypes)
	//
	// 	if stmtHasNoData(ps.Statement) {
	// 		res.SetNoDataRowDescription()
	// 	} else {
	// 		res.SetPrepStmtOutput(ctx, ps.Columns)
	// 	}
	// case pgwirebase.PreparePortal:
	// 	portal, ok := ex.prepStmtsNamespace.portals[descCmd.Name]
	// 	if !ok {
	// 		return pgerror.NewErrorf(
	// 			pgerror.CodeInvalidCursorNameError, "unknown portal %q", descCmd.Name)
	// 	}
	//
	// 	if stmtHasNoData(portal.Stmt.Statement) {
	// 		res.SetNoDataRowDescription()
	// 	} else {
	// 		res.SetPortalOutput(ctx, portal.Stmt.Columns, portal.OutFormats)
	// 	}
	// default:
	// 	return errors.Errorf("unknown describe type: %s", descCmd.Type)
	// }
	// return nil
}
