// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bytes"
	"testing"

	"gosl/chk"
)

func TestEncDec01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EncDec01")

	// data
	dat := struct{ A, B int }{123, 456}

	// ------- A -------

	// writer
	wbuf := new(bytes.Buffer)

	// encoder
	enc := NewEncoder(wbuf, "json")

	// encode
	enc.Encode(&dat)
	chk.String(tst, string(wbuf.Bytes()), `{"A":123,"B":456}`+"\n")

	// ------- B -------

	// reader
	rbuf := bytes.NewBuffer(wbuf.Bytes())

	// decoder
	dec := NewDecoder(rbuf, "json")

	// decode
	res := struct{ A, B int }{}
	dec.Decode(&res)
	if res.A != 123 {
		tst.Errorf("A should be 123")
		return
	}
	if res.B != 456 {
		tst.Errorf("A should be 456")
		return
	}
}

func TestEncDec02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("EncDec02")

	// data
	dat := struct{ A, B int }{123, 456}

	// ------- A -------

	// writer
	wbuf := new(bytes.Buffer)

	// encoder
	enc := NewEncoder(wbuf, "gob")

	// encode
	enc.Encode(&dat)
	// cannot check wbuf

	// ------- B -------

	// reader
	rbuf := bytes.NewBuffer(wbuf.Bytes())

	// decoder
	dec := NewDecoder(rbuf, "gob")

	// decode
	res := struct{ A, B int }{}
	dec.Decode(&res)
	if res.A != 123 {
		tst.Errorf("A should be 123")
		return
	}
	if res.B != 456 {
		tst.Errorf("A should be 456")
		return
	}
}
