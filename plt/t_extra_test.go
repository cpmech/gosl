// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func TestWaterfall01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Waterfall01")

	if chk.Verbose {

		// a simple (un-normalized) Gaussian shape with amplitude A.
		fG := func(x, x0, σ, A float64) float64 {
			return A * math.Exp(-math.Pow((x-x0)/σ, 2.0))
		}

		rand.Seed(int64(time.Now().Unix()))

		σ := 0.05
		nx, nt := 401, 11
		X := utl.NonlinSpace(0, 2, nx, 4.0, true)
		T := utl.LinSpace(0, 1, nt)
		Z := utl.Alloc(nt, nx)
		for i := 0; i < nt; i++ {
			for k := 0; k < 4; k++ {
				x0 := rand.Float64() * 2
				A := rand.Float64() * 10.0
				for j := 0; j < nx; j++ {
					Z[i][j] += fG(X[j], x0, σ, A)
				}
			}
		}

		Reset(true, nil)
		Waterfall(X, T, Z, nil)
		err := Save("/tmp/gosl/plt", "t_waterfall01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func TestSlopeInd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SlopeInd01. slope indicator")

	if chk.Verbose {
		x := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8}
		y := []float64{0, 2, 4, 6, 8, 6, 4, 2, 0}
		Reset(true, &A{Prop: 1.5})
		draw := func(idx int) {
			Subplot(2, 1, idx)
			Plot(x, y, &A{C: C(0, 0), M: ".", NoClip: true})
			SlopeInd(+2.0, 1.1, 2.0, 1, "2", false, false, false, nil, nil)
			SlopeInd(+2.0, 2.9, 6.0, 1, "2", true, false, false, nil, nil)
			SlopeInd(-2.0, 7.1, 2.0, 1, "2", false, false, false, nil, nil)
			SlopeInd(-2.0, 5.4, 5.0, 1, "2", true, false, false, nil, nil)
			Grid(nil)
			if idx == 1 {
				Title("with equal scale", &A{Fsz: 8})
				Equal()
			} else {
				Title("unscaled", &A{Fsz: 8})
			}
		}
		draw(1)
		draw(2)
		err := Save("/tmp/gosl/plt", "t_slopeind01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func TestSlopeInd02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SlopeInd02. slope indicator")

	if chk.Verbose {
		x := []float64{1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6}
		y := []float64{1, 1e2, 1e4, 1e6, 1e4, 1e2, 1}
		Reset(true, &A{Prop: 1.5})
		draw := func(idx int) {
			Subplot(2, 1, idx)
			Plot(x, y, &A{C: C(0, 0), M: ".", NoClip: true})
			SlopeInd(+2.0, 1.1, 2.0, 1, "2", false, true, true, nil, nil)
			SlopeInd(+2.0, 1.9, 4.0, 1, "2", true, true, true, nil, nil)
			SlopeInd(-2.0, 5.1, 2.0, 1, "2", false, true, true, nil, nil)
			SlopeInd(-2.0, 3.9, 4.0, 1, "2", true, true, true, nil, nil)
			Grid(nil)
			SetXlog()
			SetYlog()
			if idx == 1 {
				Title("with equal scale", &A{Fsz: 8})
				Equal()
			} else {
				Title("unscaled", &A{Fsz: 8})
			}
		}
		draw(1)
		draw(2)
		err := Save("/tmp/gosl/plt", "t_slopeind02")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func TestSlopeInd03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SlopeInd03. slope indicator")

	if chk.Verbose {
		x := []float64{1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6}
		y := []float64{0, 2, 4, 6, 4, 2, 0}
		Reset(true, nil)
		Plot(x, y, &A{C: C(0, 0), M: ".", NoClip: true})
		SlopeInd(+2.0, 1.1, 2.0, 1, "2", false, true, false, nil, nil)
		SlopeInd(+2.0, 1.9, 4.0, 1, "2", true, true, false, nil, nil)
		SlopeInd(-2.0, 5.1, 2.0, 1, "2", false, true, false, nil, nil)
		SlopeInd(-2.0, 3.9, 4.0, 1, "2", true, true, false, nil, nil)
		Grid(nil)
		SetXlog()
		err := Save("/tmp/gosl/plt", "t_slopeind03")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func TestSlopeInd04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SlopeInd04. slope indicator")

	if chk.Verbose {
		x := []float64{0, 1, 2, 3, 4, 5, 6}
		y := []float64{1, 1e2, 1e4, 1e6, 1e4, 1e2, 1}
		Reset(true, nil)
		Plot(x, y, &A{C: C(0, 0), M: ".", NoClip: true})
		SlopeInd(+2.0, 1.1, 2.0, 1, "2", false, false, true, nil, nil)
		SlopeInd(+2.0, 1.9, 4.0, 1, "2", true, false, true, nil, nil)
		SlopeInd(-2.0, 5.1, 2.0, 1, "2", false, false, true, nil, nil)
		SlopeInd(-2.0, 3.9, 4.0, 1, "2", true, false, true, nil, nil)
		Grid(nil)
		SetYlog()
		err := Save("/tmp/gosl/plt", "t_slopeind04")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func TestMatrix01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Matrix01. Subplot matrix")

	if chk.Verbose {
		xx := [][]float64{ // [iComponent][iTime]
			{0.0, +1.0, +2.0, +3.0}, // x0
			{0.0, 1.0, 4.0, 9.0},    // x1
		}
		count := 0
		cmds := func(i, j int) {
			io.Pforan("\nxx[%d] = %v\n", i, xx[i])
			io.Pf("xx[%d] = %v\n", j, xx[j])
			Plot(xx[i], xx[j], &A{M: ".", C: C(count, 0)})
			count++
		}
		Reset(true, &A{WidthPt: 500})
		SubplotMatrix(len(xx), len(xx), cmds)
		err := Save("/tmp/gosl/plt", "t_matrix01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func TestMatrixSym01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MatrixSym01. Subplot matrix symmetric")

	if chk.Verbose {
		xx := [][]float64{ // [iComponent][iTime]
			{0.0, +1.0, +2.0, +3.0}, // x0
			{0.0, -1.0, -2.0, -3.0}, // x1
			{0.0, 2.0, 4.0, 6.0},    // x2
			{0.0, 1.0, 4.0, 9.0},    // x3
		}
		count := 0
		cmds := func(i, j int) {
			io.Pforan("\nxx[%d] = %v\n", i, xx[i])
			io.Pf("xx[%d] = %v\n", j, xx[j])
			Plot(xx[i], xx[j], &A{M: ".", C: C(count, 0)})
			count++
		}
		Reset(true, &A{WidthPt: 500})
		SubplotMatrixSym(len(xx), len(xx), cmds, nil)
		err := Save("/tmp/gosl/plt", "t_matrixsym01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
