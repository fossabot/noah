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
    "fmt"
    "reflect"
    "testing"

    "github.com/Ready-Stock/noah/db/sql/types"
    "github.com/Ready-Stock/noah/db/sql/types/testutil"
    "github.com/Ready-Stock/pgx"
)

func TestRecordTranscode(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, conn)

	tests := []struct {
		sql      string
		expected types.Record
	}{
		{
			sql: `select row()`,
			expected: types.Record{
				Fields: []types.Value{},
				Status: types.Present,
			},
		},
		{
			sql: `select row('foo'::text, 42::int4)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select row(100.0::float4, 1.09::float4)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Float4{Float: 100, Status: types.Present},
					&types.Float4{Float: 1.09, Status: types.Present},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select row('foo'::text, array[1, 2, null, 4]::int4[], 42::int4)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4Array{
						Elements: []types.Int4{
							{Int: 1, Status: types.Present},
							{Int: 2, Status: types.Present},
							{Status: types.Null},
							{Int: 4, Status: types.Present},
						},
						Dimensions: []types.ArrayDimension{{Length: 4, LowerBound: 1}},
						Status:     types.Present,
					},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select row(null)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Unknown{Status: types.Null},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select null::record`,
			expected: types.Record{
				Status: types.Null,
			},
		},
	}

	for i, tt := range tests {
		psName := fmt.Sprintf("test%d", i)
		ps, err := conn.Prepare(psName, tt.sql)
		if err != nil {
			t.Fatal(err)
		}
		ps.FieldDescriptions[0].FormatCode = pgx.BinaryFormatCode

		var result types.Record
		if err := conn.QueryRow(psName).Scan(&result); err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(tt.expected, result) {
			t.Errorf("%d: expected %#v, got %#v", i, tt.expected, result)
		}
	}
}

func TestRecordWithUnknownOID(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, conn)

	_, err := conn.Exec(`drop type if exists floatrange;

create type floatrange as range (
  subtype = float8,
  subtype_diff = float8mi
);`)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Exec("drop type floatrange")

	var result types.Record
	err = conn.QueryRow("select row('foo'::text, floatrange(1, 10), 'bar'::text)").Scan(&result)
	if err == nil {
		t.Errorf("expected error but none")
	}
}

func TestRecordAssignTo(t *testing.T) {
	var valueSlice []types.Value
	var interfaceSlice []interface{}

	simpleTests := []struct {
		src      types.Record
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
			dst: &valueSlice,
			expected: []types.Value{
				&types.Text{String: "foo", Status: types.Present},
				&types.Int4{Int: 42, Status: types.Present},
			},
		},
		{
			src: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
			dst:      &interfaceSlice,
			expected: []interface{}{"foo", int32(42)},
		},
		{
			src:      types.Record{Status: types.Null},
			dst:      &valueSlice,
            expected: ([]types.Value)(nil),
		},
		{
			src:      types.Record{Status: types.Null},
			dst:      &interfaceSlice,
            expected: ([]interface{})(nil),
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}
}
