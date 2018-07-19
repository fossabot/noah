// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in licenses/BSD-vitess.txt.

// Portions of this file are additionally subject to the following
// license and copyright.
//
// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

// This code was derived from https://github.com/youtube/vitess.

package lex

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/Ready-Stock/Noah/Database/sql/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/Database/sql/sessiondata"
	"github.com/Ready-Stock/Noah/Database/util/stringencoding"
)

var mustQuoteMap = map[byte]bool{
	' ': true,
	',': true,
	'{': true,
	'}': true,
}

// EncodeFlags influence the formatting of strings and identifiers.
type EncodeFlags int

// HasFlags tests whether the given flags are set.
func (f EncodeFlags) HasFlags(subset EncodeFlags) bool {
	return f&subset == subset
}

const (
	// EncNoFlags indicates nothing special should happen while encoding.
	EncNoFlags EncodeFlags = 0

	// EncBareStrings indicates that strings will be rendered without
	// wrapping quotes if they contain no special characters.
	EncBareStrings EncodeFlags = 1 << iota

	// EncBareIdentifiers indicates that identifiers will be rendered
	// without wrapping quotes.
	EncBareIdentifiers

	// EncFirstFreeFlagBit needs to remain unused; it is used as base
	// bit offset for tree.FmtFlags.
	EncFirstFreeFlagBit
)

// EncodeSQLString writes a string literal to buf. All unicode and
// non-printable characters are escaped.
func EncodeSQLString(buf *bytes.Buffer, in string) {
	EncodeSQLStringWithFlags(buf, in, EncNoFlags)
}

// EscapeSQLString returns an escaped SQL representation of the given
// string. This is suitable for safely producing a SQL string valid
// for input to the parser.
func EscapeSQLString(in string) string {
	var buf bytes.Buffer
	EncodeSQLString(&buf, in)
	return buf.String()
}

// EncodeSQLStringWithFlags writes a string literal to buf. All
// unicode and non-printable characters are escaped. flags controls
// the output format: if encodeBareString is set, the output string
// will not be wrapped in quotes if the strings contains no special
// characters.
func EncodeSQLStringWithFlags(buf *bytes.Buffer, in string, flags EncodeFlags) {
	// See http://www.postgresql.org/docs/9.4/static/sql-syntax-lexical.html
	start := 0
	escapedString := false
	bareStrings := flags.HasFlags(EncBareStrings)
	// Loop through each unicode code point.
	for i, r := range in {
		ch := byte(r)
		if r >= 0x20 && r < 0x7F {
			if mustQuoteMap[ch] {
				// We have to quote this string - ignore bareStrings setting
				bareStrings = false
			}
			if !stringencoding.NeedEscape(ch) && ch != '\'' {
				continue
			}
		}

		if !escapedString {
			buf.WriteString("e'") // begin e'xxx' string
			escapedString = true
		}
		buf.WriteString(in[start:i])
		ln := utf8.RuneLen(r)
		if ln < 0 {
			start = i + 1
		} else {
			start = i + ln
		}
		stringencoding.EncodeEscapedChar(buf, in, r, ch, i, '\'')
	}

	quote := !escapedString && !bareStrings
	if quote {
		buf.WriteByte('\'') // begin 'xxx' string if nothing was escaped
	}
	buf.WriteString(in[start:])
	if escapedString || quote {
		buf.WriteByte('\'')
	}
}

// EncodeSQLStringInsideArray writes a string literal to buf using the "string
// within array" formatting.
func EncodeSQLStringInsideArray(buf *bytes.Buffer, in string) {
	buf.WriteByte('"')
	// Loop through each unicode code point.
	for i, r := range in {
		ch := byte(r)
		if unicode.IsPrint(r) && !stringencoding.NeedEscape(ch) && ch != '"' {
			// Character is printable doesn't need escaping - just print it out.
			buf.WriteRune(r)
		} else {
			stringencoding.EncodeEscapedChar(buf, in, r, ch, i, '"')
		}
	}

	buf.WriteByte('"')
}

// EncodeUnrestrictedSQLIdent writes the identifier in s to buf.
// The identifier is only quoted if the flags don't tell otherwise and
// the identifier contains special characters.
func EncodeUnrestrictedSQLIdent(buf *bytes.Buffer, s string, flags EncodeFlags) {
	if flags.HasFlags(EncBareIdentifiers) || isBareIdentifier(s) {
		buf.WriteString(s)
		return
	}
	encodeEscapedSQLIdent(buf, s)
}

// EncodeRestrictedSQLIdent writes the identifier in s to buf. The
// identifier is quoted if either the flags ask for it, the identifier
// contains special characters, or the identifier is a reserved SQL
// keyword.
func EncodeRestrictedSQLIdent(buf *bytes.Buffer, s string, flags EncodeFlags) {
	if flags.HasFlags(EncBareIdentifiers) || (!isReservedKeyword(s) && isBareIdentifier(s)) {
		buf.WriteString(s)
		return
	}
	encodeEscapedSQLIdent(buf, s)
}

