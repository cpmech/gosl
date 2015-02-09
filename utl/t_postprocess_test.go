// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"
)

func Test_postp01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("postp01")

	Tout := []float64{0, 0.1, 0.2, 0.200001, 0.201, 0.3001, 0.8, 0.99, 0.999, 1}
	Tsel := []float64{0, 0.2, 0.3, 0.6, 0.8, 0.9, 0.99, -1}

	tol := 0.001
	I, T := GetITout(Tout, Tsel, tol)
	Pfcyan("Tout = %v\n", Tout)
	Pfcyan("Tsel = %v\n", Tsel)
	Pforan("I = %v\n", I)
	Pforan("T = %v\n", T)

	CompareInts(tst, "I", I, []int{0, 2, 5, 6, 7, 9})
	CompareDbls(tst, "T", T, []float64{0, 0.2, 0.3001, 0.8, 0.99, 1})
}
