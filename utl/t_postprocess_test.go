// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_postp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("postp01")

	Tout := []float64{0, 0.1, 0.2, 0.200001, 0.201, 0.3001, 0.8, 0.99, 0.999, 1}
	Tsel := []float64{0, 0.2, 0.3, 0.6, 0.8, 0.9, 0.99, -1}

	tol := 0.001
	I, T := GetITout(Tout, Tsel, tol)
	io.Pfcyan("Tout = %v\n", Tout)
	io.Pfcyan("Tsel = %v\n", Tsel)
	io.Pforan("I = %v\n", I)
	io.Pforan("T = %v\n", T)

	chk.Ints(tst, "I", I, []int{0, 2, 5, 6, 7, 9})
	chk.Vector(tst, "T", 1e-16, T, []float64{0, 0.2, 0.3001, 0.8, 0.99, 1})
}
