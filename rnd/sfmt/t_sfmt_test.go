// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sfmt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_sfmt01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sfmt01. integers")

	Init(1234)

	nsamples := 10
	for i := 0; i < nsamples; i++ {
		gen := Rand(0, 10)
		io.Pforan("gen = %v\n", gen)
	}
}
