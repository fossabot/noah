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

func TestCIDRArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "cidr[]", []interface{}{
		&types.CIDRArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.CIDRArray{
			Elements: []types.CIDR{
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.CIDRArray{Status: types.Null},
		&types.CIDRArray{
			Elements: []types.CIDR{
				{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "192.168.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: types.Present},
				{Status: types.Null},
				{IPNet: mustParseCIDR(t, "255.0.0.0/8"), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.CIDRArray{
			Elements: []types.CIDR{
				{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "192.168.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestCIDRArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.CIDRArray
	}{
		{
			source: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
			result: types.CIDRArray{
				Elements:   []types.CIDR{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: (([]*net.IPNet)(nil)),
			result: types.CIDRArray{Status: types.Null},
		},
		{
			source: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
			result: types.CIDRArray{
				Elements:   []types.CIDR{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: (([]net.IP)(nil)),
			result: types.CIDRArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.CIDRArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestCIDRArrayAssignTo(t *testing.T) {
	var ipnetSlice []*net.IPNet
	var ipSlice []net.IP

	simpleTests := []struct {
		src      types.CIDRArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.CIDRArray{
				Elements:   []types.CIDR{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
		},
		{
			src: types.CIDRArray{
				Elements:   []types.CIDR{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{nil},
		},
		{
			src: types.CIDRArray{
				Elements:   []types.CIDR{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipSlice,
			expected: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
		},
		{
			src: types.CIDRArray{
				Elements:   []types.CIDR{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipSlice,
			expected: []net.IP{nil},
		},
		{
			src:      types.CIDRArray{Status: types.Null},
			dst:      &ipnetSlice,
			expected: (([]*net.IPNet)(nil)),
		},
		{
			src:      types.CIDRArray{Status: types.Null},
			dst:      &ipSlice,
			expected: (([]net.IP)(nil)),
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
