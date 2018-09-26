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
 */

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/pkg/errors"
)

type Polygon struct {
	P      []Vec2
	Status Status
}

func (dst *Polygon) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Polygon", src)
}

func (dst *Polygon) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Polygon) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Polygon) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Polygon{Status: Null}
		return nil
	}

	if len(src) < 7 {
		return errors.Errorf("invalid length for Polygon: %v", len(src))
	}

	points := make([]Vec2, 0)

	str := string(src[2:])

	for {
		end := strings.IndexByte(str, ',')
		x, err := strconv.ParseFloat(str[:end], 64)
		if err != nil {
			return err
		}

		str = str[end+1:]
		end = strings.IndexByte(str, ')')

		y, err := strconv.ParseFloat(str[:end], 64)
		if err != nil {
			return err
		}

		points = append(points, Vec2{x, y})

		if end+3 < len(str) {
			str = str[end+3:]
		} else {
			break
		}
	}

	*dst = Polygon{P: points, Status: Present}
	return nil
}

func (dst *Polygon) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Polygon{Status: Null}
		return nil
	}

	if len(src) < 5 {
		return errors.Errorf("invalid length for Polygon: %v", len(src))
	}

	pointCount := int(binary.BigEndian.Uint32(src))
	rp := 4

	if 4+pointCount*16 != len(src) {
		return errors.Errorf("invalid length for Polygon with %d points: %v", pointCount, len(src))
	}

	points := make([]Vec2, pointCount)
	for i := 0; i < len(points); i++ {
		x := binary.BigEndian.Uint64(src[rp:])
		rp += 8
		y := binary.BigEndian.Uint64(src[rp:])
		rp += 8
		points[i] = Vec2{math.Float64frombits(x), math.Float64frombits(y)}
	}

	*dst = Polygon{
		P:      points,
		Status: Present,
	}
	return nil
}

func (src *Polygon) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, '(')

	for i, p := range src.P {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, fmt.Sprintf(`(%s,%s)`,
			strconv.FormatFloat(p.X, 'f', -1, 64),
			strconv.FormatFloat(p.Y, 'f', -1, 64),
		)...)
	}

	return append(buf, ')'), nil
}

func (src *Polygon) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendInt32(buf, int32(len(src.P)))

	for _, p := range src.P {
		buf = pgio.AppendUint64(buf, math.Float64bits(p.X))
		buf = pgio.AppendUint64(buf, math.Float64bits(p.Y))
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Polygon) Scan(src interface{}) error {
	if src == nil {
		*dst = Polygon{Status: Null}
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
func (src *Polygon) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
