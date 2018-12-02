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
	"fmt"
    "github.com/Ready-Stock/noah/db/sql/pgwire/pgerror"
)

// Table names are used in statements like CREATE TABLE,
// INSERT INTO, etc.
// General syntax:
//    [ [ <database-name> '.' ] <schema-name> '.' ] <table-name>
//
// The other syntax nodes hold a mutable NormalizableTableName
// attribute. This is populated during parsing with an
// UnresolvedName, and gets assigned an actual TableName upon the first
// call to its Normalize() method.

// NormalizableTableName implements an editable table name.
type NormalizableTableName struct {
	TableNameReference
}

// Format implements the NodeFormatter interface.
func (nt *NormalizableTableName) Format(ctx *FmtCtx) {
	if ctx.tableNameFormatter != nil {
		ctx.tableNameFormatter(ctx, nt)
	} else {
		ctx.FormatNode(nt.TableNameReference)
	}
}
func (nt *NormalizableTableName) String() string { return AsString(nt) }

// Normalize checks if the table name is already normalized and
// normalizes it as necessary. It stores the result of normalization
// so that it does not need to occur the next time around.
func (nt *NormalizableTableName) Normalize() (*TableName, error) {
	switch t := nt.TableNameReference.(type) {
	case *TableName:
		return t, nil
	case *UnresolvedName:
		tn, err := NormalizeTableName(t)
		if err != nil {
			return nil, err
		}
		nt.TableNameReference = &tn
		return &tn, nil
	default:
		return nil, pgerror.NewErrorWithDepthf(1, pgerror.CodeInternalError,
			"programming error: unsupported table name reference: %+v (%T)",
			nt.TableNameReference, nt.TableNameReference)
	}
}

// TableName asserts that the table name has been previously normalized.
func (nt *NormalizableTableName) TableName() *TableName {
	return nt.TableNameReference.(*TableName)
}

// tableExpr implements the TableExpr interface.
func (*NormalizableTableName) tableExpr() {}

// TableNameReference implements the editable cell of a TableExpr that
// refers to a single table.
type TableNameReference interface {
	fmt.Stringer
	NodeFormatter

	tblNameReference()
}

func (t *TableName) tblNameReference()      {}
func (u *UnresolvedName) tblNameReference() {}

// NormalizableTableNames corresponds to a comma-delimited list of
// normalizable table names.
type NormalizableTableNames []NormalizableTableName

// Format implements the NodeFormatter interface.
func (t *NormalizableTableNames) Format(ctx *FmtCtx) {
	sep := ""
	for i := range *t {
		ctx.WriteString(sep)
		ctx.FormatNode(&(*t)[i])
		sep = ", "
	}
}

// TableNameWithIndex represents a "table@index", used in statements that
// specifically refer to an index.
type TableNameWithIndex struct {
	Table NormalizableTableName
	Index UnrestrictedName

	// SearchTable indicates that we have just an index (no table name); we will
	// need to search for a table that has an index with the given name.
	//
	// To allow schema-qualified index names in this case, the index is actually
	// specified in Table as the table name, and Index is empty.
	SearchTable bool
}

// Format implements the NodeFormatter interface.
func (n *TableNameWithIndex) Format(ctx *FmtCtx) {
	ctx.FormatNode(&n.Table)
	if n.Index != "" {
		ctx.WriteByte('@')
		ctx.FormatNode(&n.Index)
	}
}
func (n *TableNameWithIndex) String() string { return AsString(n) }

// TableNameWithIndexList is a list of indexes.
type TableNameWithIndexList []*TableNameWithIndex

// Format implements the NodeFormatter interface.
func (n *TableNameWithIndexList) Format(ctx *FmtCtx) {
	sep := ""
	for _, tni := range *n {
		ctx.WriteString(sep)
		ctx.FormatNode(tni)
		sep = ", "
	}
}
