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
	"log"
	"sync"
)

func deprWarning(n string) func() {
	return func() {
		log.Printf("function %s is slated for deletion and will go away June, 20 2014", n)
	}
}

var writeWarn = deprWarning("WriteDelimited")
var readWarn = deprWarning("ReadDelimited")

var readOnce, writeOnce sync.Once

func deprWriteDelimited() {
	writeOnce.Do(writeWarn)
}

func deprReadDelimited() {
	readOnce.Do(readWarn)
}
