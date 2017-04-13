// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import "github.com/cpmech/gosl/utl"

// Grid2d holds data representing a 2D grid
type Grid2d struct {
	Lx  float64 // length along x
	Ly  float64 // length along y
	Nx  int     // number of divisions along x. Number of spacings = Nx - 1
	Ny  int     // number of divisions along y. Number of spacings = Ny - 1
	N   int     // total number of points
	Dx  float64 // increments along x
	Dy  float64 // increments along y
	Dxx float64 // squared x-increment
	Dyy float64 // squared y-increment
	L   []int   // indices of points along Left edge
	R   []int   // indices of points along Right edge
	B   []int   // indices of points along Bottom edge
	T   []int   // indices of points along Top edge
}

// Init initialises the grid
func (o *Grid2d) Init(lx, ly float64, nx, ny int) {
	o.Lx, o.Ly = lx, ly
	o.Nx, o.Ny, o.N = nx, ny, nx*ny
	o.Dx, o.Dy = o.Lx/float64(nx-1), o.Ly/float64(ny-1)
	o.Dxx, o.Dyy = o.Dx*o.Dx, o.Dy*o.Dy

	o.L = utl.IntRange3(0, o.N, o.Nx)
	o.R = utl.IntAddScalar(o.L, o.Nx-1)
	o.B = utl.IntRange(o.Nx)
	o.T = utl.IntAddScalar(o.B, (o.Ny-1)*o.Nx)
}

// Generate generates coordinates and may evaluate a function over the grid
//   Input:
//     fcn -- function f(x,y) to compute F matrix (may be nil)
//       or
//     Fserial -- serialized f values F[i+j*Nx] (may be nil)
//   Output:
//     X, Y, F(optional) -- matrices of coordinates and f(x,y) values
func (o *Grid2d) Generate(fcn Cb_fxy, Fserial []float64) (X, Y, F [][]float64) {
	X = utl.DblsAlloc(o.Nx, o.Ny)
	Y = utl.DblsAlloc(o.Nx, o.Ny)
	if fcn != nil || Fserial != nil {
		F = utl.DblsAlloc(o.Nx, o.Ny)
	}
	for i := 0; i < o.Nx; i++ {
		x := float64(i) * o.Dx
		for j := 0; j < o.Ny; j++ {
			X[i][j] = x
			Y[i][j] = float64(j) * o.Dy
			if fcn != nil {
				F[i][j] = fcn(X[i][j], Y[i][j])
			}
			if Fserial != nil {
				F[i][j] = Fserial[i+j*o.Nx]
			}
		}
	}
	return
}
