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

package pbutil

import (
	"io"

	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/proto"
)

type countingReader struct {
	r io.Reader
	n int
}

// implements protodelim.Reader
func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if n > 0 {
		r.n += n
	}
	return n, err
}

// implements protodelim.Reader
func (c *countingReader) ReadByte() (byte, error) {
	var buf [1]byte
	for {
		n, err := c.Read(buf[:])
		if n == 0 && err == nil {
			// io.Reader states: Callers should treat a return of 0 and nil as
			// indicating that nothing happened; in particular it does not
			// indicate EOF.
			continue
		}
		return buf[0], err
	}
}

// ReadDelimited decodes a message from the provided length-delimited stream,
// where the length is encoded as 32-bit varint prefix to the message body.
// It returns the total number of bytes read and any applicable error.  This is
// roughly equivalent to the companion Java API's
// MessageLite#parseDelimitedFrom.  As per the reader contract, this function
// calls r.Read repeatedly as required until exactly one message including its
// prefix is read and decoded (or an error has occurred).  The function never
// reads more bytes from the stream than required.  The function never returns
// an error if a message has been read and decoded correctly, even if the end
// of the stream has been reached in doing so.  In that case, any subsequent
// calls return (0, io.EOF).
func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	cr := &countingReader{r: r}
	opts := protodelim.UnmarshalOptions{
		MaxSize: -1,
	}
	err = opts.UnmarshalFrom(cr, m)
	return cr.n, err
}
