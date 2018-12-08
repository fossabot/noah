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

package tree

import (
	"strings"

	"github.com/readystock/noah/db/sql/sessiondata"
)

// GetRenderColName computes a name for a result column.
// A name specified with AS takes priority, otherwise a name
// is derived from the expression.
//
// This function is meant to be used on untransformed syntax trees.
//
// The algorithm is borrowed from FigureColName() in PostgreSQL 10, to be
// found in src/backend/parser/parse_target.c. We reuse this algorithm
// to provide names more compatible with PostgreSQL.
func GetRenderColName(searchPath sessiondata.SearchPath, target SelectExpr) (string, error) {
	if target.As != "" {
		return string(target.As), nil
	}

	_, s, err := ComputeColNameInternal(searchPath, target.Expr)
	if err != nil {
		return s, err
	}
	if len(s) == 0 {
		s = "?column?"
	}
	return s, nil
}

// ComputeColNameInternal is the workhorse for GetRenderColName.
// The return value indicates the strength of the confidence in the result:
// 0 - no information
// 1 - second-best name choice
// 2 - good name choice
//
// The algorithm is borrowed from FigureColnameInternal in PostgreSQL 10,
// to be found in src/backend/parser/parse_target.c.
func ComputeColNameInternal(sp sessiondata.SearchPath, target Expr) (int, string, error) {
	// The order of the type cases below mirrors that of PostgreSQL's
	// own code, so that code reviews can more easily compare the two
	// implementations.
	switch e := target.(type) {
	case *UnresolvedName:
		if e.Star {
			return 0, "", nil
		}
		return 2, e.Parts[0], nil

	case *ColumnItem:
		return 2, e.Column(), nil

	case *IndirectionExpr:
		return ComputeColNameInternal(sp, e.Expr)

	case *FuncExpr:
		fd, err := e.Func.Resolve(sp)
		if err != nil {
			return 0, "", err
		}
		return 2, fd.Name, nil

	case *NullIfExpr:
		return 2, "nullif", nil

	case *IfExpr:
		return 2, "if", nil

	case *ParenExpr:
		return ComputeColNameInternal(sp, e.Expr)

	case *CastExpr:
		strength, s, err := ComputeColNameInternal(sp, e.Expr)
		if err != nil {
			return 0, "", err
		}
		if strength <= 1 {
			tname := strings.ToLower(e.Type.TypeName())
			// TTuple has no short time name, so check this
			// here. Otherwise we'll want to fall back below.
			if tname != "" {
				return 1, tname, nil
			}
		}
		return strength, s, nil

	case *AnnotateTypeExpr:
		// Ditto CastExpr.
		strength, s, err := ComputeColNameInternal(sp, e.Expr)
		if err != nil {
			return 0, "", err
		}
		if strength <= 1 {
			tname := strings.ToLower(e.Type.TypeName())
			// TTuple has no short time name, so check this
			// here. Otherwise we'll want to fall back below.
			if tname != "" {
				return 1, tname, nil
			}
		}
		return strength, s, nil

	case *CollateExpr:
		return ComputeColNameInternal(sp, e.Expr)

	case *ArrayFlatten:
		return 2, "array", nil

	case *Subquery:
		if e.Exists {
			return 2, "exists", nil
		}
		return computeColNameInternalSubquery(sp, e.Select)

	case *CaseExpr:
		strength, s, err := 0, "", error(nil)
		if e.Else != nil {
			strength, s, err = ComputeColNameInternal(sp, e.Else)
		}
		if strength <= 1 {
			s = "case"
			strength = 1
		}
		return strength, s, err

	case *Array:
		return 2, "array", nil

	case *Tuple:
		if e.Row {
			return 2, "row", nil
		}
		if len(e.Exprs) == 1 {
			if len(e.Labels) > 0 {
				return 2, e.Labels[0], nil
			}
			return ComputeColNameInternal(sp, e.Exprs[0])
		}

	case *CoalesceExpr:
		return 2, "coalesce", nil

		// CockroachDB-specific nodes follow.
	case *IfErrExpr:
		if e.Else == nil {
			return 2, "iserror", nil
		}
		return 2, "iferror", nil

	case *ColumnAccessExpr:
		return 2, e.ColName, nil

	case *DBool:
		// PostgreSQL implements the "true" and "false" literals
		// by generating the expressions 't'::BOOL and 'f'::BOOL, so
		// the derived column name is just "bool". Do the same.
		return 1, "bool", nil
	}

	return 0, "", nil
}

// computeColNameInternalSubquery handles the cases of subqueries that
// cannot be handled by the function above due to the Go typing
// differences.
func computeColNameInternalSubquery(
	sp sessiondata.SearchPath, s SelectStatement,
) (int, string, error) {
	switch e := s.(type) {
	case *ParenSelect:
		return computeColNameInternalSubquery(sp, e.Select.Select)
	case *ValuesClause:
		if len(e.Tuples) > 0 && len(e.Tuples[0].Exprs) == 1 {
			return 2, "column1", nil
		}
	case *SelectClause:
		if len(e.Exprs) == 1 {
			if len(e.Exprs[0].As) > 0 {
				return 2, string(e.Exprs[0].As), nil
			}
			return ComputeColNameInternal(sp, e.Exprs[0].Expr)
		}
	}
	return 0, "", nil
}
