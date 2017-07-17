// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestNonlinSpace01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NonlinSpace01")

	xa, xb := 0.0, 1.0
	N := 11
	R := 0.1
	x := NonlinSpace(xa, xb, N, R, false)
	if x[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x[0], xa)
	}
	if x[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x[N-1], xb)
		return
	}

	chk.Int(tst, "npts", len(x), N)

	Δx0 := x[1] - x[0]
	Δxn := x[len(x)-1] - x[len(x)-2]

	io.Pf("Δxn/Δx0 = %v\n", Δxn/Δx0)
	chk.Scalar(tst, "Δxn/Δx0", 1e-15, Δxn/Δx0, R)

	R = 1
	xuni := NonlinSpace(xa, xb, N, R, false)
	if xuni[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", xuni[0], xa)
	}
	if xuni[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", xuni[N-1], xb)
		return
	}
	chk.Array(tst, "xuni", 1e-17, xuni, LinSpace(xa, xb, N))

	N = 2
	x2 := NonlinSpace(xa, xb, N, R, false)
	chk.Array(tst, "x2", 1e-17, x2, []float64{xa, xb})

	N = 1
	x2b := NonlinSpace(xa, xb, N, R, false)
	chk.Array(tst, "x2b", 1e-17, x2b, []float64{xa, xb})

	xa, xb = -2.0, 12.0
	N = 4
	R = 4
	x4 := NonlinSpace(xa, xb, N, R, false)
	if x4[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x4[0], xa)
	}
	if x4[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x4[N-1], xb)
		return
	}
	chk.Array(tst, "x4", 1e-17, x4, []float64{-2, 0, 4, 12})
}

func TestNonlinSpace02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NonlinSpace02")

	xa, xb := 0.0, 10.0
	R := 4.0

	N := 2
	x2 := NonlinSpace(xa, xb, N, R, true)
	chk.Array(tst, "x2", 1e-17, x2, []float64{0, 10})
	if x2[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x2[0], xa)
	}
	if x2[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x2[N-1], xb)
		return
	}

	N = 3
	x3 := NonlinSpace(xa, xb, N, R, true)
	chk.Array(tst, "x3", 1e-17, x3, []float64{0, 5, 10})
	if x3[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x3[0], xa)
	}
	if x3[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x3[N-1], xb)
		return
	}

	N = 4
	x4 := NonlinSpace(xa, xb, N, R, true)
	chk.Array(tst, "x4", 1e-14, x4, []float64{0, 5.0 / 3.0, 25.0 / 3.0, 10})
	if x4[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x4[0], xa)
	}
	if x4[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x4[N-1], xb)
		return
	}

	N = 5
	x5 := NonlinSpace(xa, xb, N, R, true)
	chk.Array(tst, "x5", 1e-17, x5, []float64{0, 1, 5, 9, 10})
	if x5[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x5[0], xa)
	}
	if x5[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x5[N-1], xb)
		return
	}

	N = 6
	x6 := NonlinSpace(xa, xb, N, R, true)
	chk.Array(tst, "x6", 1e-17, x6, []float64{0, 1, 3, 7, 9, 10})
	if x6[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x6[0], xa)
	}
	if x6[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x6[N-1], xb)
		return
	}

	N = 7
	x7 := NonlinSpace(xa, xb, N, R, true)
	chk.Array(tst, "x5", 1e-15, x7, []float64{0, 10.0 / 14.0, 30.0 / 14.0, 70.0 / 14.0, 110.0 / 14.0, 130.0 / 14.0, 10.0})
	if x7[0] != xa {
		tst.Errorf("first point must be exact. %g != %g\n", x7[0], xa)
	}
	if x7[N-1] != xb {
		tst.Errorf("last point must be exact. %g != %g\n", x7[N-1], xb)
		return
	}
}

func TestNonlinSpace03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NonlinSpace03")

	xa, xb := -3.0, 27.0
	R := 4.0

	for _, N := range IntRange2(2, 12) {
		x := NonlinSpace(xa, xb, N, R, true)
		io.Pf("\n-----------------------------------------------------\n")
		io.Pf("N = %d\n", N)
		io.Pf("x = %v\n", x)
		chk.Int(tst, "len(x)", len(x), N)
		chk.Scalar(tst, "x0", 1e-17, x[0], xa)
		chk.Scalar(tst, "xL", 1e-17, x[N-1], xb)

		// some constants
		l := N - 1
		even := l%2 == 0
		imax := l/2 + 1
		if !even {
			imax = (l + 1) / 2
		}

		// check symmetry
		tolSym := 1e-17
		if N > 6 {
			tolSym = 1e-14
		}
		for i := 1; i < imax; i++ {
			io.Pforan("%d → %d  vs  %d → %d\n", i-1, i, l-i, l-i+1)
			Δxa := x[i] - x[i-1]
			Δxb := x[l-i+1] - x[l-i]
			chk.Scalar(tst, "Δxa = Δxb", tolSym, Δxa, Δxb)
		}

		// check ratio
		if N < 4 {
			continue
		}
		tolR := 1e-17
		if N > 6 {
			tolR = 1e-15
		}
		if N > 10 {
			tolR = 1e-14
		}
		Δxl := x[imax] - x[imax-1]
		Δx0 := x[1] - x[0]
		r := Δxl / Δx0
		io.Pforan("r = %v\n", r)
		chk.Scalar(tst, "ratio", tolR, r, R)
	}
}
