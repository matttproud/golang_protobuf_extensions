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
