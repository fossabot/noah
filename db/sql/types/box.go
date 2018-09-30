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

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"math"
	"strconv"
	"strings"
	"github.com/pkg/errors"
)

type Box struct {
	P      [2]Vec2
	Status Status
}

func (dst *Box) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Box", src)
}

func (dst *Box) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Box) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Box) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Box{Status: Null}
		return nil
	}

	if len(src) < 11 {
		return errors.Errorf("invalid length for Box: %v", len(src))
	}

	str := string(src[1:])

	var end int
	end = strings.IndexByte(str, ',')

	x1, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+1:]
	end = strings.IndexByte(str, ')')

	y1, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+3:]
	end = strings.IndexByte(str, ',')

	x2, err := strconv.ParseFloat(str[:end], 64)
	if err != nil {
		return err
	}

	str = str[end+1 : len(str)-1]

	y2, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*dst = Box{P: [2]Vec2{{x1, y1}, {x2, y2}}, Status: Present}
	return nil
}

func (dst *Box) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Box{Status: Null}
		return nil
	}

	if len(src) != 32 {
		return errors.Errorf("invalid length for Box: %v", len(src))
	}

	x1 := binary.BigEndian.Uint64(src)
	y1 := binary.BigEndian.Uint64(src[8:])
	x2 := binary.BigEndian.Uint64(src[16:])
	y2 := binary.BigEndian.Uint64(src[24:])

	*dst = Box{
		P: [2]Vec2{
			{math.Float64frombits(x1), math.Float64frombits(y1)},
			{math.Float64frombits(x2), math.Float64frombits(y2)},
		},
		Status: Present,
	}
	return nil
}

func (src *Box) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`(%s,%s),(%s,%s)`,
		strconv.FormatFloat(src.P[0].X, 'f', -1, 64),
		strconv.FormatFloat(src.P[0].Y, 'f', -1, 64),
		strconv.FormatFloat(src.P[1].X, 'f', -1, 64),
		strconv.FormatFloat(src.P[1].Y, 'f', -1, 64),
	)...)
	return buf, nil
}

func (src *Box) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[0].X))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[0].Y))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[1].X))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.P[1].Y))

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Box) Scan(src interface{}) error {
	if src == nil {
		*dst = Box{Status: Null}
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
func (src *Box) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
