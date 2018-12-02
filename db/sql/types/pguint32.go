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
	"database/sql/driver"
	"encoding/binary"
	"math"
	"strconv"

    "github.com/Ready-Stock/noah/db/sql/pgio"
	"github.com/pkg/errors"
)

// pguint32 is the core type that is used to implement PostgreSQL types such as
// CID and XID.
type pguint32 struct {
	Uint   uint32
	Status Status
}

// Set converts from src to dst. Note that as pguint32 is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *pguint32) Set(src interface{}) error {
	switch value := src.(type) {
	case int64:
		if value < 0 {
			return errors.Errorf("%d is less than minimum value for pguint32", value)
		}
		if value > math.MaxUint32 {
			return errors.Errorf("%d is greater than maximum value for pguint32", value)
		}
		*dst = pguint32{Uint: uint32(value), Status: Present}
	case uint32:
		*dst = pguint32{Uint: value, Status: Present}
	default:
		return errors.Errorf("cannot convert %v to pguint32", value)
	}

	return nil
}

func (dst *pguint32) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Uint
	case Null:
		return nil
	default:
		return dst.Status
	}
}

// AssignTo assigns from src to dst. Note that as pguint32 is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *pguint32) AssignTo(dst interface{}) error {
	switch v := dst.(type) {
	case *uint32:
		if src.Status == Present {
			*v = src.Uint
		} else {
			return errors.Errorf("cannot assign %v into %T", src, dst)
		}
	case **uint32:
		if src.Status == Present {
			n := src.Uint
			*v = &n
		} else {
			*v = nil
		}
	}

	return nil
}

func (dst *pguint32) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = pguint32{Status: Null}
		return nil
	}

	n, err := strconv.ParseUint(string(src), 10, 32)
	if err != nil {
		return err
	}

	*dst = pguint32{Uint: uint32(n), Status: Present}
	return nil
}

func (dst *pguint32) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = pguint32{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return errors.Errorf("invalid length: %v", len(src))
	}

	n := binary.BigEndian.Uint32(src)
	*dst = pguint32{Uint: n, Status: Present}
	return nil
}

func (src *pguint32) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, strconv.FormatUint(uint64(src.Uint), 10)...), nil
}

func (src *pguint32) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return pgio.AppendUint32(buf, src.Uint), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *pguint32) Scan(src interface{}) error {
	if src == nil {
		*dst = pguint32{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case uint32:
		*dst = pguint32{Uint: src, Status: Present}
		return nil
	case int64:
		*dst = pguint32{Uint: uint32(src), Status: Present}
		return nil
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *pguint32) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return int64(src.Uint), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
