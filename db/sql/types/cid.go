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
)

// CID is PostgreSQL's Command Identifier type.
//
// When one does
//
//     select cmin, cmax, * from some_table;
//
// it is the data type of the cmin and cmax hidden system columns.
//
// It is currently implemented as an unsigned four byte integer.
// Its definition can be found in src/include/c.h as CommandId
// in the PostgreSQL sources.
type CID pguint32

// Set converts from src to dst. Note that as CID is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *CID) Set(src interface{}) error {
    return (*pguint32)(dst).Set(src)
}

func (dst *CID) Get() interface{} {
    return (*pguint32)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as CID is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *CID) AssignTo(dst interface{}) error {
    return (*pguint32)(src).AssignTo(dst)
}

func (dst *CID) DecodeText(ci *ConnInfo, src []byte) error {
    return (*pguint32)(dst).DecodeText(ci, src)
}

func (dst *CID) DecodeBinary(ci *ConnInfo, src []byte) error {
    return (*pguint32)(dst).DecodeBinary(ci, src)
}

func (src *CID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    return (*pguint32)(src).EncodeText(ci, buf)
}

func (src *CID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    return (*pguint32)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *CID) Scan(src interface{}) error {
    return (*pguint32)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *CID) Value() (driver.Value, error) {
    return (*pguint32)(src).Value()
}
