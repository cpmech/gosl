// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

import (
	. "math"

	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestSpecProb01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SpecProb01. Problem with cos(n⋅x)⋅exp(cos(x)-1)")

	π := Pi
	for i := 60; i < 100; i++ {
		n := float64(i)
		f1 := func(x float64) float64 { return Cos(n*x) * Exp(Cos(x)-1) }
		f2 := func(x float64) float64 { return Exp(Cos(x) - 1) }
		A1, _, _, _, err1 := Agse(0, f1, -π, π, 0, 0, nil, nil, nil, nil, nil)
		A2, _, _, _, err2 := Awoe(0, f2, -π, π, n, 1, 0, 0, 0, 0, nil, nil, nil, nil, nil, nil, 0, nil)
		io.Pf("n=%2.f  A1=%23.15e A2=%23.15e diff=%g\n", n, A1, A2, A1-A2)
		if err1 != nil {
			if n > 75 {
				io.Pf("YES: it fails ⇒ %v", err1)
			} else {
				tst.Errorf("%v", err1)
				return
			}
		}
		if err2 != nil {
			tst.Errorf("%v\n", err2)
			return
		}
	}
}
