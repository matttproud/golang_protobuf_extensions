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
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"io"
)

// https://code.google.com/p/goprotobuf/source/browse/proto/encode.go?r=145dca00d164a4ab904098268d596823720702d0#69
const maxVarintBytes = 10

// ReadDelimited decodes a message from the provided length-delimited stream,
// where the length is encoded as 64-bit varint prefix to the message body.
// It returns the total number of bytes read and any applicable error.
func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	buffer := make([]byte, maxVarintBytes)

	headerLength, err := r.Read(buffer)
	if err != nil {
		return headerLength, err
	}
	if headerLength == 0 || int(buffer[0]) == -1 {
		return headerLength, io.EOF
	}

	messageLength, syncLength := proto.DecodeVarint(buffer)
	buffer = buffer[syncLength:]

	remainder := make([]byte, int(messageLength)-len(buffer))
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
