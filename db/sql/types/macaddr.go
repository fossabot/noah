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
	"net"

	"github.com/pkg/errors"
)

type Macaddr struct {
	Addr   net.HardwareAddr
	Status Status
}

func (dst *Macaddr) Set(src interface{}) error {
	if src == nil {
		*dst = Macaddr{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case net.HardwareAddr:
		addr := make(net.HardwareAddr, len(value))
		copy(addr, value)
		*dst = Macaddr{Addr: addr, Status: Present}
	case string:
		addr, err := net.ParseMAC(value)
		if err != nil {
			return err
		}
		*dst = Macaddr{Addr: addr, Status: Present}
	default:
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Macaddr", value)
	}

	return nil
}

func (dst *Macaddr) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Addr
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Macaddr) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *net.HardwareAddr:
			*v = make(net.HardwareAddr, len(src.Addr))
			copy(*v, src.Addr)
			return nil
		case *string:
			*v = src.Addr.String()
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

func (dst *Macaddr) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Macaddr{Status: Null}
		return nil
	}

	addr, err := net.ParseMAC(string(src))
	if err != nil {
		return err
	}

	*dst = Macaddr{Addr: addr, Status: Present}
	return nil
}

func (dst *Macaddr) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Macaddr{Status: Null}
		return nil
	}

	if len(src) != 6 {
		return errors.Errorf("Received an invalid size for a macaddr: %d", len(src))
	}

	addr := make(net.HardwareAddr, 6)
	copy(addr, src)

	*dst = Macaddr{Addr: addr, Status: Present}

	return nil
}

func (src *Macaddr) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.Addr.String()...), nil
}

// EncodeBinary encodes src into w.
func (src *Macaddr) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.Addr...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Macaddr) Scan(src interface{}) error {
	if src == nil {
		*dst = Macaddr{Status: Null}
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
func (src *Macaddr) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
