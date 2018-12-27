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
 */

// This code was derived from https://github.com/youtube/vitess.

package tree

import (
    "fmt"
)

// UnionClause represents a UNION statement.
type UnionClause struct {
    Type        UnionType
    Left, Right *Select
    All         bool
}

// UnionType represents one of the three set operations in sql.
type UnionType int

// Union.Type
const (
    UnionOp UnionType = iota
    IntersectOp
    ExceptOp
)

var unionTypeName = [...]string{
    UnionOp:     "UNION",
    IntersectOp: "INTERSECT",
    ExceptOp:    "EXCEPT",
}

func (i UnionType) String() string {
    if i < 0 || i > UnionType(len(unionTypeName)-1) {
        return fmt.Sprintf("UnionType(%d)", i)
    }
    return unionTypeName[i]
}

// Format implements the NodeFormatter interface.
func (node *UnionClause) Format(ctx *FmtCtx) {
    ctx.FormatNode(node.Left)
    ctx.WriteByte(' ')
    ctx.WriteString(node.Type.String())
    if node.All {
        ctx.WriteString(" ALL")
    }
    ctx.WriteByte(' ')
    ctx.FormatNode(node.Right)
}
