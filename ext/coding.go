// Copyright 2014 Matt T. Proud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ext

import (
	"encoding/binary"
	"errors"
	"io"

	"code.google.com/p/goprotobuf/proto"
)

// DelimitedBuffer provides varint record length-delimited Protocol Buffer
// message encoding, with the feature that all internal buffers are reused to
// reduce memory usage overhead.  The type is not goroutine safe.
type DelimitedBuffer struct {
	msgBuf  proto.Buffer
	buf     []byte
	headBuf []byte
}

// Marshal encodes a Protocol Buffer message to this DelimitedBuffer's internal
// byte buffer, with the binary message prefixed by the varint encoded size of
// the written message.  It returns the number of bytes written to the buffer
// along with any error it may have encountered.
func (b *DelimitedBuffer) Marshal(m proto.Message) (n int, err error) {
	if err = b.msgBuf.Marshal(m); err != nil {
		return 0, err
	}
	if len(b.headBuf) < binary.MaxVarintLen64 {
		b.headBuf = make([]byte, binary.MaxVarintLen64)
	}
	msg := b.msgBuf.Bytes()
	msgLen := len(msg)
	headLen := binary.PutUvarint(b.headBuf, uint64(msgLen))
	b.buf = append(b.buf, b.headBuf[0:headLen]...)
	b.buf = append(b.buf, msg...)
	return headLen + msgLen, nil
}

var ErrTruncMsg = errors.New("truncated message")

// Unmarshal decodes a Protocol Buffer message from this DelimitedBuffer's
// internal and writes it to the provided message.  A read advances the internal
// buffer position accordingly.  It returns the number of bytes written to the
// buffer along with any error it may have encountered.
func (b *DelimitedBuffer) Unmarshal(m proto.Message) (n int, err error) {
	// Per AbstractParser#parsePartialDelimitedFrom with
	// CodedInputStream#readRawVarint32.
	bufLen := len(b.buf)
	if bufLen == 0 || int(b.buf[0]) == -1 {
		return 0, io.EOF
	}
	scanLen := binary.MaxVarintLen64
	if bufLen < scanLen {
		scanLen = bufLen
	}
	msgLen, syncLen := proto.DecodeVarint(b.buf)
	b.buf = b.buf[syncLen:]
	if len(b.buf)-int(msgLen) < 0 {
		return syncLen, ErrTruncMsg
	}
	wind := b.buf[0:msgLen]
	b.buf = b.buf[msgLen:]
	err = proto.Unmarshal(wind, m)
	if err != nil {
		return syncLen, err
	}
	return syncLen + int(msgLen), nil
}

// Bytes yields the internal byte array buffer for this instance.
func (b *DelimitedBuffer) Bytes() []byte {
	return b.buf
}

// SetBuf replaces the internal buffer with the provided slice, thereby enabling
// the DelimitedWriter to either begin decoding from its initial position or
// begin writing thereto.
func (b *DelimitedBuffer) SetBuf(buf []byte) {
	b.buf = buf
}

// Clear clears the internal buffer, thereby deleting any values written to it.
func (b *DelimitedBuffer) Clear() {
	b.msgBuf.Reset()
	b.buf = b.buf[0:0]
}

// Reset returns the buffer to a state similar to that of a new instance.  The
// use case for this, which differs from Clear, is that you may have encoded an
// unusually large Protocol Buffer message and want to signal to the memory
// manager that the large buffer should be reaped eventually.
func (b *DelimitedBuffer) Reset() {
	b.msgBuf = proto.Buffer{}
	b.buf = nil
}

// NewDelimitedBuffer allocates a new DelimitedBuffer and initializes its
// internal buffer with the contents of the provided slice.  Any writes to
// this DelimitedBuffer will be written to the provided buffer as well.
func NewDelimitedBuffer(buf []byte) *DelimitedBuffer {
	return &DelimitedBuffer{
		buf: buf,
	}
}
