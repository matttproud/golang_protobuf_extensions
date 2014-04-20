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
	"io"
	"testing"
	"testing/quick"

	. "code.google.com/p/goprotobuf/proto"
	. "code.google.com/p/goprotobuf/proto/testdata"
)

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
			t.Fatal(err)
			return false
		}

		var buf DelimitedBuffer
		var written int

		ents := x % 100

		for i := 0; i < ents; i++ {
			n, err := buf.Marshal(input)
			if err != nil {
				t.Fatal(err)
				return false
			}
			if n < len(reference) {
				t.Fatal(err)
				return false
			}

			written += n
		}

		var read int

		for i := 0; i < ents; i++ {
			output := new(GoTest)
			n, err := buf.Unmarshal(output)
			if err != nil {
				t.Fatal(err)
				return false
			}

			raw, err := Marshal(output)
			if err != nil {
				t.Fatal(err)
				return false
			}

			if !bytes.Equal(reference, raw) {
				t.Fatal("not equal")
				return false
			}

			read += n
		}

		if written != read {
			t.Fatalf("read != written %d %d %d ", written, read, ents)
			return false
		}

		return true
	}

	if err := quick.Check(check, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkRawMarshal(b *testing.B) {
	m := initGoTest(true)
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		WriteDelimited(buf, m)
		buf.Reset()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		WriteDelimited(buf, m)
		if len(buf.Bytes()) != 184 {
			b.Fatalf("unexpected length: %d", len(buf.Bytes()))
		}
		buf.Reset()
	}
}

func BenchmarkReusedMarshal(b *testing.B) {
	m := initGoTest(true)
	bb := new(DelimitedBuffer)
	for i := 0; i < b.N; i++ {
		bb.Marshal(m)
		bb.Clear()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bb.Marshal(m)
		if len(bb.Bytes()) != 184 {
			b.Fatalf("unexpected length: %d", len(bb.Bytes()))
		}
		bb.Clear()
	}
}

func BenchmarkRawUnmarshal(b *testing.B) {
	m := initGoTest(true)
	buf := new(bytes.Buffer)
	WriteDelimited(buf, m)
	out := buf.Bytes()
	d := new(GoTest)
	bufs := make([]*bytes.Buffer, b.N)
	for i := 0; i < b.N; i++ {
		bufs[i] = bytes.NewBuffer(out)
	}
	for i := 0; i < b.N; i++ {
		ReadDelimited(bufs[i], d)
	}
	for i := 0; i < b.N; i++ {
		bufs[i] = bytes.NewBuffer(out)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ReadDelimited(bufs[i], d)
	}
}

func BenchmarkReusedUnmarshal(b *testing.B) {
	m := initGoTest(true)
	buf := new(bytes.Buffer)
	WriteDelimited(buf, m)
	out := buf.Bytes()
	d := new(GoTest)
	dec := new(DelimitedBuffer)
	bufs := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		bufs[i] = out
	}
	for i := 0; i < b.N; i++ {
		dec.SetBuf(bufs[i])
		if _, err := dec.Unmarshal(d); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dec.SetBuf(bufs[i])
		if _, err := dec.Unmarshal(d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReusedUnmarshalRepeating(b *testing.B) {
	m := initGoTest(true)
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		WriteDelimited(buf, m)
	}
	out := buf.Bytes()
	d := new(GoTest)
	dec := new(DelimitedBuffer)
	dec.SetBuf(out)
	deced := 0
outer:
	for {
		_, err := dec.Unmarshal(d)
		deced++
		switch err {
		case io.EOF:
			break outer
		case nil:
			continue
		default:
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	dec.SetBuf(out)
	deced = 0
outer2:
	for {
		_, err := dec.Unmarshal(d)
		deced++
		switch err {
		case io.EOF:
			break outer2
		case nil:
			continue
		default:
			b.Fatal(err)
		}
	}
}
