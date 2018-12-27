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

// This code was derived from https://github.com/youtube/vitess.

package stringencoding

import (
    "bytes"
    "fmt"
    "unicode/utf8"
)

// This is its own package so it can be shared among packages that parser
// depends on.

var (
    // DontEscape is a sentinel value for characters that don't need to be escaped.
    DontEscape = byte(255)
    // EncodeMap specifies how to escape binary data with '\'.
    EncodeMap [256]byte
    // HexMap is a mapping from each byte to the `\x%%` hex form as a []byte.
    HexMap [256][]byte
    // RawHexMap is a mapping from each byte to the `%%` hex form as a []byte.
    RawHexMap [256][]byte
)

func init() {
    encodeRef := map[byte]byte{
        '\b': 'b',
        '\f': 'f',
        '\n': 'n',
        '\r': 'r',
        '\t': 't',
        '\\': '\\',
    }

    for i := range EncodeMap {
        EncodeMap[i] = DontEscape
    }
    for i := range EncodeMap {
        if to, ok := encodeRef[byte(i)]; ok {
            EncodeMap[byte(i)] = to
        }
    }

    // underlyingHexMap contains the string "\x00\x01\x02..." which HexMap and
    // RawHexMap then index into.
    var underlyingHexMap bytes.Buffer
    underlyingHexMap.Grow(1024)

    for i := 0; i < 256; i++ {
        underlyingHexMap.WriteString("\\x")
        writeHexDigit(&underlyingHexMap, i/16)
        writeHexDigit(&underlyingHexMap, i%16)
    }

    underlyingHexBytes := underlyingHexMap.Bytes()

    for i := 0; i < 256; i++ {
        HexMap[i] = underlyingHexBytes[i*4 : i*4+4]
        RawHexMap[i] = underlyingHexBytes[i*4+2 : i*4+4]
    }
}

// EncodeEscapedChar is used internally to write out a character from a larger
// string that needs to be escaped to a buffer.
func EncodeEscapedChar(
    buf *bytes.Buffer,
    entireString string,
    currentRune rune,
    currentByte byte,
    currentIdx int,
    quoteChar byte,
) {
    ln := utf8.RuneLen(currentRune)
    if currentRune == utf8.RuneError {
        // Errors are due to invalid unicode points, so escape the bytes.
        // Make sure this is run at least once in case ln == -1.
        buf.Write(HexMap[entireString[currentIdx]])
        for ri := 1; ri < ln; ri++ {
            buf.Write(HexMap[entireString[currentIdx+ri]])
        }
    } else if ln == 1 {
        // For single-byte runes, do the same as encodeSQLBytes.
        if encodedChar := EncodeMap[currentByte]; encodedChar != DontEscape {
            buf.WriteByte('\\')
            buf.WriteByte(encodedChar)
        } else if currentByte == quoteChar {
            buf.WriteByte('\\')
            buf.WriteByte(quoteChar)
        } else {
            // Escape non-printable characters.
            buf.Write(HexMap[currentByte])
        }
    } else if ln == 2 {
        // For multi-byte runes, print them based on their width.
        fmt.Fprintf(buf, `\u%04X`, currentRune)
    } else {
        fmt.Fprintf(buf, `\U%08X`, currentRune)
    }
}

func writeHexDigit(buf *bytes.Buffer, v int) {
    if v < 10 {
        buf.WriteByte('0' + byte(v))
    } else {
        buf.WriteByte('a' + byte(v-10))
    }
}

// NeedEscape returns whether the given byte needs to be escaped.
func NeedEscape(ch byte) bool {
    return EncodeMap[ch] != DontEscape
}
