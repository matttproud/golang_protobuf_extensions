# Overview
This repository provides various Protocol Buffer extensions for the Go
language (golang), namely support for record length-delimited message
streaming.

| Java                           | Go                 |
| ------------------------------ | ------------------ |
| MessageLite#parseDelimitedFrom | ext.ReadDelimited  |
| MessageLite#writeDelimitedTo   | ext.WriteDelimited |

Because [Code Review 9102043](https://codereview.appspot.com/9102043/) is
destined to never be merged into mainline (i.e., never be promoted to formal
[goprotobuf features](https://code.google.com/p/goprotobuf)), this repository
will live here in the wild.

# Documentation
We have [generated Go Doc documentation](http://godoc.org/github.com/matttproud/golang_protobuf_extensions/ext) here.

# Testing
[![Build Status](https://travis-ci.org/matttproud/golang_protobuf_extensions.png?branch=master)](https://travis-ci.org/matttproud/golang_protobuf_extensions)
