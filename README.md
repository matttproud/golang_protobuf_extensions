# Overview
This repository provides various Protocol Buffer extensions for the Go
language (golang), namely support for record length-delimited message
streaming.

# Installing

    $ go get github.com/matttproud/golang_protobuf_extensions/ext

# Example

```go
package main

import (
	"io"

	"code.google.com/p/goprotobuf/proto"
	"github.com/matttproud/golang_protobuf_extensions/ext"
)

func main() {
	// You have your pre-populated Protocol Buffer messages.  Yay!
	msgs := []proto.Message{firstMsg, secondMsg, thirdMsg}

	// Destination for writing.
	buf := new(ext.DelimitedBuffer)

	for _, m := range msgs {
		buf.Marshal(m) // Write each out.
		m.Reset()      // Clear the message, since we'll read it back in.
	}

	for _, m := range msgs {
		_, err := buf.Unmarshal(m) // Read each in ...
		if err == io.EOF {         // until we hit EOF or
			break
		} else if err != nil {     // encounter an error.
			panic(err)
		}
	}
}
```

# Documentation
We have [generated Go Doc documentation](
http://godoc.org/github.com/matttproud/golang_protobuf_extensions/ext) here.

# Testing
[![Build Status](https://travis-ci.org/matttproud/golang_protobuf_extensions.png?branch=master)](https://travis-ci.org/matttproud/golang_protobuf_extensions)
