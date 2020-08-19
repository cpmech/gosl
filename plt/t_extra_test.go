// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
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
		Save("/tmp/gosl/plt", "t_waterfall01")
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
		Save("/tmp/gosl/plt", "t_slopeind01")
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
		Save("/tmp/gosl/plt", "t_slopeind02")
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
		Save("/tmp/gosl/plt", "t_slopeind03")
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
		Save("/tmp/gosl/plt", "t_slopeind04")
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
		Save("/tmp/gosl/plt", "t_matrix01")
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
		Save("/tmp/gosl/plt", "t_matrixsym01")
	}
}

func TestDrawArrow2d(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DrawArrow2d")

	if chk.Verbose {
		Reset(true, &A{WidthPt: 500})
		Circle(0, 0, 1, &A{C: "grey", Lw: 0.8})
		DrawArrow2d([]float64{0, 0}, []float64{123, 0}, true, 1, nil)
		DrawArrow2d([]float64{0, 0}, []float64{66, 66}, true, 1, nil)
		DrawArrow2d([]float64{0, 0}, []float64{0, 88}, true, 1, nil)
		Equal()
		Grid(nil)
		Save("/tmp/gosl/plt", "t_arrow2d")
	}
}

func TestDrawArrow3d(tst *testing.T) {

	//verbose()
	chk.PrintTitle("DrawArrow3d")

	if chk.Verbose {
		Reset(true, &A{WidthPt: 500})
		DrawArrow3d([]float64{0, 0, 0}, []float64{12, 0, 0}, true, 1, nil)
		DrawArrow3d([]float64{0, 0, 0}, []float64{66, 66, 66}, true, 1, nil)
		DrawArrow3d([]float64{0, 0, 0}, []float64{0, 88, 0}, true, 1, nil)
		DrawArrow3d([]float64{0, 0, 0}, []float64{0, 0, 88}, true, 1, nil)
		Camera(40, 30, nil)
		Equal()
		Save("/tmp/gosl/plt", "t_arrow3d")
	}
}
