// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestQags01(tst *testing.T) {

	verbose()
	chk.PrintTitle("Qags01.")

	/*
		y := func(x float64) (res float64, err error) {
			res = math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
			return
		}
		var err error
		Acor := 1.08268158558

		A, err := QuadGaussL10(0, 1, y)
		if err != nil {
			tst.Errorf("%v\n", err)
			return
		}
		io.Pforan("A  = %v\n", A)
		chk.Scalar(tst, "A", 1e-12, A, Acor)
	*/
}

/*
func Test_gaussLegXW01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("gaussLegXW01. Gauss-Legendre x-w data.")

	// constants
	xRef := []float64{-0.9739065285171717, -0.8650633666889845, -0.6794095682990244, -0.4333953941292472, -0.1488743389816312, 0.1488743389816312, 0.4333953941292472, 0.6794095682990244, 0.8650633666889845, 0.9739065285171717}
	wRef := []float64{0.0666713443086881, 0.1494513491505806, 0.2190863625159821, 0.2692667193099963, 0.2955242247147529, 0.2955242247147529, 0.2692667193099963, 0.2190863625159821, 0.1494513491505806, 0.0666713443086881}

	xL, wL := GaussLegendreXW(-1, 1, 10)
	chk.Vector(tst, "xL", 1e-15, xL, xRef)
	chk.Vector(tst, "wL", 1e-15, wL, wRef)

	xJ, wJ := GaussJacobiXW(0, 0, 10)
	chk.Vector(tst, "xJ", 1e-15, xJ, xRef)
	chk.Vector(tst, "wJ", 1e-14, wJ, wRef)
}
*/
