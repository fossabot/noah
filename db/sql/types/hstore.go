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

package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/readystock/noah/db/sql/pgio"
)

// Hstore represents an hstore column that can be null or have null values
// associated with its keys.
type Hstore struct {
	Map    map[string]Text
	Status Status
}

func (dst *Hstore) Set(src interface{}) error {
	if src == nil {
		*dst = Hstore{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case map[string]string:
		m := make(map[string]Text, len(value))
		for k, v := range value {
			m[k] = Text{String: v, Status: Present}
		}
		*dst = Hstore{Map: m, Status: Present}
	default:
		return errors.Errorf("cannot convert %v to Hstore", src)
	}

	return nil
}

func (dst *Hstore) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Map
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Hstore) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *map[string]string:
			*v = make(map[string]string, len(src.Map))
			for k, val := range src.Map {
				if val.Status != Present {
					return errors.Errorf("cannot decode %#v into %T", src, dst)
				}
				(*v)[k] = val.String
			}
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

func (dst *Hstore) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Hstore{Status: Null}
		return nil
	}

	keys, values, err := parseHstore(string(src))
	if err != nil {
		return err
	}

	m := make(map[string]Text, len(keys))
	for i := range keys {
		m[keys[i]] = values[i]
	}

	*dst = Hstore{Map: m, Status: Present}
	return nil
}

func (dst *Hstore) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Hstore{Status: Null}
		return nil
	}

	rp := 0

	if len(src[rp:]) < 4 {
		return errors.Errorf("hstore incomplete %v", src)
	}
	pairCount := int(int32(binary.BigEndian.Uint32(src[rp:])))
	rp += 4

	m := make(map[string]Text, pairCount)

	for i := 0; i < pairCount; i++ {
		if len(src[rp:]) < 4 {
			return errors.Errorf("hstore incomplete %v", src)
		}
		keyLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4

		if len(src[rp:]) < keyLen {
			return errors.Errorf("hstore incomplete %v", src)
		}
		key := string(src[rp : rp+keyLen])
		rp += keyLen

		if len(src[rp:]) < 4 {
			return errors.Errorf("hstore incomplete %v", src)
		}
		valueLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4

		var valueBuf []byte
		if valueLen >= 0 {
			valueBuf = src[rp : rp+valueLen]
		}
		rp += valueLen

		var value Text
		err := value.DecodeBinary(ci, valueBuf)
		if err != nil {
			return err
		}
		m[key] = value
	}

	*dst = Hstore{Map: m, Status: Present}

	return nil
}

func (src *Hstore) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	firstPair := true

	for k, v := range src.Map {
		if firstPair {
			firstPair = false
		} else {
			buf = append(buf, ',')
		}

		buf = append(buf, quoteHstoreElementIfNeeded(k)...)
		buf = append(buf, "=>"...)

		elemBuf, err := v.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}

		if elemBuf == nil {
			buf = append(buf, "NULL"...)
		} else {
			buf = append(buf, quoteHstoreElementIfNeeded(string(elemBuf))...)
		}
	}

	return buf, nil
}

func (src *Hstore) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendInt32(buf, int32(len(src.Map)))

	var err error
	for k, v := range src.Map {
		buf = pgio.AppendInt32(buf, int32(len(k)))
		buf = append(buf, k...)

		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		elemBuf, err := v.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		}
		if elemBuf != nil {
			buf = elemBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
	}

	return buf, err
}

var quoteHstoreReplacer = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func quoteHstoreElement(src string) string {
	return `"` + quoteArrayReplacer.Replace(src) + `"`
}

func quoteHstoreElementIfNeeded(src string) string {
	if src == "" || (len(src) == 4 && strings.ToLower(src) == "null") || strings.ContainsAny(src, ` {},"\=>`) {
		return quoteArrayElement(src)
	}
	return src
}

const (
	hsPre = iota
	hsKey
	hsSep
	hsVal
	hsNul
	hsNext
)

type hstoreParser struct {
	str string
	pos int
}

func newHSP(in string) *hstoreParser {
	return &hstoreParser{
		pos: 0,
		str: in,
	}
}

func (p *hstoreParser) Consume() (r rune, end bool) {
	if p.pos >= len(p.str) {
		end = true
		return
	}
	r, w := utf8.DecodeRuneInString(p.str[p.pos:])
	p.pos += w
	return
}

