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

package pgproto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	"github.com/readystock/pgx/pgio"
)

type Bind struct {
	DestinationPortal    string
	PreparedStatement    string
	ParameterFormatCodes []int16
	Parameters           [][]byte
	ResultFormatCodes    []int16
}

func (*Bind) Frontend() {}

func (dst *Bind) Decode(src []byte) error {
	*dst = Bind{}

	idx := bytes.IndexByte(src, 0)
	if idx < 0 {
		return &invalidMessageFormatErr{messageType: "Bind"}
	}
	dst.DestinationPortal = string(src[:idx])
	rp := idx + 1

	idx = bytes.IndexByte(src[rp:], 0)
	if idx < 0 {
		return &invalidMessageFormatErr{messageType: "Bind"}
	}
	dst.PreparedStatement = string(src[rp : rp+idx])
	rp += idx + 1

	if len(src[rp:]) < 2 {
		return &invalidMessageFormatErr{messageType: "Bind"}
	}
	parameterFormatCodeCount := int(binary.BigEndian.Uint16(src[rp:]))
	rp += 2

	if parameterFormatCodeCount > 0 {
		dst.ParameterFormatCodes = make([]int16, parameterFormatCodeCount)

		if len(src[rp:]) < len(dst.ParameterFormatCodes)*2 {
			return &invalidMessageFormatErr{messageType: "Bind"}
		}
		for i := 0; i < parameterFormatCodeCount; i++ {
			dst.ParameterFormatCodes[i] = int16(binary.BigEndian.Uint16(src[rp:]))
			rp += 2
		}
	}

	if len(src[rp:]) < 2 {
		return &invalidMessageFormatErr{messageType: "Bind"}
	}
	parameterCount := int(binary.BigEndian.Uint16(src[rp:]))
	rp += 2

	if parameterCount > 0 {
		dst.Parameters = make([][]byte, parameterCount)

		for i := 0; i < parameterCount; i++ {
			if len(src[rp:]) < 4 {
				return &invalidMessageFormatErr{messageType: "Bind"}
			}

			msgSize := int(int32(binary.BigEndian.Uint32(src[rp:])))
			rp += 4

			// null
			if msgSize == -1 {
				continue
			}

			if len(src[rp:]) < msgSize {
				return &invalidMessageFormatErr{messageType: "Bind"}
			}

			dst.Parameters[i] = src[rp : rp+msgSize]
			rp += msgSize
		}
	}

	if len(src[rp:]) < 2 {
		return &invalidMessageFormatErr{messageType: "Bind"}
	}
	resultFormatCodeCount := int(binary.BigEndian.Uint16(src[rp:]))
	rp += 2

	dst.ResultFormatCodes = make([]int16, resultFormatCodeCount)
	if len(src[rp:]) < len(dst.ResultFormatCodes)*2 {
		return &invalidMessageFormatErr{messageType: "Bind"}
	}
	for i := 0; i < resultFormatCodeCount; i++ {
		dst.ResultFormatCodes[i] = int16(binary.BigEndian.Uint16(src[rp:]))
		rp += 2
	}

	return nil
}

func (src *Bind) Encode(dst []byte) []byte {
	dst = append(dst, 'B')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.DestinationPortal...)
	dst = append(dst, 0)
	dst = append(dst, src.PreparedStatement...)
	dst = append(dst, 0)

	dst = pgio.AppendUint16(dst, uint16(len(src.ParameterFormatCodes)))
	for _, fc := range src.ParameterFormatCodes {
		dst = pgio.AppendInt16(dst, fc)
	}

	dst = pgio.AppendUint16(dst, uint16(len(src.Parameters)))
	for _, p := range src.Parameters {
		if p == nil {
			dst = pgio.AppendInt32(dst, -1)
			continue
		}

		dst = pgio.AppendInt32(dst, int32(len(p)))
		dst = append(dst, p...)
	}

	dst = pgio.AppendUint16(dst, uint16(len(src.ResultFormatCodes)))
	for _, fc := range src.ResultFormatCodes {
		dst = pgio.AppendInt16(dst, fc)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *Bind) MarshalJSON() ([]byte, error) {
	formattedParameters := make([]map[string]string, len(src.Parameters))
	for i, p := range src.Parameters {
		if p == nil {
			continue
		}

		if src.ParameterFormatCodes[i] == 0 {
			formattedParameters[i] = map[string]string{"text": string(p)}
		} else {
			formattedParameters[i] = map[string]string{"binary": hex.EncodeToString(p)}
		}
	}

	return json.Marshal(struct {
		Type                 string
		DestinationPortal    string
		PreparedStatement    string
		ParameterFormatCodes []int16
		Parameters           []map[string]string
		ResultFormatCodes    []int16
	}{
		Type:                 "Bind",
		DestinationPortal:    src.DestinationPortal,
		PreparedStatement:    src.PreparedStatement,
		ParameterFormatCodes: src.ParameterFormatCodes,
		Parameters:           formattedParameters,
		ResultFormatCodes:    src.ResultFormatCodes,
	})
}
