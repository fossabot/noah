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
 */

package types_test

import (
    "reflect"
    "testing"

    "github.com/Ready-Stock/noah/db/sql/types"
    "github.com/Ready-Stock/noah/db/sql/types/testutil"
    "github.com/Ready-Stock/pgx"
)

func TestHstoreArrayTranscode(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, conn)

	text := func(s string) types.Text {
		return types.Text{String: s, Status: types.Present}
	}

	values := []types.Hstore{
		{Map: map[string]types.Text{}, Status: types.Present},
		{Map: map[string]types.Text{"foo": text("bar")}, Status: types.Present},
		{Map: map[string]types.Text{"foo": text("bar"), "baz": text("quz")}, Status: types.Present},
		{Map: map[string]types.Text{"NULL": text("bar")}, Status: types.Present},
		{Map: map[string]types.Text{"foo": text("NULL")}, Status: types.Present},
		{Status: types.Null},
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
		values = append(values, types.Hstore{Map: map[string]types.Text{s + "foo": text("bar")}, Status: types.Present})         // at beginning
		values = append(values, types.Hstore{Map: map[string]types.Text{"foo" + s + "bar": text("bar")}, Status: types.Present}) // in middle
		values = append(values, types.Hstore{Map: map[string]types.Text{"foo" + s: text("bar")}, Status: types.Present})         // at end
		values = append(values, types.Hstore{Map: map[string]types.Text{s: text("bar")}, Status: types.Present})                 // is key

		// Special value values
		values = append(values, types.Hstore{Map: map[string]types.Text{"foo": text(s + "bar")}, Status: types.Present})         // at beginning
		values = append(values, types.Hstore{Map: map[string]types.Text{"foo": text("foo" + s + "bar")}, Status: types.Present}) // in middle
		values = append(values, types.Hstore{Map: map[string]types.Text{"foo": text("foo" + s)}, Status: types.Present})         // at end
		values = append(values, types.Hstore{Map: map[string]types.Text{"foo": text(s)}, Status: types.Present})                 // is key
	}

	src := &types.HstoreArray{
		Elements:   values,
		Dimensions: []types.ArrayDimension{{Length: int32(len(values)), LowerBound: 1}},
		Status:     types.Present,
	}

	ps, err := conn.Prepare("test", "select $1::hstore[]")
	if err != nil {
		t.Fatal(err)
	}

	formats := []struct {
		name       string
		formatCode int16
	}{
		{name: "TextFormat", formatCode: pgx.TextFormatCode},
		{name: "BinaryFormat", formatCode: pgx.BinaryFormatCode},
	}

	for _, fc := range formats {
		ps.FieldDescriptions[0].FormatCode = fc.formatCode
		vEncoder := testutil.ForceEncoder(src, fc.formatCode)
		if vEncoder == nil {
			t.Logf("%#v does not implement %v", src, fc.name)
			continue
		}

		var result types.HstoreArray
		err := conn.QueryRow("test", vEncoder).Scan(&result)
		if err != nil {
			t.Errorf("%v: %v", fc.name, err)
			continue
		}

		if result.Status != src.Status {
			t.Errorf("%v: expected Status %v, got %v", fc.formatCode, src.Status, result.Status)
			continue
		}

		if len(result.Elements) != len(src.Elements) {
			t.Errorf("%v: expected %v elements, got %v", fc.formatCode, len(src.Elements), len(result.Elements))
			continue
		}

		for i := range result.Elements {
			a := src.Elements[i]
			b := result.Elements[i]

			if a.Status != b.Status {
				t.Errorf("%v element idx %d: expected status %v, got %v", fc.formatCode, i, a.Status, b.Status)
			}

			if len(a.Map) != len(b.Map) {
				t.Errorf("%v element idx %d: expected %v pairs, got %v", fc.formatCode, i, len(a.Map), len(b.Map))
			}

			for k := range a.Map {
				if a.Map[k] != b.Map[k] {
					t.Errorf("%v element idx %d: expected key %v to be %v, got %v", fc.formatCode, i, k, a.Map[k], b.Map[k])
				}
			}
		}
	}
}

func TestHstoreArraySet(t *testing.T) {
	successfulTests := []struct {
		src    []map[string]string
		result types.HstoreArray
	}{
		{
			src: []map[string]string{{"foo": "bar"}},
			result: types.HstoreArray{
				Elements: []types.Hstore{
					{
						Map:    map[string]types.Text{"foo": {String: "bar", Status: types.Present}},
						Status: types.Present,
					},
				},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
		},
	}

	for i, tt := range successfulTests {
		var dst types.HstoreArray
		err := dst.Set(tt.src)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(dst, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.src, tt.result, dst)
		}
	}
}

func TestHstoreArrayAssignTo(t *testing.T) {
	var m []map[string]string

	simpleTests := []struct {
		src      types.HstoreArray
		dst      *[]map[string]string
		expected []map[string]string
	}{
		{
			src: types.HstoreArray{
				Elements: []types.Hstore{
					{
						Map:    map[string]types.Text{"foo": {String: "bar", Status: types.Present}},
						Status: types.Present,
					},
				},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &m,
			expected: []map[string]string{{"foo": "bar"}}},
        {src: types.HstoreArray{Status: types.Null}, dst: &m, expected: ([]map[string]string)(nil)},
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
