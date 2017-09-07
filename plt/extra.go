// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// Waterfall draws parallel lines @ t along x with height = z. z[len(t)][len(x)]
func Waterfall(X, T []float64, Z [][]float64, args *A) {
	if args == nil {
		args = &A{Fc: "w", Ec: "k", Closed: false}
	}
	createAxes3d()
	uid := genUID()
	sx := io.Sf("X%d", uid)
	sz := io.Sf("Z%d", uid)
	genArray(&bufferPy, sx, X)
	genMat(&bufferPy, sz, Z)
	nx := len(X)
	nt := len(T)
	tt := make([]float64, nx)
	P := utl.Alloc(nx, 3)
	xmin, xmax, tmin, tmax, zmin, zmax := X[0], X[0], T[0], T[0], Z[0][0], Z[0][0]
	for i := nt - 1; i >= 0; i-- {
		t := T[i]
		utl.Fill(tt, t)
		uid = genUID()
		st := io.Sf("T%d", uid)
		genArray(&bufferPy, st, tt)
		for j, x := range X {
			P[j][0] = x
			P[j][1] = t
			P[j][2] = Z[i][j]
			zmin = utl.Min(zmin, Z[i][j])
			zmax = utl.Max(zmax, Z[i][j])
		}
		tmin = utl.Min(tmin, t)
		tmax = utl.Max(tmax, t)
		Polygon3d(P, args)
	}
	for _, x := range X {
		xmin = utl.Min(xmin, x)
		xmax = utl.Max(xmax, x)
	}
	AxisRange3d(xmin, xmax, tmin, tmax, zmin, zmax)
}

// SlopeInd draws indicator of line slope
func SlopeInd(m, xc, yc, xlen float64, lbl string, flip, xlog, ylog bool, args, argsLbl *A) {
	if args == nil {
		args = &A{C: "k"}
	}
	args.NoClip = true
	l := 0.5 * xlen
	x := []float64{xc - l, xc + l, xc + l, xc - l}
	y := []float64{yc - m*l, yc - m*l, yc + m*l, yc - m*l}
	if flip {
		x[1] = xc - l
		y[1] = yc + m*l
	}
	dx, dy := x[2]-x[0], y[2]-y[0]
	d := 0.03 * math.Sqrt(dx*dx+dy*dy)
	xm := xc - l - d
	xp := xc + l + d
	ym := yc + m*l - d
	yp := yc + m*l + d
	yr := yc - m*l + d
	ys := yc - m*l - d
	if xlog {
		for i := 0; i < 4; i++ {
			x[i] = math.Pow(10.0, x[i])
		}
		xc = math.Pow(10.0, xc)
		xm = math.Pow(10.0, xm)
		xp = math.Pow(10.0, xp)
	}
	if ylog {
		for i := 0; i < 4; i++ {
			y[i] = math.Pow(10.0, y[i])
		}
		yc = math.Pow(10.0, yc)
		ym = math.Pow(10.0, ym)
		yp = math.Pow(10.0, yp)
		yr = math.Pow(10.0, yr)
		ys = math.Pow(10.0, ys)
	}
	Plot(x, y, args)
	if lbl != "" {
		if argsLbl == nil {
			argsLbl = &A{C: "k", Fsz: 6}
		}
		argsLbl.NoClip = true
		if flip {
			argsLbl.Ha = "center"
			if m < 0 {
				argsLbl.Va = "top"
				Text(xc, ym, "1", argsLbl)
			} else {
				argsLbl.Va = "bottom"
				Text(xc, yp, "1", argsLbl)
			}
			argsLbl.Ha = "right"
			argsLbl.Va = "center"
			Text(xm, yc, lbl, argsLbl)
		} else {
			argsLbl.Ha = "center"
			if m < 0 {
				argsLbl.Va = "bottom"
				Text(xc, yr, "1", argsLbl)
			} else {
				argsLbl.Va = "top"
				Text(xc, ys, "1", argsLbl)
			}
			argsLbl.Ha = "left"
			argsLbl.Va = "center"
			Text(xp, yc, lbl, argsLbl)
		}
	}
}

// SubplotMatrix call plotCommands to fill a matrix of plots
//   Input:
//     nrow -- number of rows in matrix [must be at least 1]
//     ncol -- number of columns in matrix [must be at least 1]
//     cmds -- plotting commants for each (i,j) grid point (NOT Subplot indices)
func SubplotMatrix(nrow, ncol int, cmds func(i, j int)) {
	if nrow < 1 {
		chk.Panic("At least one row is required. nrow = %d is invalid\n", nrow)
	}
	if ncol < 1 {
		chk.Panic("At least one column is required. ncol = %d is invalid\n", ncol)
	}
	idx := nrow * ncol
	for row := 0; row < nrow; row++ {
		for col := ncol - 1; col >= 0; col-- {
			Subplot(nrow, ncol, idx)
			SplotGap(0.0, 0.0)
			cmds(col, row)
			if col > 0 {
				SetNoYtickLabels()
			} else {
				SetYlabel(io.Sf("$x_{%d}$", row), nil)
			}
			if row > 0 {
				SetNoXtickLabels()
			} else {
				SetXlabel(io.Sf("$x_{%d}$", col), nil)
			}
			Grid(&A{C: "grey", Z: -1000})
			idx--
		}
	}
}

// SubplotMatrixSym call plotCommands to fill a symmetric matrix of plots (excluding the diagonal)
//   Input:
//     nrow   -- number of rows in matrix including diagonal [must be at least 2]
//     ncol   -- number of columns in matrix  including diagonal[must be at least 2]
//     cmds   -- plotting commants for each (i,j) grid point (NOT Subplot indices)
//     corner -- function to be called to draw on "corner" of matrix (left-bottom side)
//   Output:
//     n -- the effective number of columns == min(nrow-1,ncol-1)
//     cornerIdx -- the subplot index of "corner" graph
//     NOTE: the resulting grid will be (n×n)
func SubplotMatrixSym(nrow, ncol int, cmds func(i, j int), corner func()) (n, cornerIdx int) {
	if nrow < 2 {
		chk.Panic("At least two rows are required. nrow = %d is invalid\n", nrow)
	}
	if ncol < 2 {
		chk.Panic("At least two columns are required. ncol = %d is invalid\n", ncol)
	}
	n = utl.Imin(nrow, ncol)
	n--
	idx := 1
	for row := 0; row < nrow; row++ {
		idx += row
		for col := row + 1; col < nrow; col++ {
			Subplot(n, n, idx)
			SplotGap(0.0, 0.0)
			cmds(col, row)
			if col > row+1 {
				SetNoXtickLabels()
				SetNoYtickLabels()
				Grid(&A{C: "grey", Z: -1000})
			} else {
				Gll(io.Sf("$x_{%d}$", col), io.Sf("$x_{%d}$", row), nil)
			}
			idx++
		}
	}
	cornerIdx = n*(n-1) + 1
	Subplot(n, n, cornerIdx)
	AxisOff()
	if corner != nil {
		corner()
	}
	return
}
