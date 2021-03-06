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
	"github.com/pkg/errors"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/pgwire/pgerror"
	"github.com/readystock/noah/db/sql/pgwire/pgwirebase"
)

func (ex *connExecutor) execDescribe(
	ctx context.Context, descCmd DescribeStmt, res DescribeResult,
) error {
	// return pgerror.NewErrorf(
	//     pgerror.CodeDataExceptionError,
	//     "unknown prepared statement %q", descCmd.Name)

	switch descCmd.Type {
	case pgwirebase.PrepareStatement:
		ps, ok := ex.prepStmtsNamespace.prepStmts[descCmd.Name]
		if !ok {
			return pgerror.NewErrorf(
				pgerror.CodeInvalidSQLStatementNameError,
				"unknown prepared statement %q", descCmd.Name)
		}

		res.SetInTypes(ps.InTypes)

		if stmtHasNoData(*ps.Statement) {
			res.SetNoDataRowDescription()
		} else {
			// golog.Warnf("cannot handle prepared statement output")
			// res.SetPrepStmtOutput(ctx, ps.Columns)
		}
	case pgwirebase.PreparePortal:
		portal, ok := ex.prepStmtsNamespace.portals[descCmd.Name]
		if !ok {
			return pgerror.NewErrorf(
				pgerror.CodeInvalidCursorNameError, "unknown portal %q", descCmd.Name)
		}

		if stmtHasNoData(*portal.Stmt.Statement) {
			res.SetNoDataRowDescription()
		} else {
			golog.Warnf("cannot handle portal output of statement")
			// res.SetPortalOutput(ctx, portal.Stmt.Columns, portal.OutFormats)
		}
	default:
		return errors.Errorf("unknown describe type: %s", descCmd.Type)
	}
	return nil
}
