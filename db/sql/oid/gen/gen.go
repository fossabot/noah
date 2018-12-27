// +build ignore

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

// Generate the table of OID values
// Run with 'go run gen.go'.
package gen

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    _ "github.com/lib/pq"
)

// OID represent a postgres Object Identifier Type.
type OID struct {
    ID   int
    Type string
}

// Name returns an upper case version of the oid type.
func (o OID) Name() string {
    return strings.ToUpper(o.Type)
}

func main() {
    datname := os.Getenv("PGDATABASE")
    sslmode := os.Getenv("PGSSLMODE")

    if datname == "" {
        os.Setenv("PGDATABASE", "pqgotest")
    }

    if sslmode == "" {
        os.Setenv("PGSSLMODE", "disable")
    }

    db, err := sql.Open("postgres", "")
    if err != nil {
        log.Fatal(err)
    }
    rows, err := db.Query(`
        SELECT typname, oid
        FROM pg_type WHERE oid < 10000
        ORDER BY oid;
    `)
    if err != nil {
        log.Fatal(err)
    }
    oids := make([]*OID, 0)
    for rows.Next() {
        var oid OID
        if err = rows.Scan(&oid.Type, &oid.ID); err != nil {
            log.Fatal(err)
        }
        oids = append(oids, &oid)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
    cmd := exec.Command("gofmt")
    cmd.Stderr = os.Stderr
    w, err := cmd.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }
    f, err := os.Create("types.go")
    if err != nil {
        log.Fatal(err)
    }
    cmd.Stdout = f
    err = cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Fprintln(w, "// Code generated by gen.go. DO NOT EDIT.")
    fmt.Fprintln(w, "\npackage oid")
    fmt.Fprintln(w, "const (")
    for _, oid := range oids {
        fmt.Fprintf(w, "T_%s Oid = %d\n", oid.Type, oid.ID)
    }
    fmt.Fprintln(w, ")")
    fmt.Fprintln(w, "var TypeName = map[Oid]string{")
    for _, oid := range oids {
        fmt.Fprintf(w, "T_%s: \"%s\",\n", oid.Type, oid.Name())
    }
    fmt.Fprintln(w, "}")
    w.Close()
    cmd.Wait()
}
