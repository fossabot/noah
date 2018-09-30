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
 */

package uint128

import (
	"bytes"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	b := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	i := FromBytes(b)

	if !bytes.Equal(i.GetBytes(), b) {
		t.Errorf("incorrect bytes representation for num: %v", i)
	}
}

func TestString(t *testing.T) {
	s := "a95e31998f38490651c02b97c7f2acca"

	i, _ := FromString(s)

	if s != i.String() {
		t.Errorf("incorrect string representation for num: %v", i)
	}
}

func TestStringTooLong(t *testing.T) {
	s := "ba95e31998f38490651c02b97c7f2acca"

	_, err := FromString(s)

	if err == nil || !strings.Contains(err.Error(), "too large") {
		t.Error("did not get error for encoding invalid uint128 string")
	}
}

func TestStringInvalidHex(t *testing.T) {
	s := "bazz95e31998849051c02b97c7f2acca"

	_, err := FromString(s)

	if err == nil || !strings.Contains(err.Error(), "could not decode") {
		t.Error("did not get error for encoding invalid uint128 string")
	}
}

func TestSub(t *testing.T) {
	testData := []struct {
		num      Uint128
		expected Uint128
		sub      uint64
	}{
		{Uint128{0, 1}, Uint128{0, 0}, 1},
		{Uint128{18446744073709551615, 18446744073709551615}, Uint128{18446744073709551615, 18446744073709551614}, 1},
		{Uint128{0, 18446744073709551615}, Uint128{0, 18446744073709551614}, 1},
		{Uint128{18446744073709551615, 0}, Uint128{18446744073709551614, 18446744073709551615}, 1},
		{Uint128{18446744073709551615, 0}, Uint128{18446744073709551614, 18446744073709551591}, 25},
	}

	for _, test := range testData {
		res := test.num.Sub(test.sub)
		if res != test.expected {
			t.Errorf("expected: %v - %d = %v but got %v", test.num, test.sub, test.expected, res)
		}
	}
}

func TestAdd(t *testing.T) {
	testData := []struct {
		num      Uint128
		expected Uint128
		add      uint64
	}{
		{Uint128{0, 0}, Uint128{0, 1}, 1},
		{Uint128{18446744073709551615, 18446744073709551614}, Uint128{18446744073709551615, 18446744073709551615}, 1},
		{Uint128{0, 18446744073709551615}, Uint128{1, 0}, 1},
		{Uint128{18446744073709551615, 0}, Uint128{18446744073709551615, 1}, 1},
		{Uint128{0, 18446744073709551615}, Uint128{1, 24}, 25},
	}

	for _, test := range testData {
		res := test.num.Add(test.add)
		if res != test.expected {
			t.Errorf("expected: %v + %d = %v but got %v", test.num, test.add, test.expected, res)
		}
	}
}

func TestEqual(t *testing.T) {
	testData := []struct {
		u1       Uint128
		u2       Uint128
		expected bool
	}{
		{Uint128{0, 0}, Uint128{0, 1}, false},
		{Uint128{1, 0}, Uint128{0, 1}, false},
		{Uint128{18446744073709551615, 18446744073709551614}, Uint128{18446744073709551615, 18446744073709551615}, false},
		{Uint128{0, 1}, Uint128{0, 1}, true},
		{Uint128{0, 0}, Uint128{0, 0}, true},
		{Uint128{314, 0}, Uint128{314, 0}, true},
		{Uint128{18446744073709551615, 18446744073709551615}, Uint128{18446744073709551615, 18446744073709551615}, true},
	}

	for _, test := range testData {

		if actual := test.u1.Equal(test.u2); actual != test.expected {
			t.Errorf("expected: %v.Equal(%v) expected %v but got %v", test.u1, test.u2, test.expected, actual)
		}
	}
}

func TestAnd(t *testing.T) {
	u1 := Uint128{14799720563850130797, 11152134164166830811}
	u2 := Uint128{10868624793753271583, 6542293553298186666}

	expected := Uint128{9529907221165552909, 1927615693132931210}
	if !(u1.And(u2)).Equal(expected) {
		t.Errorf("incorrect AND computation: %v & %v != %v", u1, u2, expected)
	}
}

func TestOr(t *testing.T) {
	u1 := Uint128{14799720563850130797, 11152134164166830811}
	u2 := Uint128{10868624793753271583, 6542293553298186666}

	expected := Uint128{16138438136437849471, 15766812024332086267}
	if !(u1.Or(u2)).Equal(expected) {
		t.Errorf("incorrect OR computation: %v | %v != %v", u1, u2, expected)
	}
}

func TestXor(t *testing.T) {
	u1 := Uint128{14799720563850130797, 11152134164166830811}
	u2 := Uint128{10868624793753271583, 6542293553298186666}

	expected := Uint128{6608530915272296562, 13839196331199155057}
	if !(u1.Xor(u2)).Equal(expected) {
		t.Errorf("incorrect XOR computation: %v ^ %v != %v", u1, u2, expected)
	}
}
