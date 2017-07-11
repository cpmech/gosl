// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"

	goio "io"
)

type testWriter []byte
type testReader []byte
type testStruct struct{ A, B int }

func (o *testWriter) Write(p []byte) (n int, err error) {
	(*o) = p
	n = len(p)
	return
}

func (o testReader) Read(p []byte) (n int, err error) {
	io.Pforan("p = %v\n", string(p))
	return 0, goio.EOF
}

func Test_encoder01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("encoder01")

	// writer and encoder
	var w testWriter
	enc := GetEncoder(&w, "json")
	sIn := testStruct{A: 1, B: 2}
	enc.Encode(&sIn)
	io.Pforan("w = %q\n", string(w))
	chk.String(tst, "{\"A\":1,\"B\":2}\n", string(w))

	// reader and decoder
	var r testReader
	dec := GetDecoder(&r, "json")
	var sOut testStruct
	dec.Decode(&sOut)
	io.Pforan("r = %v\n", r)
	// TODO: implement test here
}
