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

package types

import (
	"database/sql/driver"
    "github.com/Ready-Stock/noah/db/sql/pgio"
	"github.com/pkg/errors"
)

type Daterange struct {
	Lower     Date
	Upper     Date
	LowerType BoundType
	UpperType BoundType
	Status    Status
}

func (dst *Daterange) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Daterange", src)
}

func (dst *Daterange) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Daterange) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Daterange) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Daterange{Status: Null}
		return nil
	}

	utr, err := ParseUntypedTextRange(string(src))
	if err != nil {
		return err
	}

	*dst = Daterange{Status: Present}

	dst.LowerType = utr.LowerType
	dst.UpperType = utr.UpperType

	if dst.LowerType == Empty {
		return nil
	}

	if dst.LowerType == Inclusive || dst.LowerType == Exclusive {
		if err := dst.Lower.DecodeText(ci, []byte(utr.Lower)); err != nil {
			return err
		}
	}

	if dst.UpperType == Inclusive || dst.UpperType == Exclusive {
		if err := dst.Upper.DecodeText(ci, []byte(utr.Upper)); err != nil {
			return err
		}
	}

	return nil
}

func (dst *Daterange) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Daterange{Status: Null}
		return nil
	}

	ubr, err := ParseUntypedBinaryRange(src)
	if err != nil {
		return err
	}

	*dst = Daterange{Status: Present}

	dst.LowerType = ubr.LowerType
	dst.UpperType = ubr.UpperType

	if dst.LowerType == Empty {
		return nil
	}

	if dst.LowerType == Inclusive || dst.LowerType == Exclusive {
		if err := dst.Lower.DecodeBinary(ci, ubr.Lower); err != nil {
			return err
		}
	}

	if dst.UpperType == Inclusive || dst.UpperType == Exclusive {
		if err := dst.Upper.DecodeBinary(ci, ubr.Upper); err != nil {
			return err
		}
	}

	return nil
}

func (src Daterange) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	switch src.LowerType {
	case Exclusive, Unbounded:
		buf = append(buf, '(')
	case Inclusive:
		buf = append(buf, '[')
	case Empty:
		return append(buf, "empty"...), nil
	default:
		return nil, errors.Errorf("unknown lower bound type %v", src.LowerType)
	}

	var err error

	if src.LowerType != Unbounded {
		buf, err = src.Lower.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		} else if buf == nil {
			return nil, errors.Errorf("Lower cannot be null unless LowerType is Unbounded")
		}
	}

	buf = append(buf, ',')

	if src.UpperType != Unbounded {
		buf, err = src.Upper.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		} else if buf == nil {
			return nil, errors.Errorf("Upper cannot be null unless UpperType is Unbounded")
		}
	}

	switch src.UpperType {
	case Exclusive, Unbounded:
		buf = append(buf, ')')
	case Inclusive:
		buf = append(buf, ']')
	default:
		return nil, errors.Errorf("unknown upper bound type %v", src.UpperType)
	}

	return buf, nil
}

func (src Daterange) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var rangeType byte
	switch src.LowerType {
	case Inclusive:
		rangeType |= lowerInclusiveMask
	case Unbounded:
		rangeType |= lowerUnboundedMask
	case Exclusive:
	case Empty:
		return append(buf, emptyMask), nil
	default:
		return nil, errors.Errorf("unknown LowerType: %v", src.LowerType)
	}

	switch src.UpperType {
	case Inclusive:
		rangeType |= upperInclusiveMask
	case Unbounded:
		rangeType |= upperUnboundedMask
	case Exclusive:
	default:
		return nil, errors.Errorf("unknown UpperType: %v", src.UpperType)
	}

	buf = append(buf, rangeType)

	var err error

	if src.LowerType != Unbounded {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		buf, err = src.Lower.EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, errors.Errorf("Lower cannot be null unless LowerType is Unbounded")
		}

		pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
	}

	if src.UpperType != Unbounded {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		buf, err = src.Upper.EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, errors.Errorf("Upper cannot be null unless UpperType is Unbounded")
		}

		pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Daterange) Scan(src interface{}) error {
	if src == nil {
		*dst = Daterange{Status: Null}
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
func (src Daterange) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
