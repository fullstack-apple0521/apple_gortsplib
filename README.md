
# gortsplib

[![Build Status](https://travis-ci.org/aler9/gortsplib.svg?branch=master)](https://travis-ci.org/aler9/gortsplib)
[![Go Report Card](https://goreportcard.com/badge/github.com/aler9/gortsplib)](https://goreportcard.com/report/github.com/aler9/gortsplib)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/aler9/gortsplib?tab=doc)

RTSP 1.0 library for the Go programming language, written for [rtsp-simple-server](https://github.com/aler9/rtsp-simple-server).

Features:
* Read streams via TCP or UDP
* Publish streams via TCP or UDP
* Provides primitives, a class for building clients (`ConnClient`) and a class for building servers (`ConnServer`)

## Examples

* [read-tcp](examples/read-tcp.go)
* [read-udp](examples/read-udp.go)
* [publish-tcp](examples/publish-tcp.go)
* [publish-udp](examples/publish-udp.go)

## Documentation

https://pkg.go.dev/github.com/aler9/gortsplib

## Links

Related projects
* https://github.com/aler9/rtsp-simple-server
* https://github.com/pion/sdp (SDP library used internally)
* https://github.com/pion/rtcp (RTCP library used internally)

IETF Standards
* RTSP 1.0 https://tools.ietf.org/html/rfc2326
* RTSP 2.0 https://tools.ietf.org/html/rfc7826
* HTTP 1.1 https://tools.ietf.org/html/rfc2616
