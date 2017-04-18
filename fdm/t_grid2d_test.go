// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_grid2d_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("grid2d. test 01")

	var g Grid2d
	g.Init(-6.0, 6.0, -3.0, 3.0, 5, 4)

	chk.Int(tst, "N", g.N, 20)
	chk.Int(tst, "Nx", g.Nx, 5)
	chk.Int(tst, "Ny", g.Ny, 4)

	chk.Scalar(tst, "Lx", 1e-15, g.Lx, 12.0)
	chk.Scalar(tst, "Ly", 1e-15, g.Ly, 6.0)
	chk.Scalar(tst, "Dx", 1e-15, g.Dx, 3.0)
	chk.Scalar(tst, "Dy", 1e-15, g.Dy, 2.0)

	chk.Ints(tst, "L", g.L, []int{0, 5, 10, 15})
	chk.Ints(tst, "B", g.B, []int{0, 1, 2, 3, 4})
	chk.Ints(tst, "R", g.R, []int{4, 9, 14, 19})
	chk.Ints(tst, "T", g.T, []int{15, 16, 17, 18, 19})
}

func Test_grid2d_02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("grid2d. test 02")

	var g Grid2d
	g.Init(-6.0, 6.0, -3.0, 3.0, 4, 3)

	dx, dy := 4.0, 3.0
	chk.Scalar(tst, "Dx", 1e-15, g.Dx, dx)
	chk.Scalar(tst, "Dy", 1e-15, g.Dy, dy)

	S := func(v float64) float64 { return v * v }

	Fserial := []float64{
		S(-6) + S(-3), S(-6+dx) + S(-3), S(-6+2*dx) + S(-3), S(-6+3*dx) + S(-3),
		S(-6) + S(-3+dy), S(-6+dx) + S(-3+dy), S(-6+2*dx) + S(-3+dy), S(-6+3*dx) + S(-3+dy),
		S(-6) + S(-3+2*dy), S(-6+dx) + S(-3+2*dy), S(-6+2*dx) + S(-3+2*dy), S(-6+3*dx) + S(-3+2*dy),
	}

	fxy := func(x, y float64) float64 { return x*x + y*y }
	X, Y, F := g.Generate(fxy, nil)
	io.Pforan("X = %v\n", X)
	io.Pforan("Y = %v\n", Y)
	io.Pforan("F = %v\n", F)
	chk.Matrix(tst, "X", 1e-15, X, [][]float64{
		{-6, -6, -6},
		{-6 + dx, -6 + dx, -6 + dx},
		{-6 + 2*dx, -6 + 2*dx, -6 + 2*dx},
		{-6 + 3*dx, -6 + 3*dx, -6 + 3*dx},
	})
	chk.Matrix(tst, "Y", 1e-15, Y, [][]float64{
		{-3, -3 + dy, -3 + 2*dy},
		{-3, -3 + dy, -3 + 2*dy},
		{-3, -3 + dy, -3 + 2*dy},
		{-3, -3 + dy, -3 + 2*dy},
	})
	for i := 0; i < g.Nx; i++ {
		for j := 0; j < g.Ny; j++ {
			chk.Scalar(tst, io.Sf("F%d%d", i, j), 1e-15, F[i][j], Fserial[i+j*g.Nx])
		}
	}
}

func Test_grid2d_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("grid2d. test 03")

	var g Grid2d
	g.Init(2.0, 14.0, 2.0, 8.0, 21, 11)

	if chk.Verbose {
		fxy := func(x, y float64) float64 { return x*x + y*y }
		X, Y, F := g.Generate(fxy, nil)
		plt.SetForPng(0.4, 500, 150, nil)
		plt.ContourF(X, Y, F, nil)
		plt.Equal()
		plt.Gll("x", "y", nil)
		plt.SaveD("/tmp/gosl", "fig_grid2d_03.png")
	}
}
