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

func Test_quadGaussL01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quadGaussL01. Gauss-Legendre quadrature.")

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
}
