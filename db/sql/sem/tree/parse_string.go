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
    "time"

    "github.com/Ready-Stock/noah/db/sql/coltypes"
    "github.com/Ready-Stock/noah/db/sql/pgwire/pgerror"
    "github.com/Ready-Stock/noah/db/sql/sem/types"
)

// ParseStringAs reads s as type t. If t is Bytes or String, s is returned
// unchanged. Otherwise s is parsed with the given type's Parse func.
func ParseStringAs(t types.T, s string, evalCtx *EvalContext) (Datum, error) {
	var d Datum
	var err error
	switch t {
	case types.Bytes:
		d = NewDBytes(DBytes(s))
	default:
		switch t := t.(type) {
		case types.TArray:
			typ, err := coltypes.DatumTypeToColumnType(t.Typ)
			if err != nil {
				return nil, err
			}
			d, err = ParseDArrayFromString(evalCtx, s, typ)
			if err != nil {
				return nil, err
			}
		case types.TCollatedString:
			d = NewDCollatedString(s, t.Locale, &evalCtx.collationEnv)
		default:
			d, err = parseStringAs(t, s, evalCtx)
			if d == nil && err == nil {
				return nil, pgerror.NewErrorf(pgerror.CodeInternalError, "unknown type %s (%T)", t, t)
			}
		}
	}
	return d, err
}

// ParseDatumStringAs parses s as type t. This function is guaranteed to
// round-trip when printing a Datum with FmtParseDatums.
func ParseDatumStringAs(t types.T, s string, evalCtx *EvalContext) (Datum, error) {
	switch t {
	case types.Bytes:
		return ParseDByte(s)
	default:
		return ParseStringAs(t, s, evalCtx)
	}
}

type locationContext interface {
	GetLocation() *time.Location
}

var _ locationContext = &EvalContext{}
var _ locationContext = &SemaContext{}

// parseStringAs parses s as type t for simple types. Bytes, arrays, collated
// strings are not handled. nil, nil is returned if t is not a supported type.
func parseStringAs(t types.T, s string, ctx locationContext) (Datum, error) {
	switch t {
	case types.Bool:
		return ParseDBool(s)
	case types.Date:
		return ParseDDate(s, ctx.GetLocation())
	case types.Decimal:
		return ParseDDecimal(s)
	case types.Float:
		return ParseDFloat(s)
	case types.INet:
		return ParseDIPAddrFromINetString(s)
	case types.Int:
		return ParseDInt(s)
	case types.Interval:
		return ParseDInterval(s)
	case types.JSON:
		return ParseDJSON(s)
	case types.String:
		return NewDString(s), nil
	case types.Time:
		return ParseDTime(s)
	case types.TimeTZ:
		return ParseDTimeTZ(s, ctx.GetLocation())
	case types.Timestamp:
		return ParseDTimestamp(s, time.Microsecond)
	case types.TimestampTZ:
		return ParseDTimestampTZ(s, ctx.GetLocation(), time.Microsecond)
	case types.UUID:
		return ParseDUuidFromString(s)
	default:
		return nil, nil
	}
}
