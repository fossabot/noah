/*
 * Copyright (c) 2019 Ready Stock
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

package types

import (
	"math"
	"strconv"

	"github.com/pkg/errors"
)

// QChar is for PostgreSQL's special 8-bit-only "char" type more akin to the C
// language's char type, or Go's byte type. (Note that the name in PostgreSQL
// itself is "char", in double-quotes, and not char.) It gets used a lot in
// PostgreSQL's system tables to hold a single ASCII character value (eg
// pg_class.relkind). It is named Qchar for quoted char to disambiguate from SQL
// standard type char.
//
// Not all possible values of QChar are representable in the text format.
// Therefore, QChar does not implement TextEncoder and TextDecoder. In
// addition, database/sql Scanner and database/sql/driver Value are not
// implemented.
type QChar struct {
	Int    int8
	Status Status
}

func (dst *QChar) Set(src interface{}) error {
	if src == nil {
		*dst = QChar{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case int8:
		*dst = QChar{Int: value, Status: Present}
	case uint8:
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case int16:
		if value < math.MinInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case uint16:
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case int32:
		if value < math.MinInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case uint32:
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case int64:
		if value < math.MinInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case uint64:
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case int:
		if value < math.MinInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case uint:
		if value > math.MaxInt8 {
			return errors.Errorf("%d is greater than maximum value for QChar", value)
		}
		*dst = QChar{Int: int8(value), Status: Present}
	case string:
		num, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		*dst = QChar{Int: int8(num), Status: Present}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to QChar", value)
	}

	return nil
}

func (dst *QChar) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Int
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *QChar) AssignTo(dst interface{}) error {
	return int64AssignTo(int64(src.Int), src.Status, dst)
}

func (dst *QChar) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = QChar{Status: Null}
		return nil
	}

	if len(src) != 1 {
		return errors.Errorf(`invalid length for "char": %v`, len(src))
	}

	*dst = QChar{Int: int8(src[0]), Status: Present}
	return nil
}

func (src *QChar) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, byte(src.Int)), nil
}

func (dst *QChar) DecodeText(ci *ConnInfo, src []byte) error {
	return dst.DecodeBinary(ci, src)
}
