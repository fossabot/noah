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

package tree

// Backup represents a BACKUP statement.
type Backup struct {
    Targets         TargetList
    To              Expr
    IncrementalFrom Exprs
    AsOf            AsOfClause
    Options         KVOptions
}

var _ Statement = &Backup{}

// Format implements the NodeFormatter interface.
func (node *Backup) Format(ctx *FmtCtx) {
    ctx.WriteString("BACKUP ")
    ctx.FormatNode(&node.Targets)
    ctx.WriteString(" TO ")
    ctx.FormatNode(node.To)
    if node.AsOf.Expr != nil {
        ctx.WriteString(" ")
        ctx.FormatNode(&node.AsOf)
    }
    if node.IncrementalFrom != nil {
        ctx.WriteString(" INCREMENTAL FROM ")
        ctx.FormatNode(&node.IncrementalFrom)
    }
    if node.Options != nil {
        ctx.WriteString(" WITH ")
        ctx.FormatNode(&node.Options)
    }
}

// Restore represents a RESTORE statement.
type Restore struct {
    Targets TargetList
    From    Exprs
    AsOf    AsOfClause
    Options KVOptions
}

var _ Statement = &Restore{}

// Format implements the NodeFormatter interface.
func (node *Restore) Format(ctx *FmtCtx) {
    ctx.WriteString("RESTORE ")
    ctx.FormatNode(&node.Targets)
    ctx.WriteString(" FROM ")
    ctx.FormatNode(&node.From)
    if node.AsOf.Expr != nil {
        ctx.WriteString(" ")
        ctx.FormatNode(&node.AsOf)
    }
    if node.Options != nil {
        ctx.WriteString(" WITH ")
        ctx.FormatNode(&node.Options)
    }
}

// KVOption is a key-value option.
type KVOption struct {
    Key   Name
    Value Expr
}

// KVOptions is a list of KVOptions.
type KVOptions []KVOption

// Format implements the NodeFormatter interface.
func (o *KVOptions) Format(ctx *FmtCtx) {
    for i := range *o {
        n := &(*o)[i]
        if i > 0 {
            ctx.WriteString(", ")
        }
        ctx.FormatNode(&n.Key)
        if n.Value != nil {
            ctx.WriteString(` = `)
            ctx.FormatNode(n.Value)
        }
    }
}
