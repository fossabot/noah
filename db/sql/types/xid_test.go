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

package types_test

import (
    "reflect"
    "testing"

    "github.com/readystock/noah/db/sql/types"
    "github.com/readystock/noah/db/sql/types/testutil"
)

func TestXIDTranscode(t *testing.T) {
    typesName := "xid"
    values := []interface{}{
        &types.XID{Uint: 42, Status: types.Present},
        &types.XID{Status: types.Null},
    }
    eqFunc := func(a, b interface{}) bool {
        return reflect.DeepEqual(a, b)
    }

    testutil.TestPgxSuccessfulTranscodeEqFunc(t, typesName, values, eqFunc)

    // No direct conversion from int to xid, convert through text
    testutil.TestPgxSimpleProtocolSuccessfulTranscodeEqFunc(t, "text::"+typesName, values, eqFunc)

    for _, driverName := range []string{"github.com/lib/pq", "github.com/jackc/pgx/stdlib"} {
        testutil.TestDatabaseSQLSuccessfulTranscodeEqFunc(t, driverName, typesName, values, eqFunc)
    }
}

func TestXIDSet(t *testing.T) {
    successfulTests := []struct {
        source interface{}
        result types.XID
    }{
        {source: uint32(1), result: types.XID{Uint: 1, Status: types.Present}},
    }

    for i, tt := range successfulTests {
        var r types.XID
        err := r.Set(tt.source)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if r != tt.result {
            t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
        }
    }
}

func TestXIDAssignTo(t *testing.T) {
    var ui32 uint32
    var pui32 *uint32

    simpleTests := []struct {
        src      types.XID
        dst      interface{}
        expected interface{}
    }{
        {src: types.XID{Uint: 42, Status: types.Present}, dst: &ui32, expected: uint32(42)},
        {src: types.XID{Status: types.Null}, dst: &pui32, expected: (*uint32)(nil)},
    }

    for i, tt := range simpleTests {
        err := tt.src.AssignTo(tt.dst)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
            t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
        }
    }

    pointerAllocTests := []struct {
        src      types.XID
        dst      interface{}
        expected interface{}
    }{
        {src: types.XID{Uint: 42, Status: types.Present}, dst: &pui32, expected: uint32(42)},
    }

    for i, tt := range pointerAllocTests {
        err := tt.src.AssignTo(tt.dst)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); dst != tt.expected {
            t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
        }
    }

    errorTests := []struct {
        src types.XID
        dst interface{}
    }{
        {src: types.XID{Status: types.Null}, dst: &ui32},
    }

    for i, tt := range errorTests {
        err := tt.src.AssignTo(tt.dst)
        if err == nil {
            t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
        }
    }
}