func (p *hstoreParser) Peek() (r rune, end bool) {
	if p.pos >= len(p.str) {
		end = true
		return
	}
	r, _ = utf8.DecodeRuneInString(p.str[p.pos:])
	return
}

// parseHstore parses the string representation of an hstore column (the same
// you would get from an ordinary SELECT) into two slices of keys and values. it
// is used internally in the default parsing of hstores.
func parseHstore(s string) (k []string, v []Text, err error) {
	if s == "" {
		return
	}

	buf := bytes.Buffer{}
	keys := []string{}
	values := []Text{}
	p := newHSP(s)

	r, end := p.Consume()
	state := hsPre

	for !end {
		switch state {
		case hsPre:
			if r == '"' {
				state = hsKey
			} else {
				err = errors.New("String does not begin with \"")
			}
		case hsKey:
			switch r {
			case '"': //End of the key
				if buf.Len() == 0 {
					err = errors.New("Empty Key is invalid")
				} else {
					keys = append(keys, buf.String())
					buf = bytes.Buffer{}
					state = hsSep
				}
			case '\\': //Potential escaped character
				n, end := p.Consume()
				switch {
				case end:
					err = errors.New("Found EOS in key, expecting character or \"")
				case n == '"', n == '\\':
					buf.WriteRune(n)
				default:
					buf.WriteRune(r)
					buf.WriteRune(n)
				}
			default: //Any other character
				buf.WriteRune(r)
			}
		case hsSep:
			if r == '=' {
				r, end = p.Consume()
				switch {
				case end:
					err = errors.New("Found EOS after '=', expecting '>'")
				case r == '>':
					r, end = p.Consume()
					switch {
					case end:
						err = errors.New("Found EOS after '=>', expecting '\"' or 'NULL'")
					case r == '"':
						state = hsVal
					case r == 'N':
						state = hsNul
					default:
						err = errors.Errorf("Invalid character '%c' after '=>', expecting '\"' or 'NULL'", r)
					}
				default:
					err = errors.Errorf("Invalid character after '=', expecting '>'")
				}
			} else {
				err = errors.Errorf("Invalid character '%c' after value, expecting '='", r)
			}
		case hsVal:
			switch r {
			case '"': //End of the value
				values = append(values, Text{String: buf.String(), Status: Present})
				buf = bytes.Buffer{}
				state = hsNext
			case '\\': //Potential escaped character
				n, end := p.Consume()
				switch {
				case end:
					err = errors.New("Found EOS in key, expecting character or \"")
				case n == '"', n == '\\':
					buf.WriteRune(n)
				default:
					buf.WriteRune(r)
					buf.WriteRune(n)
				}
			default: //Any other character
				buf.WriteRune(r)
			}
		case hsNul:
			nulBuf := make([]rune, 3)
			nulBuf[0] = r
			for i := 1; i < 3; i++ {
				r, end = p.Consume()
				if end {
					err = errors.New("Found EOS in NULL value")
					return
				}
				nulBuf[i] = r
			}
			if nulBuf[0] == 'U' && nulBuf[1] == 'L' && nulBuf[2] == 'L' {
				values = append(values, Text{Status: Null})
				state = hsNext
			} else {
				err = errors.Errorf("Invalid NULL value: 'N%s'", string(nulBuf))
			}
		case hsNext:
			if r == ',' {
				r, end = p.Consume()
				switch {
				case end:
					err = errors.New("Found EOS after ',', expcting space")
				case unicode.IsSpace(r):
					r, end = p.Consume()
					state = hsKey
				default:
					err = errors.Errorf("Invalid character '%c' after ', ', expecting \"", r)
				}
			} else {
				err = errors.Errorf("Invalid character '%c' after value, expecting ','", r)
			}
		}

		if err != nil {
			return
		}
		r, end = p.Consume()
	}
	if state != hsNext {
		err = errors.New("Improperly formatted hstore")
		return
	}
	k = keys
	v = values
	return
}

// Scan implements the database/sql Scanner interface.
func (dst *Hstore) Scan(src interface{}) error {
	if src == nil {
		*dst = Hstore{Status: Null}
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
func (src *Hstore) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
