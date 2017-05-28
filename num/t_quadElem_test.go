// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_QuadElem01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("QuadElem01. Trapz and Simpson Elementary")

	y := func(x float64) (res float64, err error) {
		res = math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
		return
	}
	var err error
	Acor := 1.08268158558

	// trapezoidal rule
	var T QuadElementary
	T = new(ElementaryTrapz)
	T.Init(y, 0, 1, 1e-11)
	A, err := T.Integrate()
	if err != nil {
		io.Pforan(err.Error())
	}
	io.Pforan("A  = %v\n", A)
	chk.Scalar(tst, "A", 1e-11, A, Acor)

	// Simpson's rule
	var S QuadElementary
	S = new(ElementarySimpson)
	S.Init(y, 0, 1, 1e-11)
	A, err = S.Integrate()
	if err != nil {
		io.Pforan(err.Error())
	}
	io.Pforan("A  = %v\n", A)
	chk.Scalar(tst, "A", 1e-11, A, Acor)
}
