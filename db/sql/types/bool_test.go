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

package types_test

import (
    "reflect"
    "testing"

    "github.com/readystock/noah/db/sql/types"
    "github.com/readystock/noah/db/sql/types/testutil"
)

func TestBoolTranscode(t *testing.T) {
    testutil.TestSuccessfulTranscode(t, "bool", []interface{}{
        &types.Bool{Bool: false, Status: types.Present},
        &types.Bool{Bool: true, Status: types.Present},
        &types.Bool{Bool: false, Status: types.Null},
    })
}

func TestBoolSet(t *testing.T) {
    successfulTests := []struct {
        source interface{}
        result types.Bool
    }{
        {source: true, result: types.Bool{Bool: true, Status: types.Present}},
        {source: false, result: types.Bool{Bool: false, Status: types.Present}},
        {source: "true", result: types.Bool{Bool: true, Status: types.Present}},
        {source: "false", result: types.Bool{Bool: false, Status: types.Present}},
        {source: "t", result: types.Bool{Bool: true, Status: types.Present}},
        {source: "f", result: types.Bool{Bool: false, Status: types.Present}},
        {source: _bool(true), result: types.Bool{Bool: true, Status: types.Present}},
        {source: _bool(false), result: types.Bool{Bool: false, Status: types.Present}},
        {source: nil, result: types.Bool{Status: types.Null}},
    }

    for i, tt := range successfulTests {
        var r types.Bool
        err := r.Set(tt.source)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if r != tt.result {
            t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
        }
    }
}

func TestBoolAssignTo(t *testing.T) {
    var b bool
    var _b _bool
    var pb *bool
    var _pb *_bool

    simpleTests := []struct {
        src      types.Bool
        dst      interface{}
        expected interface{}
    }{
        {src: types.Bool{Bool: false, Status: types.Present}, dst: &b, expected: false},
        {src: types.Bool{Bool: true, Status: types.Present}, dst: &b, expected: true},
        {src: types.Bool{Bool: false, Status: types.Present}, dst: &_b, expected: _bool(false)},
        {src: types.Bool{Bool: true, Status: types.Present}, dst: &_b, expected: _bool(true)},
        {src: types.Bool{Bool: false, Status: types.Null}, dst: &pb, expected: (*bool)(nil)},
        {src: types.Bool{Bool: false, Status: types.Null}, dst: &_pb, expected: (*_bool)(nil)},
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
        src      types.Bool
        dst      interface{}
        expected interface{}
    }{
        {src: types.Bool{Bool: true, Status: types.Present}, dst: &pb, expected: true},
        {src: types.Bool{Bool: true, Status: types.Present}, dst: &_pb, expected: _bool(true)},
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
}
