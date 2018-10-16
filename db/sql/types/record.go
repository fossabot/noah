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

package types

import (
	"encoding/binary"
	"reflect"

	"github.com/pkg/errors"
)

// Record is the generic PostgreSQL record type such as is created with the
// "row" function. Record only implements BinaryEncoder and Value. The text
// format output format from PostgreSQL does not include type information and is
// therefore impossible to decode. No encoders are implemented because
// PostgreSQL does not support input of generic records.
type Record struct {
	Fields []Value
	Status Status
}

func (dst *Record) Set(src interface{}) error {
	if src == nil {
		*dst = Record{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case []Value:
		*dst = Record{Fields: value, Status: Present}
	default:
		return errors.Errorf("cannot convert %v to Record", src)
	}

	return nil
}

func (dst *Record) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Fields
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Record) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *[]Value:
			*v = make([]Value, len(src.Fields))
			copy(*v, src.Fields)
			return nil
		case *[]interface{}:
			*v = make([]interface{}, len(src.Fields))
			for i := range *v {
				(*v)[i] = src.Fields[i].Get()
			}
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
		}
	case Null:
		return NullAssignTo(dst)
	}

	return errors.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *Record) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Record{Status: Null}
		return nil
	}

	rp := 0

	if len(src[rp:]) < 4 {
		return errors.Errorf("Record incomplete %v", src)
	}
	fieldCount := int(int32(binary.BigEndian.Uint32(src[rp:])))
	rp += 4

	fields := make([]Value, fieldCount)

	for i := 0; i < fieldCount; i++ {
		if len(src[rp:]) < 8 {
			return errors.Errorf("Record incomplete %v", src)
		}
		fieldOID := OID(binary.BigEndian.Uint32(src[rp:]))
		rp += 4

		fieldLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4

		var binaryDecoder BinaryDecoder
		if dt, ok := ci.DataTypeForOID(fieldOID); ok {
			binaryDecoder, _ = dt.Value.(BinaryDecoder)
		}
		if binaryDecoder == nil {
			return errors.Errorf("unknown oid while decoding record: %v", fieldOID)
		}

		var fieldBytes []byte
		if fieldLen >= 0 {
			if len(src[rp:]) < fieldLen {
				return errors.Errorf("Record incomplete %v", src)
			}
			fieldBytes = src[rp : rp+fieldLen]
			rp += fieldLen
		}

		// Duplicate struct to scan into
		binaryDecoder = reflect.New(reflect.ValueOf(binaryDecoder).Elem().Type()).Interface().(BinaryDecoder)

		if err := binaryDecoder.DecodeBinary(ci, fieldBytes); err != nil {
			return err
		}

		fields[i] = binaryDecoder.(Value)
	}

	*dst = Record{Fields: fields, Status: Present}

	return nil
}
