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

package types_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestInetTranscode(t *testing.T) {
	for _, typesName := range []string{"inet", "cidr"} {
		testutil.TestSuccessfulTranscode(t, typesName, []interface{}{
			&types.Inet{IPNet: mustParseCIDR(t, "0.0.0.0/32"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "192.168.1.0/24"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "255.0.0.0/8"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "255.255.255.255/32"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "::/128"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "::/0"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "::1/128"), Status: types.Present},
			&types.Inet{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: types.Present},
			&types.Inet{Status: types.Null},
		})
	}
}

func TestInetSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Inet
	}{
		{source: mustParseCIDR(t, "127.0.0.1/32"), result: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
		{source: mustParseCIDR(t, "127.0.0.1/32").IP, result: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
		{source: "127.0.0.1/32", result: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.Inet
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInetAssignTo(t *testing.T) {
	var ipnet net.IPNet
	var pipnet *net.IPNet
	var ip net.IP
	var pip *net.IP

	simpleTests := []struct {
		src      types.Inet
		dst      interface{}
		expected interface{}
	}{
		{src: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}, dst: &ipnet, expected: *mustParseCIDR(t, "127.0.0.1/32")},
		{src: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}, dst: &ip, expected: mustParseCIDR(t, "127.0.0.1/32").IP},
		{src: types.Inet{Status: types.Null}, dst: &pipnet, expected: ((*net.IPNet)(nil))},
		{src: types.Inet{Status: types.Null}, dst: &pip, expected: ((*net.IP)(nil))},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %#v, but result was %#v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.Inet
		dst      interface{}
		expected interface{}
	}{
		{src: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}, dst: &pipnet, expected: *mustParseCIDR(t, "127.0.0.1/32")},
		{src: types.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}, dst: &pip, expected: mustParseCIDR(t, "127.0.0.1/32").IP},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src types.Inet
		dst interface{}
	}{
		{src: types.Inet{IPNet: mustParseCIDR(t, "192.168.0.0/16"), Status: types.Present}, dst: &ip},
		{src: types.Inet{Status: types.Null}, dst: &ipnet},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
