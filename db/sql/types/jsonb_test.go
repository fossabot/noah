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
	"bytes"
	"reflect"
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestJSONBTranscode(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, conn)
	if _, ok := conn.ConnInfo.DataTypeForName("jsonb"); !ok {
		t.Skip("Skipping due to no jsonb type")
	}

	testutil.TestSuccessfulTranscode(t, "jsonb", []interface{}{
		&types.JSONB{Bytes: []byte("{}"), Status: types.Present},
		&types.JSONB{Bytes: []byte("null"), Status: types.Present},
		&types.JSONB{Bytes: []byte("42"), Status: types.Present},
		&types.JSONB{Bytes: []byte(`"hello"`), Status: types.Present},
		&types.JSONB{Status: types.Null},
	})
}

func TestJSONBSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.JSONB
	}{
		{source: "{}", result: types.JSONB{Bytes: []byte("{}"), Status: types.Present}},
		{source: []byte("{}"), result: types.JSONB{Bytes: []byte("{}"), Status: types.Present}},
		{source: ([]byte)(nil), result: types.JSONB{Status: types.Null}},
		{source: (*string)(nil), result: types.JSONB{Status: types.Null}},
		{source: []int{1, 2, 3}, result: types.JSONB{Bytes: []byte("[1,2,3]"), Status: types.Present}},
		{source: map[string]interface{}{"foo": "bar"}, result: types.JSONB{Bytes: []byte(`{"foo":"bar"}`), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var d types.JSONB
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(d, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestJSONBAssignTo(t *testing.T) {
	var s string
	var ps *string
	var b []byte

	rawStringTests := []struct {
		src      types.JSONB
		dst      *string
		expected string
	}{
		{src: types.JSONB{Bytes: []byte("{}"), Status: types.Present}, dst: &s, expected: "{}"},
	}

	for i, tt := range rawStringTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if *tt.dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}

	rawBytesTests := []struct {
		src      types.JSONB
		dst      *[]byte
		expected []byte
	}{
		{src: types.JSONB{Bytes: []byte("{}"), Status: types.Present}, dst: &b, expected: []byte("{}")},
		{src: types.JSONB{Status: types.Null}, dst: &b, expected: (([]byte)(nil))},
	}

	for i, tt := range rawBytesTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if bytes.Compare(tt.expected, *tt.dst) != 0 {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}

	var mapDst map[string]interface{}
	type structDst struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var strDst structDst

	unmarshalTests := []struct {
		src      types.JSONB
		dst      interface{}
		expected interface{}
	}{
		{src: types.JSONB{Bytes: []byte(`{"foo":"bar"}`), Status: types.Present}, dst: &mapDst, expected: map[string]interface{}{"foo": "bar"}},
		{src: types.JSONB{Bytes: []byte(`{"name":"John","age":42}`), Status: types.Present}, dst: &strDst, expected: structDst{Name: "John", Age: 42}},
	}
	for i, tt := range unmarshalTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.JSONB
		dst      **string
		expected *string
	}{
		{src: types.JSONB{Status: types.Null}, dst: &ps, expected: ((*string)(nil))},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if *tt.dst == tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}
}
