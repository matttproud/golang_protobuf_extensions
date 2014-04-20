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
	"fmt"
	"io"

	"code.google.com/p/goprotobuf/proto"
)

// ReadDelimited decodes a message from the provided length-delimited stream,
// where the length is encoded as 32-bit varint prefix to the message body.
// It returns the total number of bytes read and any applicable error.
//
// This API is slated for removal.
func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	deprReadDelimited()

	// Per AbstractParser#parsePartialDelimitedFrom with
	// CodedInputStream#readRawVarint32.
	buffer := make([]byte, binary.MaxVarintLen32)
	headerLength, err := r.Read(buffer)
	if err != nil {
		return headerLength, err
	}
	if headerLength == 0 || int(buffer[0]) == -1 {
		return headerLength, io.EOF
	}

	messageLength, syncLength := proto.DecodeVarint(buffer)
	buffer = buffer[syncLength:]

	remainderBufSize := int(messageLength) - len(buffer)
	if remainderBufSize <= 0 {
		return headerLength, proto.Unmarshal(buffer, m)
	}

	remainder := make([]byte, remainderBufSize)
	remainderLength, err := r.Read(remainder)
	if err != nil {
		return headerLength + remainderLength, err
	}
	if len(buffer)+len(remainder) != int(messageLength) {
		return headerLength + remainderLength, fmt.Errorf("truncated message")
	}

	buffer = append(buffer, remainder...)

	return headerLength + remainderLength, proto.Unmarshal(buffer, m)
}
