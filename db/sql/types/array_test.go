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

package types_test

import (
	"reflect"
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
)

func TestParseUntypedTextArray(t *testing.T) {
	tests := []struct {
		source string
		result types.UntypedTextArray
	}{
		{
			source: "{}",
			result: types.UntypedTextArray{
				Elements:   nil,
				Dimensions: nil,
			},
		},
		{
			source: "{1}",
			result: types.UntypedTextArray{
				Elements:   []string{"1"},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: "{a,b}",
			result: types.UntypedTextArray{
				Elements:   []string{"a", "b"},
				Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			},
		},
		{
			source: `{"NULL"}`,
			result: types.UntypedTextArray{
				Elements:   []string{"NULL"},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: `{""}`,
			result: types.UntypedTextArray{
				Elements:   []string{""},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: `{"He said, \"Hello.\""}`,
			result: types.UntypedTextArray{
				Elements:   []string{`He said, "Hello."`},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: "{{a,b},{c,d},{e,f}}",
			result: types.UntypedTextArray{
				Elements:   []string{"a", "b", "c", "d", "e", "f"},
				Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			},
		},
		{
			source: "{{{a,b},{c,d},{e,f}},{{a,b},{c,d},{e,f}}}",
			result: types.UntypedTextArray{
				Elements: []string{"a", "b", "c", "d", "e", "f", "a", "b", "c", "d", "e", "f"},
				Dimensions: []types.ArrayDimension{
					{Length: 2, LowerBound: 1},
					{Length: 3, LowerBound: 1},
					{Length: 2, LowerBound: 1},
				},
			},
		},
		{
			source: "[4:4]={1}",
			result: types.UntypedTextArray{
				Elements:   []string{"1"},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 4}},
			},
		},
		{
			source: "[4:5][2:3]={{a,b},{c,d}}",
			result: types.UntypedTextArray{
				Elements: []string{"a", "b", "c", "d"},
				Dimensions: []types.ArrayDimension{
					{Length: 2, LowerBound: 4},
					{Length: 2, LowerBound: 2},
				},
			},
		},
	}

	for i, tt := range tests {
		r, err := types.ParseUntypedTextArray(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(*r, tt.result) {
			t.Errorf("%d: expected %+v to be parsed to %+v, but it was %+v", i, tt.source, tt.result, *r)
		}
	}
}
