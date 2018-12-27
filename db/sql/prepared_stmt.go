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

package sql

import (
    "context"
    "github.com/readystock/noah/db/sql/pgwire/pgwirebase"
    "github.com/readystock/noah/db/sql/sem/tree"
    nodes "github.com/readystock/pg_query_go/nodes"
)

// PreparedStatement is a SQL statement that has been parsed and the types
// of arguments and results have been determined.
type PreparedStatement struct {
    // Str is the statement string prior to parsing, used to generate
    // error messages. This may be used in
    // the future to present a contextual error message based on location
    // information.
    Str string

    // Statement is the parse tree from pg_query.
    // This is used later to modify the query on the fly.
    Statement *nodes.Stmt

    // TODO(andrei): The connExecutor doesn't use this. Delete it once the
    // Executor is gone.
    portalNames map[string]struct{}
}

func (p *PreparedStatement) close() {
}

// preparedStatementsAccessor gives a planner access to a session's collection
// of prepared statements.
type preparedStatementsAccessor interface {
    // Get returns the prepared statement with the given name. The returned bool
    // is false if a statement with the given name doesn't exist.
    Get(name string) (*PreparedStatement, bool)
    // Delete removes the PreparedStatement with the provided name from the
    // collection. If a portal exists for that statement, it is also removed.
    // The method returns true if statement with that name was found and removed,
    // false otherwise.
    Delete(ctx context.Context, name string) bool
    // DeleteAll removes all prepared statements and portals from the coolection.
    DeleteAll(ctx context.Context)
}

// PreparedPortal is a PreparedStatement that has been bound with query arguments.
type PreparedPortal struct {
    Stmt  *PreparedStatement
    Qargs tree.QueryArguments

    // OutFormats contains the requested formats for the output columns.
    OutFormats []pgwirebase.FormatCode
}

// newPreparedPortal creates a new PreparedPortal.
//
// When no longer in use, the PreparedPortal needs to be close()d.
func (ex *connExecutor) newPreparedPortal(stmt *PreparedStatement) PreparedPortal {
    portal := PreparedPortal{
        Stmt: stmt,
    }
    return portal
}

func (p *PreparedPortal) close() {
}
