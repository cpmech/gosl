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

type test_writer []byte
type test_reader []byte
type test_struct struct{ A, B int }

func (o *test_writer) Write(p []byte) (n int, err error) {
	(*o) = p
	n = len(p)
	return
}

func (o test_reader) Read(p []byte) (n int, err error) {
	io.Pforan("p = %v\n", string(p))
	return 0, goio.EOF
}

func Test_encoder01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("encoder01")

	// writer and encoder
	var w test_writer
	enc := GetEncoder(&w, "json")
	s_in := test_struct{A: 1, B: 2}
	enc.Encode(&s_in)
	io.Pforan("w = %q\n", string(w))
	chk.String(tst, "{\"A\":1,\"B\":2}\n", string(w))

	// reader and decoder
	var r test_reader
	dec := GetDecoder(&r, "json")
	var s_out test_struct
	dec.Decode(&s_out)
	io.Pforan("r = %v\n", r)
	// TODO: implement test here
}
