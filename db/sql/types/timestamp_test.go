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
    "time"

    "github.com/readystock/noah/db/sql/types"
    "github.com/readystock/noah/db/sql/types/testutil"
)

func TestTimestampTranscode(t *testing.T) {
    testutil.TestSuccessfulTranscodeEqFunc(t, "timestamp", []interface{}{
        &types.Timestamp{Time: time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(1905, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
        &types.Timestamp{Status: types.Null},
        &types.Timestamp{Status: types.Present, InfinityModifier: types.Infinity},
        &types.Timestamp{Status: types.Present, InfinityModifier: -types.Infinity},
    }, func(a, b interface{}) bool {
        at := a.(types.Timestamp)
        bt := b.(types.Timestamp)

        return at.Time.Equal(bt.Time) && at.Status == bt.Status && at.InfinityModifier == bt.InfinityModifier
    })
}

func TestTimestampSet(t *testing.T) {
    type _time time.Time

    successfulTests := []struct {
        source interface{}
        result types.Timestamp
    }{
        {source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
        {source: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
        {source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC), result: types.Timestamp{Time: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC), Status: types.Present}},
        {source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
        {source: time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), result: types.Timestamp{Time: time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), Status: types.Present}},
        {source: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
        {source: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), result: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
        {source: _time(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)), result: types.Timestamp{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
    }

    for i, tt := range successfulTests {
        var r types.Timestamp
        err := r.Set(tt.source)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if r != tt.result {
            t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
        }
    }
}

func TestTimestampAssignTo(t *testing.T) {
    var tim time.Time
    var ptim *time.Time

    simpleTests := []struct {
        src      types.Timestamp
        dst      interface{}
        expected interface{}
    }{
        {src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}, dst: &tim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
        {src: types.Timestamp{Time: time.Time{}, Status: types.Null}, dst: &ptim, expected: (*time.Time)(nil)},
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
        src      types.Timestamp
        dst      interface{}
        expected interface{}
    }{
        {src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}, dst: &ptim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)},
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
        src types.Timestamp
        dst interface{}
    }{
        {src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.Infinity, Status: types.Present}, dst: &tim},
        {src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.NegativeInfinity, Status: types.Present}, dst: &tim},
        {src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Null}, dst: &tim},
    }

    for i, tt := range errorTests {
        err := tt.src.AssignTo(tt.dst)
        if err == nil {
            t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
        }
    }
}