// encodeEscapedSQLIdent writes the identifier in s to buf. The
// identifier is always quoted. Double quotes inside the identifier
// are escaped.
func encodeEscapedSQLIdent(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	start := 0
	for i, n := 0, len(s); i < n; i++ {
		ch := s[i]
		// The only character that requires escaping is a double quote.
		if ch == '"' {
			if start != i {
				buf.WriteString(s[start:i])
			}
			start = i + 1
			buf.WriteByte(ch)
			buf.WriteByte(ch) // add extra copy of ch
		}
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
}

// EncodeSQLBytes encodes the SQL byte array in 'in' to buf, to a
// format suitable for re-scanning. We don't use a straightforward hex
// encoding here with x'...'  because the result would be less
// compact. We are trading a little more time during the encoding to
// have a little less bytes on the wire.
func EncodeSQLBytes(buf *bytes.Buffer, in string) {
	start := 0
	buf.WriteString("b'")
	// Loop over the bytes of the string (i.e., don't use range over unicode
	// code points).
	for i, n := 0, len(in); i < n; i++ {
		ch := in[i]
		if encodedChar := stringencoding.EncodeMap[ch]; encodedChar != stringencoding.DontEscape {
			buf.WriteString(in[start:i])
			buf.WriteByte('\\')
			buf.WriteByte(encodedChar)
			start = i + 1
		} else if ch == '\'' {
			// We can't just fold this into stringencoding.EncodeMap because
			// stringencoding.EncodeMap is also used for strings which
			// aren't quoted with single-quotes
			buf.WriteString(in[start:i])
			buf.WriteByte('\\')
			buf.WriteByte(ch)
			start = i + 1
		} else if ch < 0x20 || ch >= 0x7F {
			buf.WriteString(in[start:i])
			// Escape non-printable characters.
			buf.Write(stringencoding.HexMap[ch])
			start = i + 1
		}
	}
	buf.WriteString(in[start:])
	buf.WriteByte('\'')
}

// EncodeByteArrayToRawBytes converts a SQL-level byte array into raw
// bytes according to the encoding specification in "be".
// If the skipHexPrefix argument is set, the hexadecimal encoding does not
// prefix the output with "\x". This is suitable e.g. for the encode()
// built-in.
func EncodeByteArrayToRawBytes(
	data string, be sessiondata.BytesEncodeFormat, skipHexPrefix bool,
) string {
	switch be {
	case sessiondata.BytesEncodeHex:
		head := 2
		if skipHexPrefix {
			head = 0
		}
		res := make([]byte, head+hex.EncodedLen(len(data)))
		if !skipHexPrefix {
			res[0] = '\\'
			res[1] = 'x'
		}
		hex.Encode(res[head:], []byte(data))
		return string(res)

	case sessiondata.BytesEncodeEscape:
		// PostgreSQL does not allow all the escapes formats recognized by
		// CockroachDB's scanner. It only recognizes octal and \\ for the
		// backslash itself.
		// See https://www.postgresql.org/docs/current/static/datatype-binary.html#AEN5667
		res := make([]byte, 0, len(data))
		for _, c := range []byte(data) {
			if c == '\\' {
				res = append(res, '\\', '\\')
			} else if c < 32 || c >= 127 {
				// Escape the character in octal.
				//
				// Note: CockroachDB only supports UTF-8 for which all values
				// below 128 are ASCII. There is no locale-dependent escaping
				// in that case.
				res = append(res, '\\', '0'+(c>>6), '0'+((c>>3)&7), '0'+(c&7))
			} else {
				res = append(res, c)
			}
		}
		return string(res)

	case sessiondata.BytesEncodeBase64:
		return base64.StdEncoding.EncodeToString([]byte(data))

	default:
		panic(fmt.Sprintf("unhandled format: %s", be))
	}
}

// DecodeRawBytesToByteArray converts raw bytes to a SQL-level byte array
// according to the encoding specification in "be".
// When using the Hex format, the caller is responsible for skipping the
// "\x" prefix, if any. See DecodeRawBytesToByteArrayAuto() below for
// an alternative.
func DecodeRawBytesToByteArray(data string, be sessiondata.BytesEncodeFormat) ([]byte, error) {
	switch be {
	case sessiondata.BytesEncodeHex:
		return hex.DecodeString(data)

	case sessiondata.BytesEncodeEscape:
		// PostgreSQL does not allow all the escapes formats recognized by
		// CockroachDB's scanner. It only recognizes octal and \\ for the
		// backslash itself.
		// See https://www.postgresql.org/docs/current/static/datatype-binary.html#AEN5667
		res := make([]byte, 0, len(data))
		for i := 0; i < len(data); i++ {
			ch := data[i]
			if ch != '\\' {
				res = append(res, ch)
				continue
			}
			if i >= len(data)-1 {
				return nil, pgerror.NewError(pgerror.CodeInvalidEscapeSequenceError,
					"bytea encoded value ends with escape character")
			}
			if data[i+1] == '\\' {
				res = append(res, '\\')
				i++
				continue
			}
			if i+3 >= len(data) {
				return nil, pgerror.NewError(pgerror.CodeInvalidEscapeSequenceError,
					"bytea encoded value ends with incomplete escape sequence")
			}
			b := byte(0)
			for j := 1; j <= 3; j++ {
				octDigit := data[i+j]
				if octDigit < '0' || octDigit > '7' {
					return nil, pgerror.NewError(pgerror.CodeInvalidEscapeSequenceError,
						"invalid bytea escape sequence")
				}
				b = (b << 3) | (octDigit - '0')
			}
			res = append(res, b)
			i += 3
		}
		return res, nil

	case sessiondata.BytesEncodeBase64:
		return base64.StdEncoding.DecodeString(data)

	default:
		return nil, pgerror.NewErrorf(pgerror.CodeInternalError,
			"programming error: unhandled format: %s", be)
	}
}

// DecodeRawBytesToByteArrayAuto detects which format to use with
// DecodeRawBytesToByteArray(). It only supports hex ("\x" prefix)
// and escape.
func DecodeRawBytesToByteArrayAuto(data []byte) ([]byte, error) {
	if len(data) >= 2 && data[0] == '\\' && (data[1] == 'x' || data[1] == 'X') {
		return DecodeRawBytesToByteArray(string(data[2:]), sessiondata.BytesEncodeHex)
	}
	return DecodeRawBytesToByteArray(string(data), sessiondata.BytesEncodeEscape)
}
