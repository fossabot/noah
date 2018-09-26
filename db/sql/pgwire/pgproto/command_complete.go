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

package pgproto

import (
	"bytes"
	"encoding/json"

	"github.com/Ready-Stock/pgx/pgio"
)

type CommandComplete struct {
	CommandTag string
}

func (*CommandComplete) Backend() {}

func (dst *CommandComplete) Decode(src []byte) error {
	idx := bytes.IndexByte(src, 0)
	if idx != len(src)-1 {
		return &invalidMessageFormatErr{messageType: "CommandComplete"}
	}

	dst.CommandTag = string(src[:idx])

	return nil
}

func (src *CommandComplete) Encode(dst []byte) []byte {
	dst = append(dst, 'C')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.CommandTag...)
	dst = append(dst, 0)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *CommandComplete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string
		CommandTag string
	}{
		Type:       "CommandComplete",
		CommandTag: src.CommandTag,
	})
}
