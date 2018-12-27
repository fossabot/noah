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

const (
    TextFormat   = 0
    BinaryFormat = 1
)

type FieldDescription struct {
    Name                 string
    TableOID             uint32
    TableAttributeNumber uint16
    DataTypeOID          uint32
    DataTypeSize         int16
    TypeModifier         uint32
    Format               int16
}

type RowDescription struct {
    Fields []FieldDescription
}

func (*RowDescription) Backend() {}

func (dst *RowDescription) Decode(src []byte) error {
    buf := bytes.NewBuffer(src)

    if buf.Len() < 2 {
        return &invalidMessageFormatErr{messageType: "RowDescription"}
    }
    fieldCount := int(binary.BigEndian.Uint16(buf.Next(2)))

    *dst = RowDescription{Fields: make([]FieldDescription, fieldCount)}

    for i := 0; i < fieldCount; i++ {
        var fd FieldDescription
        bName, err := buf.ReadBytes(0)
        if err != nil {
            return err
        }
        fd.Name = string(bName[:len(bName)-1])

        // Since buf.Next() doesn't return an error if we hit the end of the buffer
        // check Len ahead of time
        if buf.Len() < 18 {
            return &invalidMessageFormatErr{messageType: "RowDescription"}
        }

        fd.TableOID = binary.BigEndian.Uint32(buf.Next(4))
        fd.TableAttributeNumber = binary.BigEndian.Uint16(buf.Next(2))
        fd.DataTypeOID = binary.BigEndian.Uint32(buf.Next(4))
        fd.DataTypeSize = int16(binary.BigEndian.Uint16(buf.Next(2)))
        fd.TypeModifier = binary.BigEndian.Uint32(buf.Next(4))
        fd.Format = int16(binary.BigEndian.Uint16(buf.Next(2)))

        dst.Fields[i] = fd
    }

    return nil
}

func (src *RowDescription) Encode(dst []byte) []byte {
    dst = append(dst, 'T')
    sp := len(dst)
    dst = pgio.AppendInt32(dst, -1)

    dst = pgio.AppendUint16(dst, uint16(len(src.Fields)))
    for _, fd := range src.Fields {
        dst = append(dst, fd.Name...)
        dst = append(dst, 0)

        dst = pgio.AppendUint32(dst, fd.TableOID)
        dst = pgio.AppendUint16(dst, fd.TableAttributeNumber)
        dst = pgio.AppendUint32(dst, fd.DataTypeOID)
        dst = pgio.AppendInt16(dst, fd.DataTypeSize)
        dst = pgio.AppendUint32(dst, fd.TypeModifier)
        dst = pgio.AppendInt16(dst, fd.Format)
    }

    pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

    return dst
}

func (src *RowDescription) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Type   string
        Fields []FieldDescription
    }{
        Type:   "RowDescription",
        Fields: src.Fields,
    })
}
