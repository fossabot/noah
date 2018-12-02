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
)

func TestByteaTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "bytea", []interface{}{
		&types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present},
		&types.Bytea{Bytes: []byte{}, Status: types.Present},
		&types.Bytea{Bytes: nil, Status: types.Null},
	})
}

func TestByteaSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Bytea
	}{
		{source: []byte{1, 2, 3}, result: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}},
		{source: []byte{}, result: types.Bytea{Bytes: []byte{}, Status: types.Present}},
		{source: []byte(nil), result: types.Bytea{Status: types.Null}},
		{source: _byteSlice{1, 2, 3}, result: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}},
		{source: _byteSlice(nil), result: types.Bytea{Status: types.Null}},
	}

	for i, tt := range successfulTests {
		var r types.Bytea
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestByteaAssignTo(t *testing.T) {
	var buf []byte
	var _buf _byteSlice
	var pbuf *[]byte
	var _pbuf *_byteSlice

	simpleTests := []struct {
		src      types.Bytea
		dst      interface{}
		expected interface{}
	}{
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &buf, expected: []byte{1, 2, 3}},
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &_buf, expected: _byteSlice{1, 2, 3}},
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &pbuf, expected: &[]byte{1, 2, 3}},
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &_pbuf, expected: &_byteSlice{1, 2, 3}},
        {src: types.Bytea{Status: types.Null}, dst: &pbuf, expected: (*[]byte)(nil)},
        {src: types.Bytea{Status: types.Null}, dst: &_pbuf, expected: (*_byteSlice)(nil)},
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
