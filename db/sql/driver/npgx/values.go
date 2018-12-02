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

package npgx

import (
	"database/sql/driver"
	"fmt"
    "github.com/Ready-Stock/noah/db/sql/pgio"
    "github.com/Ready-Stock/noah/db/sql/types"
	"github.com/Ready-Stock/pgx/pgtype"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"time"
)

// PostgreSQL format codes
const (
	TextFormatCode   = 0
	BinaryFormatCode = 1
)

// SerializationError occurs on failure to encode or decode a value
type SerializationError string

func (e SerializationError) Error() string {
	return string(e)
}

func convertSimpleArgument(ci *types.ConnInfo, arg interface{}) (interface{}, error) {
	if arg == nil {
		return nil, nil
	}

	switch arg := arg.(type) {

	// https://github.com/jackc/pgx/issues/409 Changed JSON and JSONB to surface
	// []byte to database/sql instead of string. But that caused problems with the
	// simple protocol because the driver.Valuer case got taken before the
	// pgtype.TextEncoder case. And driver.Valuer needed to be first in the usual
	// case because of https://github.com/jackc/pgx/issues/339. So instead we
	// special case JSON and JSONB.
	case *types.JSON:
		buf, err := arg.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil
	case *types.JSONB:
		buf, err := arg.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil

	case driver.Valuer:
		return callValuerValue(arg)
	case types.TextEncoder:
		buf, err := arg.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, nil
		}
		return string(buf), nil
	case int64:
		return arg, nil
	case float64:
		return arg, nil
	case bool:
		return arg, nil
	case time.Time:
		return arg, nil
	case string:
		return arg, nil
	case []byte:
		return arg, nil
	case int8:
		return int64(arg), nil
	case int16:
		return int64(arg), nil
	case int32:
		return int64(arg), nil
	case int:
		return int64(arg), nil
	case uint8:
		return int64(arg), nil
	case uint16:
		return int64(arg), nil
	case uint32:
		return int64(arg), nil
	case uint64:
		if arg > math.MaxInt64 {
			return nil, errors.Errorf("arg too big for int64: %v", arg)
		}
		return int64(arg), nil
	case uint:
		if uint64(arg) > math.MaxInt64 {
			return nil, errors.Errorf("arg too big for int64: %v", arg)
		}
		return int64(arg), nil
	case float32:
		return float64(arg), nil
	}

	refVal := reflect.ValueOf(arg)

	if refVal.Kind() == reflect.Ptr {
		if refVal.IsNil() {
			return nil, nil
		}
		arg = refVal.Elem().Interface()
		return convertSimpleArgument(ci, arg)
	}

	if strippedArg, ok := stripNamedType(&refVal); ok {
		return convertSimpleArgument(ci, strippedArg)
	}
	return nil, SerializationError(fmt.Sprintf("Cannot encode %T in simple protocol - %T must implement driver.Valuer, pgtype.TextEncoder, or be a native type", arg, arg))
}

func encodePreparedStatementArgument(ci *types.ConnInfo, buf []byte, oid types.OID, arg interface{}) ([]byte, error) {
	if arg == nil {
		return pgio.AppendInt32(buf, -1), nil
	}

	switch arg := arg.(type) {
	case types.BinaryEncoder:
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)
		argBuf, err := arg.EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if argBuf != nil {
			buf = argBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
		return buf, nil
	case types.TextEncoder:
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)
		argBuf, err := arg.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		}
		if argBuf != nil {
			buf = argBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
		return buf, nil
	case string:
		buf = pgio.AppendInt32(buf, int32(len(arg)))
		buf = append(buf, arg...)
		return buf, nil
	}

	refVal := reflect.ValueOf(arg)

	if refVal.Kind() == reflect.Ptr {
		if refVal.IsNil() {
			return pgio.AppendInt32(buf, -1), nil
		}
		arg = refVal.Elem().Interface()
		return encodePreparedStatementArgument(ci, buf, oid, arg)
	}

	if dt, ok := ci.DataTypeForOID(oid); ok {
		value := dt.Value
		err := value.Set(arg)
		if err != nil {
			{
				if arg, ok := arg.(driver.Valuer); ok {
					v, err := callValuerValue(arg)
					if err != nil {
						return nil, err
					}
					return encodePreparedStatementArgument(ci, buf, oid, v)
				}
			}

			return nil, err
		}

		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)
		argBuf, err := value.(types.BinaryEncoder).EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if argBuf != nil {
			buf = argBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
		return buf, nil
	}

	if arg, ok := arg.(driver.Valuer); ok {
		v, err := callValuerValue(arg)
		if err != nil {
			return nil, err
		}
		return encodePreparedStatementArgument(ci, buf, oid, v)
	}

	if strippedArg, ok := stripNamedType(&refVal); ok {
		return encodePreparedStatementArgument(ci, buf, oid, strippedArg)
	}
	return nil, SerializationError(fmt.Sprintf("Cannot encode %T into oid %v - %T must implement Encoder or be converted to a string", arg, oid, arg))
}

// chooseParameterFormatCode determines the correct format code for an
// argument to a prepared statement. It defaults to TextFormatCode if no
// determination can be made.
func chooseParameterFormatCode(ci *types.ConnInfo, oid types.OID, arg interface{}) int16 {
	switch arg.(type) {
	case types.BinaryEncoder:
		return BinaryFormatCode
	case string, *string, types.TextEncoder:
		return TextFormatCode
	}

	if dt, ok := ci.DataTypeForOID(oid); ok {
		if _, ok := dt.Value.(pgtype.BinaryEncoder); ok {
			if arg, ok := arg.(driver.Valuer); ok {
				if err := dt.Value.Set(arg); err != nil {
					if value, err := callValuerValue(arg); err == nil {
						if _, ok := value.(string); ok {
							return TextFormatCode
						}
					}
				}
			}

			return BinaryFormatCode
		}
	}

	return TextFormatCode
}

func stripNamedType(val *reflect.Value) (interface{}, bool) {
	switch val.Kind() {
	case reflect.Int:
		convVal := int(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int8:
		convVal := int8(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int16:
		convVal := int16(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int32:
		convVal := int32(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Int64:
		convVal := int64(val.Int())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint:
		convVal := uint(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint8:
		convVal := uint8(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint16:
		convVal := uint16(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint32:
		convVal := uint32(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.Uint64:
		convVal := uint64(val.Uint())
		return convVal, reflect.TypeOf(convVal) != val.Type()
	case reflect.String:
		convVal := val.String()
		return convVal, reflect.TypeOf(convVal) != val.Type()
	}

	return nil, false
}
