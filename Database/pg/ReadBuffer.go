package pg

import (
	"io"
	"encoding/binary"
	"github.com/Ready-Stock/Noah/Database/pg/pgwire/pgerror"
	"bytes"
	"unsafe"
)

const secondsInDay = 24 * 60 * 60
const maxMessageSize = 1 << 24

type ReadBuffer struct {
	Msg []byte
	tmp [4]byte
}

func (b *ReadBuffer) ReadUntypedMsg(rd io.Reader) (int, error) {
	nread, err := io.ReadFull(rd, b.tmp[:])
	if err != nil {
		return nread, err
	}
	size := int(binary.BigEndian.Uint32(b.tmp[:]))
	// size includes itself.
	size -= 4
	if size > maxMessageSize || size < 0 {
		return nread, NewProtocolViolationErrorf("message size %d out of bounds (0..%d)",
			size, maxMessageSize)
	}

	b.reset(size)
	n, err := io.ReadFull(rd, b.Msg)
	return nread + n, err
}

func (b *ReadBuffer) GetString() (string, error) {
	pos := bytes.IndexByte(b.Msg, 0)
	if pos == -1 {
		return "", NewProtocolViolationErrorf("NUL terminator not found")
	}
	// Note: this is a conversion from a byte slice to a string which avoids
	// allocation and copying. It is safe because we never reuse the bytes in our
	// read buffer. It is effectively the same as: "s := string(b.Msg[:pos])"
	s := b.Msg[:pos]
	b.Msg = b.Msg[pos+1:]
	return *((*string)(unsafe.Pointer(&s))), nil
}



// GetUint32 returns the buffer's contents as a uint32.
func (b *ReadBuffer) GetUint32() (uint32, error) {
	if len(b.Msg) < 4 {
		return 0, NewProtocolViolationErrorf("insufficient data: %d", len(b.Msg))
	}
	v := binary.BigEndian.Uint32(b.Msg[:4])
	b.Msg = b.Msg[4:]
	return v, nil
}



func (b *ReadBuffer) reset(size int) {
	if b.Msg != nil {
		b.Msg = b.Msg[len(b.Msg):]
	}

	if cap(b.Msg) >= size {
		b.Msg = b.Msg[:size]
		return
	}

	allocSize := size
	if allocSize < 4096 {
		allocSize = 4096
	}
	b.Msg = make([]byte, size, allocSize)
}


// NewProtocolViolationErrorf creates a pgwire ProtocolViolationError.
func NewProtocolViolationErrorf(format string, args ...interface{}) error {
	return pgerror.NewErrorf(pgerror.CodeProtocolViolationError, format, args...)
}