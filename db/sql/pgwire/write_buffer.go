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

package pgwire

import (
	"bytes"
	"encoding/binary"
    "github.com/Ready-Stock/noah/db/sql/pgwire/pgwirebase"
    "github.com/Ready-Stock/noah/db/sql/sem/tree"
    "github.com/Ready-Stock/noah/db/sql/types"
    "github.com/Ready-Stock/noah/db/util"
	"io"
)

// writeBuffer is a wrapper around bytes.Buffer that provides a convenient interface
// for writing PGWire results. The buffer preserves any errors it encounters when writing,
// and will turn all subsequent write attempts into no-ops until finishMsg is called.
type writeBuffer struct {
	noCopy util.NoCopy

	wrapped bytes.Buffer
	err     error

	// These two buffers are used as temporary storage. Use putbuf when the
	// length of the required temp space is known. Use variablePutbuf when the length
	// of the required temp space is unknown, or when a bytes.Buffer is needed.
	//
	// We keep both of these because there are operations that are only possible to
	// perform (efficiently) with one or the other, such as strconv.AppendInt with
	// putbuf or Datum.Format with variablePutbuf.
	putbuf          [64]byte
	variablePutbuf  bytes.Buffer
	simpleFormatter tree.FmtCtx
	arrayFormatter  tree.FmtCtx
}

func newWriteBuffer() *writeBuffer {
	b := &writeBuffer{
	}
	b.simpleFormatter = tree.MakeFmtCtx(&b.variablePutbuf, tree.FmtSimple)
	b.arrayFormatter = tree.MakeFmtCtx(&b.variablePutbuf, tree.FmtArrays)
	return b
}

// Write implements the io.Write interface.
func (b *writeBuffer) Write(p []byte) (int, error) {
	b.write(p)
	return len(p), b.err
}

func (b *writeBuffer) writeByte(c byte) {
	if b.err == nil {
		b.err = b.wrapped.WriteByte(c)
	}
}

func (b *writeBuffer) write(p []byte) {
	if b.err == nil {
		_, b.err = b.wrapped.Write(p)
	}
}

func (b *writeBuffer) writeString(s string) {
	if b.err == nil {
		_, b.err = b.wrapped.WriteString(s)
	}
}

func (b *writeBuffer) nullTerminate() {
	if b.err == nil {
		b.err = b.wrapped.WriteByte(0)
	}
}

// writeLengthPrefixedVariablePutbuf writes the current contents of
// variablePutbuf with a length prefix. The function will reset
// variablePutbuf.
func (b *writeBuffer) writeLengthPrefixedVariablePutbuf() {
	if b.err == nil {
		b.putInt32(int32(b.variablePutbuf.Len()))

		// bytes.Buffer.WriteTo resets the Buffer.
		_, b.err = b.variablePutbuf.WriteTo(&b.wrapped)
	}
}

// writeLengthPrefixedBuffer writes the contents of a bytes.Buffer with a
// length prefix.
func (b *writeBuffer) writeLengthPrefixedBuffer(buf *bytes.Buffer) {
	if b.err == nil {
		b.putInt32(int32(buf.Len()))

		// bytes.Buffer.WriteTo resets the Buffer.
		_, b.err = buf.WriteTo(&b.wrapped)
	}
}

func (b *writeBuffer) writeLengthPrefixedBytes(bytes []byte) {
	if b.err == nil {
		b.putInt32(int32(len(bytes)))
		b.write(bytes)
	}
}

// writeLengthPrefixedString writes a length-prefixed string. The
// length is encoded as an int32.
func (b *writeBuffer) writeLengthPrefixedString(s string) {
	b.putInt32(int32(len(s)))
	b.writeString(s)
}

// writeLengthPrefixedDatum writes a length-prefixed Datum in its
// string representation. The length is encoded as an int32.
func (b *writeBuffer) writeLengthPrefixedDatum(d tree.Datum) {
	fmtCtx := tree.MakeFmtCtx(&b.variablePutbuf, tree.FmtSimple)
	fmtCtx.FormatNode(d)
	b.writeLengthPrefixedVariablePutbuf()
}

// writeLengthPrefixedPGValue writes a length-prefixed Value in its
// string representation. The length is encoded as an int32.
func (b *writeBuffer) writeLengthPrefixedPGValue(d types.Value) {
	// fmtCtx := tree.MakeFmtCtx(&b.variablePutbuf, tree.FmtSimple)
	// fmtCtx.FormatNode(d)
	b.writeLengthPrefixedVariablePutbuf()
}

// writeTerminatedString writes a null-terminated string.
func (b *writeBuffer) writeTerminatedString(s string) {
	b.writeString(s)
	b.nullTerminate()
}

func (b *writeBuffer) putInt16(v int16) {
	if b.err == nil {
		binary.BigEndian.PutUint16(b.putbuf[:], uint16(v))
		_, b.err = b.wrapped.Write(b.putbuf[:2])
	}
}

func (b *writeBuffer) putInt32(v int32) {
	if b.err == nil {
		binary.BigEndian.PutUint32(b.putbuf[:], uint32(v))
		_, b.err = b.wrapped.Write(b.putbuf[:4])
	}
}

func (b *writeBuffer) putInt64(v int64) {
	if b.err == nil {
		binary.BigEndian.PutUint64(b.putbuf[:], uint64(v))
		_, b.err = b.wrapped.Write(b.putbuf[:8])
	}
}

func (b *writeBuffer) putErrFieldMsg(field pgwirebase.ServerErrFieldType) {
	if b.err == nil {
		b.err = b.wrapped.WriteByte(byte(field))
	}
}

func (b *writeBuffer) reset() {
	b.wrapped.Reset()
	b.err = nil
}

// initMsg begins writing a message into the writeBuffer with the provided type.
func (b *writeBuffer) initMsg(typ pgwirebase.ServerMessageType) {
	b.reset()
	b.putbuf[0] = byte(typ)
	_, b.err = b.wrapped.Write(b.putbuf[:5]) // message type + message length
}

// finishMsg attempts to write the data it has accumulated to the provided io.Writer.
// If the writeBuffer previously encountered an error since the last call to initMsg,
// or if it encounters an error while writing to w, it will return an error.
func (b *writeBuffer) finishMsg(w io.Writer) error {
	defer b.reset()
	if b.err != nil {
		return b.err
	}
	bytes := b.wrapped.Bytes()
	binary.BigEndian.PutUint32(bytes[1:5], uint32(b.wrapped.Len()-1))
	_, err := w.Write(bytes)
	return err
}

// setError sets the writeBuffer's error, if it does not already have one.
func (b *writeBuffer) setError(err error) {
	if b.err == nil {
		b.err = err
	}
}
