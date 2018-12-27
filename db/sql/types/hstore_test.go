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

func TestHstoreTranscode(t *testing.T) {
    text := func(s string) types.Text {
        return types.Text{String: s, Status: types.Present}
    }

    values := []interface{}{
        &types.Hstore{Map: map[string]types.Text{}, Status: types.Present},
        &types.Hstore{Map: map[string]types.Text{"foo": text("bar")}, Status: types.Present},
        &types.Hstore{Map: map[string]types.Text{"foo": text("bar"), "baz": text("quz")}, Status: types.Present},
        &types.Hstore{Map: map[string]types.Text{"NULL": text("bar")}, Status: types.Present},
        &types.Hstore{Map: map[string]types.Text{"foo": text("NULL")}, Status: types.Present},
        &types.Hstore{Status: types.Null},
    }

    specialStrings := []string{
        `"`,
        `'`,
        `\`,
        `\\`,
        `=>`,
        ` `,
        `\ / / \\ => " ' " '`,
    }
    for _, s := range specialStrings {
        // Special key values
        values = append(values, &types.Hstore{Map: map[string]types.Text{s + "foo": text("bar")}, Status: types.Present})         // at beginning
        values = append(values, &types.Hstore{Map: map[string]types.Text{"foo" + s + "bar": text("bar")}, Status: types.Present}) // in middle
        values = append(values, &types.Hstore{Map: map[string]types.Text{"foo" + s: text("bar")}, Status: types.Present})         // at end
        values = append(values, &types.Hstore{Map: map[string]types.Text{s: text("bar")}, Status: types.Present})                 // is key

        // Special value values
        values = append(values, &types.Hstore{Map: map[string]types.Text{"foo": text(s + "bar")}, Status: types.Present})         // at beginning
        values = append(values, &types.Hstore{Map: map[string]types.Text{"foo": text("foo" + s + "bar")}, Status: types.Present}) // in middle
        values = append(values, &types.Hstore{Map: map[string]types.Text{"foo": text("foo" + s)}, Status: types.Present})         // at end
        values = append(values, &types.Hstore{Map: map[string]types.Text{"foo": text(s)}, Status: types.Present})                 // is key
    }

    testutil.TestSuccessfulTranscodeEqFunc(t, "hstore", values, func(ai, bi interface{}) bool {
        a := ai.(types.Hstore)
        b := bi.(types.Hstore)

        if len(a.Map) != len(b.Map) || a.Status != b.Status {
            return false
        }

        for k := range a.Map {
            if a.Map[k] != b.Map[k] {
                return false
            }
        }

        return true
    })
}

func TestHstoreSet(t *testing.T) {
    successfulTests := []struct {
        src    map[string]string
        result types.Hstore
    }{
        {src: map[string]string{"foo": "bar"}, result: types.Hstore{Map: map[string]types.Text{"foo": {String: "bar", Status: types.Present}}, Status: types.Present}},
    }

    for i, tt := range successfulTests {
        var dst types.Hstore
        err := dst.Set(tt.src)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if !reflect.DeepEqual(dst, tt.result) {
            t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.src, tt.result, dst)
        }
    }
}

func TestHstoreAssignTo(t *testing.T) {
    var m map[string]string

    simpleTests := []struct {
        src      types.Hstore
        dst      *map[string]string
        expected map[string]string
    }{
        {src: types.Hstore{Map: map[string]types.Text{"foo": {String: "bar", Status: types.Present}}, Status: types.Present}, dst: &m, expected: map[string]string{"foo": "bar"}},
        {src: types.Hstore{Status: types.Null}, dst: &m, expected: (map[string]string)(nil)},
    }

    for i, tt := range simpleTests {
        err := tt.src.AssignTo(tt.dst)
        if err != nil {
            t.Errorf("%d: %v", i, err)
        }

        if !reflect.DeepEqual(*tt.dst, tt.expected) {
            t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
        }
    }
}
