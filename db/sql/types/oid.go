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
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package types

import (
    "database/sql/driver"
    "encoding/binary"
    "strconv"

    "github.com/pkg/errors"
    "github.com/readystock/noah/db/sql/pgio"
)

// OID (Object Identifier Type) is, according to
// https://www.postgresql.org/docs/current/static/datatype-oid.html, used
// internally by PostgreSQL as a primary key for various system tables. It is
// currently implemented as an unsigned four-byte integer. Its definition can be
// found in src/include/postgres_ext.h in the PostgreSQL sources. Because it is
// so frequently required to be in a NOT NULL condition OID cannot be NULL. To
// allow for NULL OIDs use OIDValue.
type OID uint32

func (dst *OID) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        return errors.Errorf("cannot decode nil into OID")
    }

    n, err := strconv.ParseUint(string(src), 10, 32)
    if err != nil {
        return err
    }

    *dst = OID(n)
    return nil
}

func (dst *OID) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        return errors.Errorf("cannot decode nil into OID")
    }

    if len(src) != 4 {
        return errors.Errorf("invalid length: %v", len(src))
    }

    n := binary.BigEndian.Uint32(src)
    *dst = OID(n)
    return nil
}

func (src OID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    return append(buf, strconv.FormatUint(uint64(src), 10)...), nil
}

func (src OID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    return pgio.AppendUint32(buf, uint32(src)), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *OID) Scan(src interface{}) error {
    if src == nil {
        return errors.Errorf("cannot scan NULL into %T", src)
    }

    switch src := src.(type) {
    case int64:
        *dst = OID(src)
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
func (src OID) Value() (driver.Value, error) {
    return int64(src), nil
}
