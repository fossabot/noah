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

package sql

import (
    "context"
    "github.com/readystock/noah/db/sql/pgwire/pgwirebase"
    "github.com/readystock/noah/db/sql/plan"
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

    // TypeHints contains the types of the placeholders set by the client. It
    // dictates how input parameters for those placeholders will be parsed. If a
    // placeholder has no type hint, it will be populated during type checking.
    TypeHints plan.PlaceholderTypes

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
    // Qargs types.QueryArguments

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
