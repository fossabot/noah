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

package tree

import (
    "fmt"
    "time"

    "github.com/Ready-Stock/noah/db/sql/sem/types"
    "github.com/Ready-Stock/noah/db/util/timeofday"
    "github.com/Ready-Stock/noah/db/util/timeutil"
)

// presetTypesForTesting is a mapping of qualified names to types that can be mocked out
// for tests to allow the qualified names to be type checked without throwing an error.
var presetTypesForTesting map[string]types.T

// MockNameTypes populates presetTypesForTesting for a test.
func MockNameTypes(types map[string]types.T) func() {
	presetTypesForTesting = types
	return func() {
		presetTypesForTesting = nil
	}
}

// SampleDatum is intended to be a more lightweight version of RandDatum for
// when you just need one consistent example of a datum.
func SampleDatum(t types.T) Datum {
	switch t {
	case types.Bool:
		return MakeDBool(true)
	case types.Int:
		return NewDInt(123)
	case types.Float:
		f := DFloat(123.456)
		return &f
	case types.Decimal:
		d := &DDecimal{}
		d.Decimal.SetExponent(6)
		// int64(rng.Uint64()) to get negative numbers, too
		d.Decimal.SetCoefficient(3)
		return d
	case types.String:
		return NewDString("Carl")
	case types.Bytes:
		return NewDBytes("Princess")
	case types.Date:
		return NewDDate(123123)
	case types.Time:
		return MakeDTime(timeofday.FromInt(789))
	case types.TimeTZ:
		return MakeDTimeTZ(timeofday.FromInt(789), time.UTC)
	case types.Timestamp:
		return MakeDTimestamp(timeutil.Unix(123, 123), time.Second)
	case types.TimestampTZ:
		return MakeDTimestampTZ(timeutil.Unix(123, 123), time.Second)
	case types.Interval:
		i, _ := ParseDInterval("1h1m1s")
		return i
	case types.UUID:
		u, _ := ParseDUuidFromString("3189ad07-52f2-4d60-83e8-4a8347fef718")
		return u
	case types.INet:
		i, _ := ParseDIPAddrFromINetString("127.0.0.1")
		return i
	case types.JSON:
		j, _ := ParseDJSON(`{"a": "b"}`)
		return j
	case types.Oid:
		return NewDOid(DInt(1009))
	default:
		panic(fmt.Sprintf("SampleDatum not implemented for %s", t))
	}
}
