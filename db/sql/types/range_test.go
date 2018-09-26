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
 */

package types

import (
	"bytes"
	"testing"
)

func TestParseUntypedTextRange(t *testing.T) {
	tests := []struct {
		src    string
		result UntypedTextRange
		err    error
	}{
		{
			src:    `[1,2)`,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `[1,2]`,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: Inclusive, UpperType: Inclusive},
			err:    nil,
		},
		{
			src:    `(1,3)`,
			result: UntypedTextRange{Lower: "1", Upper: "3", LowerType: Exclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    ` [1,2) `,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `[ foo , bar )`,
			result: UntypedTextRange{Lower: " foo ", Upper: " bar ", LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `["foo","bar")`,
			result: UntypedTextRange{Lower: "foo", Upper: "bar", LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `["f""oo","b""ar")`,
			result: UntypedTextRange{Lower: `f"oo`, Upper: `b"ar`, LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `["f""oo","b""ar")`,
			result: UntypedTextRange{Lower: `f"oo`, Upper: `b"ar`, LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `["","bar")`,
			result: UntypedTextRange{Lower: ``, Upper: `bar`, LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `[f\"oo\,,b\\ar\))`,
			result: UntypedTextRange{Lower: `f"oo,`, Upper: `b\ar)`, LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    `empty`,
			result: UntypedTextRange{Lower: "", Upper: "", LowerType: Empty, UpperType: Empty},
			err:    nil,
		},
	}

	for i, tt := range tests {
		r, err := ParseUntypedTextRange(tt.src)
		if err != tt.err {
			t.Errorf("%d. `%v`: expected err %v, got %v", i, tt.src, tt.err, err)
			continue
		}

		if r.LowerType != tt.result.LowerType {
			t.Errorf("%d. `%v`: expected result lower type %v, got %v", i, tt.src, string(tt.result.LowerType), string(r.LowerType))
		}

		if r.UpperType != tt.result.UpperType {
			t.Errorf("%d. `%v`: expected result upper type %v, got %v", i, tt.src, string(tt.result.UpperType), string(r.UpperType))
		}

		if r.Lower != tt.result.Lower {
			t.Errorf("%d. `%v`: expected result lower %v, got %v", i, tt.src, tt.result.Lower, r.Lower)
		}

		if r.Upper != tt.result.Upper {
			t.Errorf("%d. `%v`: expected result upper %v, got %v", i, tt.src, tt.result.Upper, r.Upper)
		}
	}
}

func TestParseUntypedBinaryRange(t *testing.T) {
	tests := []struct {
		src    []byte
		result UntypedBinaryRange
		err    error
	}{
		{
			src:    []byte{0, 0, 0, 0, 2, 0, 4, 0, 0, 0, 2, 0, 5},
			result: UntypedBinaryRange{Lower: []byte{0, 4}, Upper: []byte{0, 5}, LowerType: Exclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    []byte{1},
			result: UntypedBinaryRange{Lower: nil, Upper: nil, LowerType: Empty, UpperType: Empty},
			err:    nil,
		},
		{
			src:    []byte{2, 0, 0, 0, 2, 0, 4, 0, 0, 0, 2, 0, 5},
			result: UntypedBinaryRange{Lower: []byte{0, 4}, Upper: []byte{0, 5}, LowerType: Inclusive, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    []byte{4, 0, 0, 0, 2, 0, 4, 0, 0, 0, 2, 0, 5},
			result: UntypedBinaryRange{Lower: []byte{0, 4}, Upper: []byte{0, 5}, LowerType: Exclusive, UpperType: Inclusive},
			err:    nil,
		},
		{
			src:    []byte{6, 0, 0, 0, 2, 0, 4, 0, 0, 0, 2, 0, 5},
			result: UntypedBinaryRange{Lower: []byte{0, 4}, Upper: []byte{0, 5}, LowerType: Inclusive, UpperType: Inclusive},
			err:    nil,
		},
		{
			src:    []byte{8, 0, 0, 0, 2, 0, 5},
			result: UntypedBinaryRange{Lower: nil, Upper: []byte{0, 5}, LowerType: Unbounded, UpperType: Exclusive},
			err:    nil,
		},
		{
			src:    []byte{12, 0, 0, 0, 2, 0, 5},
			result: UntypedBinaryRange{Lower: nil, Upper: []byte{0, 5}, LowerType: Unbounded, UpperType: Inclusive},
			err:    nil,
		},
		{
			src:    []byte{16, 0, 0, 0, 2, 0, 4},
			result: UntypedBinaryRange{Lower: []byte{0, 4}, Upper: nil, LowerType: Exclusive, UpperType: Unbounded},
			err:    nil,
		},
		{
			src:    []byte{18, 0, 0, 0, 2, 0, 4},
			result: UntypedBinaryRange{Lower: []byte{0, 4}, Upper: nil, LowerType: Inclusive, UpperType: Unbounded},
			err:    nil,
		},
		{
			src:    []byte{24},
			result: UntypedBinaryRange{Lower: nil, Upper: nil, LowerType: Unbounded, UpperType: Unbounded},
			err:    nil,
		},
	}

	for i, tt := range tests {
		r, err := ParseUntypedBinaryRange(tt.src)
		if err != tt.err {
			t.Errorf("%d. `%v`: expected err %v, got %v", i, tt.src, tt.err, err)
			continue
		}

		if r.LowerType != tt.result.LowerType {
			t.Errorf("%d. `%v`: expected result lower type %v, got %v", i, tt.src, string(tt.result.LowerType), string(r.LowerType))
		}

		if r.UpperType != tt.result.UpperType {
			t.Errorf("%d. `%v`: expected result upper type %v, got %v", i, tt.src, string(tt.result.UpperType), string(r.UpperType))
		}

		if bytes.Compare(r.Lower, tt.result.Lower) != 0 {
			t.Errorf("%d. `%v`: expected result lower %v, got %v", i, tt.src, tt.result.Lower, r.Lower)
		}

		if bytes.Compare(r.Upper, tt.result.Upper) != 0 {
			t.Errorf("%d. `%v`: expected result upper %v, got %v", i, tt.src, tt.result.Upper, r.Upper)
		}
	}
}
