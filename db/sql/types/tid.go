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
	"fmt"
	"strconv"
	"strings"

    "github.com/Ready-Stock/noah/db/sql/pgio"
	"github.com/pkg/errors"
)

// TID is PostgreSQL's Tuple Identifier type.
//
// When one does
//
// 	select ctid, * from some_table;
//
// it is the data type of the ctid hidden system column.
//
// It is currently implemented as a pair unsigned two byte integers.
// Its conversion functions can be found in src/backend/utils/adt/tid.c
// in the PostgreSQL sources.
type TID struct {
	BlockNumber  uint32
	OffsetNumber uint16
	Status       Status
}

func (dst *TID) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to TID", src)
}

func (dst *TID) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *TID) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *TID) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = TID{Status: Null}
		return nil
	}

	if len(src) < 5 {
		return errors.Errorf("invalid length for tid: %v", len(src))
	}

	parts := strings.SplitN(string(src[1:len(src)-1]), ",", 2)
	if len(parts) < 2 {
		return errors.Errorf("invalid format for tid")
	}

	blockNumber, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return err
	}

	offsetNumber, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return err
	}

	*dst = TID{BlockNumber: uint32(blockNumber), OffsetNumber: uint16(offsetNumber), Status: Present}
	return nil
}

func (dst *TID) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = TID{Status: Null}
		return nil
	}

	if len(src) != 6 {
		return errors.Errorf("invalid length for tid: %v", len(src))
	}

	*dst = TID{
		BlockNumber:  binary.BigEndian.Uint32(src),
		OffsetNumber: binary.BigEndian.Uint16(src[4:]),
		Status:       Present,
	}
	return nil
}

func (src *TID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`(%d,%d)`, src.BlockNumber, src.OffsetNumber)...)
	return buf, nil
}

func (src *TID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint32(buf, src.BlockNumber)
	buf = pgio.AppendUint16(buf, src.OffsetNumber)
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *TID) Scan(src interface{}) error {
	if src == nil {
		*dst = TID{Status: Null}
		return nil
	}

	switch src := src.(type) {
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
func (src *TID) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
