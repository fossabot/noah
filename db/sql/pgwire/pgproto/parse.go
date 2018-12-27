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

package pgproto

import (
    "bytes"
    "encoding/binary"
    "encoding/json"

    "github.com/readystock/pgx/pgio"
)

type Parse struct {
    Name          string
    Query         string
    ParameterOIDs []uint32
}

func (*Parse) Frontend() {}

func (dst *Parse) Decode(src []byte) error {
    *dst = Parse{}

    buf := bytes.NewBuffer(src)

    b, err := buf.ReadBytes(0)
    if err != nil {
        return err
    }
    dst.Name = string(b[:len(b)-1])

    b, err = buf.ReadBytes(0)
    if err != nil {
        return err
    }
    dst.Query = string(b[:len(b)-1])

    if buf.Len() < 2 {
        return &invalidMessageFormatErr{messageType: "Parse"}
    }
    parameterOIDCount := int(binary.BigEndian.Uint16(buf.Next(2)))

    for i := 0; i < parameterOIDCount; i++ {
        if buf.Len() < 4 {
            return &invalidMessageFormatErr{messageType: "Parse"}
        }
        dst.ParameterOIDs = append(dst.ParameterOIDs, binary.BigEndian.Uint32(buf.Next(4)))
    }

    return nil
}

func (src *Parse) Encode(dst []byte) []byte {
    dst = append(dst, 'P')
    sp := len(dst)
    dst = pgio.AppendInt32(dst, -1)

    dst = append(dst, src.Name...)
    dst = append(dst, 0)
    dst = append(dst, src.Query...)
    dst = append(dst, 0)

    dst = pgio.AppendUint16(dst, uint16(len(src.ParameterOIDs)))
    for _, oid := range src.ParameterOIDs {
        dst = pgio.AppendUint32(dst, oid)
    }

    pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

    return dst
}

func (src *Parse) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Type          string
        Name          string
        Query         string
        ParameterOIDs []uint32
    }{
        Type:          "Parse",
        Name:          src.Name,
        Query:         src.Query,
        ParameterOIDs: src.ParameterOIDs,
    })
}
