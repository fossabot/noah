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

type Decimal Numeric

func (dst *Decimal) Set(src interface{}) error {
	return (*Numeric)(dst).Set(src)
}

func (dst *Decimal) Get() interface{} {
	return (*Numeric)(dst).Get()
}

func (src *Decimal) AssignTo(dst interface{}) error {
	return (*Numeric)(src).AssignTo(dst)
}

func (dst *Decimal) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Numeric)(dst).DecodeText(ci, src)
}

func (dst *Decimal) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Numeric)(dst).DecodeBinary(ci, src)
}

func (src *Decimal) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Numeric)(src).EncodeText(ci, buf)
}

func (src *Decimal) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Numeric)(src).EncodeBinary(ci, buf)
}
