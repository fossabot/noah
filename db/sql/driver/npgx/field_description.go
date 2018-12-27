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
 */

package npgx

import (
    "github.com/readystock/noah/db/sql/types"
    "math"
    "reflect"
    "time"
)

const (
    copyData      = 'd'
    copyFail      = 'f'
    copyDone      = 'c'
    varHeaderSize = 4
)

type FieldDescription struct {
    Name            string
    Table           types.OID
    AttributeNumber uint16
    DataType        types.OID
    DataTypeSize    int16
    DataTypeName    string
    Modifier        uint32
    FormatCode      int16
}

func (fd FieldDescription) Length() (int64, bool) {
    switch fd.DataType {
    case types.TextOID, types.ByteaOID:
        return math.MaxInt64, true
    case types.VarcharOID, types.BPCharArrayOID:
        return int64(fd.Modifier - varHeaderSize), true
    default:
        return 0, false
    }
}

func (fd FieldDescription) PrecisionScale() (precision, scale int64, ok bool) {
    switch fd.DataType {
    case types.NumericOID:
        mod := fd.Modifier - varHeaderSize
        precision = int64((mod >> 16) & 0xffff)
        scale = int64(mod & 0xffff)
        return precision, scale, true
    default:
        return 0, 0, false
    }
}

func (fd FieldDescription) Type() reflect.Type {
    switch fd.DataType {
    case types.Int8OID:
        return reflect.TypeOf(int64(0))
    case types.Int4OID:
        return reflect.TypeOf(int32(0))
    case types.Int2OID:
        return reflect.TypeOf(int16(0))
    case types.VarcharOID, types.BPCharArrayOID, types.TextOID:
        return reflect.TypeOf("")
    case types.BoolOID:
        return reflect.TypeOf(false)
    case types.NumericOID:
        return reflect.TypeOf(float64(0))
    case types.DateOID, types.TimestampOID, types.TimestamptzOID:
        return reflect.TypeOf(time.Time{})
    case types.ByteaOID:
        return reflect.TypeOf([]byte(nil))
    default:
        return reflect.TypeOf(new(interface{})).Elem()
    }
}
