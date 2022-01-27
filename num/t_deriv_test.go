// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
)

func Test_deriv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("deriv01")

	names := []string{
		"x²",       // 1
		"exp(x)",   // 2
		"exp(-x²)", // 3
		"1/x",      // 4
		"x⋅√x",     // 5
		"sin(1/x)", // 6
	}

	fcns := []fun.Ss{
		func(x float64) (float64, error) { // 1
			return x * x, nil
		},
		func(x float64) (float64, error) { // 2
			return math.Exp(x), nil
		},
		func(x float64) (float64, error) { // 3
			return math.Exp(-x * x), nil
		},
		func(x float64) (float64, error) { // 4
			return 1.0 / x, nil
		},
		func(x float64) (float64, error) { // 5
			return x * math.Sqrt(x), nil
		},
		func(x float64) (float64, error) { // 6
			return math.Sin(1.0 / x), nil
		},
	}

	danas := []fun.Ss{
		func(x float64) (float64, error) { // 1
			return 2 * x, nil
		},
		func(x float64) (float64, error) { // 2
			return math.Exp(x), nil
		},
		func(x float64) (float64, error) { // 3
			return -2 * x * math.Exp(-x*x), nil
		},
		func(x float64) (float64, error) { // 4
			return -1.0 / (x * x), nil
		},
		func(x float64) (float64, error) { // 5
			return 1.5 * math.Sqrt(x), nil
		},
		func(x float64) (float64, error) { // 6
			return -math.Cos(1.0/x) / (x * x), nil
		},
	}

	xvals := [][]float64{
		{-0.5, 0, 0.5},   // 1
		{-0.5, 0, 0.5},   // 2
		{-0.5, 0, 0.5},   // 3
		{-0.5, 0.2, 0.5}, // 4
		{0.2, 0.5, 1.0},  // 5
		{-0.5, 0.2, 0.5}, // 6
	}

	//                    1      2      3     4      5      6
	tols := []float64{1e-15, 1e-10, 1e-11, 1e-9, 1e-11, 1e-10}

	//                 1     2     3     4     5     6
	hs := []float64{1e-1, 1e-1, 1e-1, 1e-1, 1e-1, 1e-1}

	// check
	smethods := []string{"DerivCen5", "DerivFwd4", "DerivBwd4"}
	methods := []func(float64, float64, fun.Ss) (float64, error){
		DerivCen5, DerivFwd4, DerivBwd4,
	}
	for idx, method := range methods {
		if idx > 0 {
			tols = []float64{1e-7, 1e-7, 1e-7, 1e-7, 1e-7, 1e-6}
		}
		io.Pf("\ncheck %s:\n", smethods[idx])
		for i, f := range fcns {
			for _, x := range xvals[i] {
				dnum, err := method(x, hs[i], f)
				if err != nil {
					tst.Errorf("%v\n", err)
					return
				}
				dana, err := danas[i](x)
				if err != nil {
					tst.Errorf("%v\n", err)
					return
				}
				chk.Scalar(tst, "    "+names[i], tols[i], dnum, dana)
			}
		}
	}
}

func Test_deriv02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("deriv02")

	// scalar field
	fcn := func(x, y float64) (float64, error) {
		return -math.Pow(math.Pow(math.Cos(x), 2.0)+math.Pow(math.Cos(y), 2.0), 2.0), nil
	}

	// gradient. u=dfdx, v=dfdy
	grad := func(x, y float64) (u, v float64) {
		m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
		u = 4.0 * math.Cos(x) * math.Sin(x) * m
		v = 4.0 * math.Cos(y) * math.Sin(y) * m
		return
	}

	// grid size
	xmin, xmax, N := -math.Pi/2.0, math.Pi/2.0, 7
	dx := (xmax - xmin) / float64(N-1)

	// step size for numerical differentiation
	h := 0.01

	// tolerance
	tol := 1e-9

	// loop over grid points
	for i := 0; i < N; i++ {
		x := xmin + float64(i)*dx
		for j := 0; j < N; j++ {
			y := xmin + float64(i)*dx

			// scalar and vector field
			f, _ := fcn(x, y)
			u, v := grad(x, y)

			// numerical dfdx @ (x,y)
			unum, err := DerivCen5(x, h, func(xvar float64) (float64, error) {
				return fcn(xvar, y)
			})
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}

			// numerical dfdy @ (x,y)
			vnum, err := DerivCen5(y, h, func(yvar float64) (float64, error) {
				return fcn(x, yvar)
			})
			if err != nil {
				tst.Errorf("%v\n", err)
				return
			}

			// output
			if chk.Verbose {
				io.Pforan("x=%+.3f y=%+.3f f=%+.3f u=%+.3f v=%+.3f unum=%+.3f vnum=%+.3f\n", x, y, f, u, v, unum, vnum)
			}

			// check
			if math.Abs(unum-u) > tol {
				tst.Errorf("x=%v y=%v f=%v u=%v v=%v unum=%v vnum=%v error=%v\n", x, y, f, u, v, unum, vnum, math.Abs(unum-u))
				return
			}
			if math.Abs(vnum-v) > tol {
				tst.Errorf("x=%v y=%v f=%v u=%v v=%v unum=%v vnum=%v error=%v\n", x, y, f, u, v, unum, vnum, math.Abs(vnum-v))
				return
			}
		}
	}
}
