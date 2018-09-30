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

package uuid

import (
	"database/sql/driver"

	"github.com/pkg/errors"

	"github.com/Ready-Stock/Noah/db/sql/types"
	uuid "github.com/satori/go.uuid"
)

var errUndefined = errors.New("cannot encode status undefined")

type UUID struct {
	UUID   uuid.UUID
	Status types.Status
}

func (dst *UUID) Set(src interface{}) error {
	switch value := src.(type) {
	case uuid.UUID:
		*dst = UUID{UUID: value, Status: types.Present}
	case [16]byte:
		*dst = UUID{UUID: uuid.UUID(value), Status: types.Present}
	case []byte:
		if len(value) != 16 {
			return errors.Errorf("[]byte must be 16 bytes to convert to UUID: %d", len(value))
		}
		*dst = UUID{Status: types.Present}
		copy(dst.UUID[:], value)
	case string:
		uuid, err := uuid.FromString(value)
		if err != nil {
			return err
		}
		*dst = UUID{UUID: uuid, Status: types.Present}
	default:
		// If all else fails see if types.UUID can handle it. If so, translate through that.
		pgUUID := &types.UUID{}
		if err := pgUUID.Set(value); err != nil {
			return errors.Errorf("cannot convert %v to UUID", value)
		}

		*dst = UUID{UUID: uuid.UUID(pgUUID.Bytes), Status: pgUUID.Status}
	}

	return nil
}

func (dst *UUID) Get() interface{} {
	switch dst.Status {
	case types.Present:
		return dst.UUID
	case types.Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *UUID) AssignTo(dst interface{}) error {
	switch src.Status {
	case types.Present:
		switch v := dst.(type) {
		case *uuid.UUID:
			*v = src.UUID
		case *[16]byte:
			*v = [16]byte(src.UUID)
			return nil
		case *[]byte:
			*v = make([]byte, 16)
			copy(*v, src.UUID[:])
			return nil
		case *string:
			*v = src.UUID.String()
			return nil
		default:
			if nextDst, retry := types.GetAssignToDstType(v); retry {
				return src.AssignTo(nextDst)
			}
		}
	case types.Null:
		return types.NullAssignTo(dst)
	}

	return errors.Errorf("cannot assign %v into %T", src, dst)
}

func (dst *UUID) DecodeText(ci *types.ConnInfo, src []byte) error {
	if src == nil {
		*dst = UUID{Status: types.Null}
		return nil
	}

	u, err := uuid.FromString(string(src))
	if err != nil {
		return err
	}

	*dst = UUID{UUID: u, Status: types.Present}
	return nil
}

func (dst *UUID) DecodeBinary(ci *types.ConnInfo, src []byte) error {
	if src == nil {
		*dst = UUID{Status: types.Null}
		return nil
	}

	if len(src) != 16 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}

	*dst = UUID{Status: types.Present}
	copy(dst.UUID[:], src)
	return nil
}

func (src *UUID) EncodeText(ci *types.ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case types.Null:
		return nil, nil
	case types.Undefined:
		return nil, errUndefined
	}

	return append(buf, src.UUID.String()...), nil
}

func (src *UUID) EncodeBinary(ci *types.ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case types.Null:
		return nil, nil
	case types.Undefined:
		return nil, errUndefined
	}

	return append(buf, src.UUID[:]...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *UUID) Scan(src interface{}) error {
	if src == nil {
		*dst = UUID{Status: types.Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		return dst.DecodeText(nil, src)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *UUID) Value() (driver.Value, error) {
	return types.EncodeValueText(src)
}
