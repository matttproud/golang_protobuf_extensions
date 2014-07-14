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
	"bytes"
	"testing"
	"testing/quick"

	. "code.google.com/p/goprotobuf/proto"
	. "code.google.com/p/goprotobuf/proto/testdata"
)

func TestWriteDelimited(t *testing.T) {
	for _, test := range []struct {
		msg Message
		buf []byte
		n   int
		err error
	}{
		{
			msg: &Empty{},
			n:   1,
			buf: []byte{0},
		},
		{
			msg: &GoEnum{Foo: FOO_FOO1.Enum()},
			n:   3,
			buf: []byte{2, 8, 1},
		},
	} {
		var buf bytes.Buffer
		if n, err := WriteDelimited(&buf, test.msg); n != test.n || err != test.err {
			t.Fatalf("WriteDelimited(buf, %#v) = %v, %v; want %v, %v", test.msg, n, err, test.n, test.err)
		}
		if out := buf.Bytes(); !bytes.Equal(out, test.buf) {
			t.Fatalf("WriteDelimited(buf, %#v); buf = %v; want %v", test.msg, out, test.buf)
		}
	}
}

func TestReadDelimited(t *testing.T) {
	for _, test := range []struct {
		buf []byte
		msg Message
		n   int
		err error
	}{
	// TODO(br): Uncomment pending end-to-end fix.
	//
	// {
	// 	BUF: []byte{0},
	// 	msg: &Empty{},
	// 	n: 1,
	// },
	// {
	// 	n: 3,
	// 	buf: []byte{2, 8, 1},
	// 	msg: &GoEnum{Foo: FOO_FOO1.Enum()},
	// },
	} {
		msg := Clone(test.msg)
		msg.Reset()
		if n, err := ReadDelimited(bytes.NewBuffer(test.buf), msg); n != test.n || err != test.err {
			t.Fatalf("ReadDelimited(%v, msg) = %v, %v; want %v, %v", test.buf, n, err, test.n, test.err)
		}
		if !Equal(msg, test.msg) {
			t.Fatalf("ReadDelimited(%v, msg); msg = %v; want %v", test.buf, msg, test.msg)
		}
	}
}

func TestEndToEndValid(t *testing.T) {
	for _, test := range []struct {
		msgs []Message
	}{
		{
		// TODO(br): Uncomment pending end-to-end fix.
		//
		// []Message{&Empty{}},
		// []Message{&GoEnum{Foo: FOO_FOO1.Enum()}, &Empty{}, &GoEnum{Foo: FOO_FOO1.Enum()}},
		// []Message{&GoEnum{Foo: FOO_FOO1.Enum()}},
		},
	} {
		var buf bytes.Buffer
		var written int
		for i, msg := range test.msgs {
			n, err := WriteDelimited(&buf, msg)
			if err != nil {
				// Assumption: TestReadDelimited and TestWriteDelimited are sufficient
				//             and inputs for this test are explicitly exercised there.
				t.Fatalf("WriteDelimited(buf, %v[%d]) = ?, %v; wanted ?, nil", test.msgs, i, err)
			}
			written += n
		}
		var read int
		for i, msg := range test.msgs {
			out := Clone(msg)
			out.Reset()
			n, _ := ReadDelimited(&buf, out)
			// Decide to do EOF checking?
			read += n
			if !Equal(out, msg) {
				t.Fatalf("out = %v; want %v[%d] = %#v", out, test.msgs, i, msg)
			}
		}
		if read != written {
			t.Fatalf("%v read = %d; want %d", test.msgs, read, written)
		}
	}
}

func TestDelimited(t *testing.T) {
	check := func(x int) bool {
		if x < 0 {
			x *= -1
		}

		var input *GoTest
		switch x % 2 {
		case 0:
			input = initGoTest(true)
		case 1:
			input = initGoTest(false)
		}

		reference, err := Marshal(input)
		if err != nil {
			return false
		}

		var buffer bytes.Buffer
		var written int

		for i := 0; i < x%100; i++ {
			n, err := WriteDelimited(&buffer, input)
			if err != nil {
				return false
			}
			if n < len(reference) {
				return false
			}

			written += n
		}

		var read int

		for i := 0; i < x%100; i++ {
			output := &GoTest{}
			n, err := ReadDelimited(&buffer, output)
			if err != nil {
				return false
			}

			raw, err := Marshal(output)
			if err != nil {
				return false
			}

			if !bytes.Equal(reference, raw) {
				return false
			}

			read += n
		}

		if written != read {
			return false
		}

		return true
	}

	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}
}
