// Copyright 2013 Matt T. Proud
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

var errInvalidVarint = errors.New("invalid varint32 encountered")

// ReadDelimited decodes a message from the provided length-delimited stream,
// where the length is encoded as 32-bit varint prefix to the message body.
// It returns the total number of bytes read and any applicable error.
func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	// Per AbstractParser#parsePartialDelimitedFrom with
	// CodedInputStream#readRawVarint32.
	headerBuf := make([]byte, binary.MaxVarintLen32)
	var bytesRead, varIntBytes int
	var messageLength uint64
	for varIntBytes == 0 {
		if bytesRead >= len(headerBuf) {
			return bytesRead, errInvalidVarint
		}
		newBytesRead, err := r.Read(headerBuf[bytesRead:])
		if newBytesRead == 0 {
			if err != nil {
				return bytesRead, err
			}
			// A Reader should not return (0, nil), but if it does,
			// it should be treated as no-op (according to the
			// Reader contract). So let's go on...
			continue
		}
		bytesRead += newBytesRead
		messageLength, varIntBytes = proto.DecodeVarint(headerBuf)
	}

	headerBuf = headerBuf[varIntBytes:] // Need to process what's not used yet in headerBuf.

	if messageLength-uint64(len(headerBuf)) <= 0 {
		return bytesRead, proto.Unmarshal(headerBuf, m)
	}

	messageBuf := make([]byte, messageLength)
	copy(messageBuf, headerBuf)
	newBytesRead, err := io.ReadFull(r, messageBuf[len(headerBuf):])
	bytesRead += newBytesRead
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, proto.Unmarshal(messageBuf, m)
}
