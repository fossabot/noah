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

// formatNodeOrHideConstants recurses into a node for pretty-printing,
// unless hideConstants is set in the flags and the node is a datum or
// a literal.
func (ctx *FmtCtx) formatNodeOrHideConstants(n NodeFormatter) {
    if ctx.flags.HasFlags(FmtHideConstants) {
        switch v := n.(type) {
        case *ValuesClause:
            v.formatHideConstants(ctx)
            return
        case *ComparisonExpr:
            if v.Operator == In || v.Operator == NotIn {
                if t, ok := v.Right.(*Tuple); ok {
                    v.formatInTupleAndHideConstants(ctx, t)
                    return
                }
            }
        case *Placeholder:
            // Placeholders should be printed as placeholder markers.
            // Deliberately empty so we format as normal.
        case Datum, Constant:
            ctx.WriteByte('_')
            return
        }
    }
    n.Format(ctx)
}

// formatInTupleAndHideConstants formats an "a IN (...)" expression
// and collapses the tuple on the right to contain at most 2 elements
// if it otherwise only contains literals or placeholders.
// e.g.:
//    a IN (1, 2, 3)       -> a IN (_, _)
//    a IN (x+1, x+2, x+3) -> a IN (x+_, x+_, x+_)
func (node *ComparisonExpr) formatInTupleAndHideConstants(ctx *FmtCtx, rightTuple *Tuple) {
    exprFmtWithParen(ctx, node.Left)
    ctx.WriteByte(' ')
    ctx.WriteString(node.Operator.String())
    ctx.WriteByte(' ')
    rightTuple.formatHideConstants(ctx)
}

// formatHideConstants shortens multi-valued VALUES clauses to a
// VALUES clause with a single value.
// e.g. VALUES (a,b,c), (d,e,f) -> VALUES (_, _, _)
func (node *ValuesClause) formatHideConstants(ctx *FmtCtx) {
    ctx.WriteString("VALUES ")
    node.Tuples[0].Format(ctx)
}

// formatHideConstants formats tuples containing only literals or
// placeholders and longer than 1 element as a tuple of its first
// two elements, scrubbed.
// e.g. (1)               -> (_)
//      (1, 2)            -> (_, _)
//      (1, 2, 3)         -> (_, _)
//      ROW()             -> ROW()
//      ROW($1, $2, $3)   -> ROW($1, $2)
//      (1+2, 2+3, 3+4)   -> (_ + _, _ + _, _ + _)
//      (1+2, b, c)       -> (_ + _, b, c)
func (node *Tuple) formatHideConstants(ctx *FmtCtx) {
    if len(node.Exprs) < 2 {
        node.Format(ctx)
        return
    }

    // First, determine if there are only literals/placeholders.
    var i int
    for i = 0; i < len(node.Exprs); i++ {
        switch node.Exprs[i].(type) {
        case Datum, Constant, *Placeholder:
            continue
        }
        break
    }
    // If so, then use the special representation.
    if i == len(node.Exprs) {
        // We copy the node to preserve the "row" boolean flag.
        v2 := *node
        v2.Exprs = v2.Exprs[:2]
        v2.Format(ctx)
        return
    }
    node.Format(ctx)
}
