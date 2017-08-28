// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/utl"
)

// Grid implements a 2D or 3D grid of points (based on Bins)
type Grid struct {
	Bins // derived

	N    int     // total number of points
	Edge [][]int // ids of points on edges: [edge0, edge1, edge2, edge3]
	Face [][]int // ids of points on faces: [face0, face1, face2, face3, face4, face5]
}

// Init initialise Grid structure
//   xmin -- [ndim] min/initial coordinates of the whole space (box/cube)
//   xmax -- [ndim] max/final coordinates of the whole space (box/cube)
//   ndiv -- [ndim] number of divisions for xmax-xmin
func (o *Grid) Init(xmin, xmax []float64, ndiv []int) (err error) {
	err = o.Bins.Init(xmin, xmax, ndiv)
	o.N = 1
	for k := 0; k < o.Ndim; k++ {
		o.N *= o.Npts[k]
	}
	nx := o.Npts[0]
	ny := o.Npts[1]
	if o.Ndim == 2 {
		o.Edge = make([][]int, 4)   // bottom, right, top, left
		o.Edge[0] = make([]int, nx) // bottom
		o.Edge[1] = make([]int, ny) // right
		o.Edge[2] = make([]int, nx) // top
		o.Edge[3] = make([]int, ny) // left
		for i := 0; i < nx; i++ {
			o.Edge[0][i] = i             // bottom
			o.Edge[2][i] = i + (ny-1)*nx // top
		}
		for j := 0; j < ny; j++ {
			o.Edge[1][j] = j*nx + nx - 1 // right
			o.Edge[3][j] = j * nx        // left
		}
		return
	}
	nz := o.Npts[2]
	o.Face = make([][]int, 6)      // xmin,xmax, ymin,ymax, zmin,zmax
	o.Face[0] = make([]int, ny*nz) // xmin
	o.Face[1] = make([]int, ny*nz) // xmax
	o.Face[2] = make([]int, nx*nz) // ymin
	o.Face[3] = make([]int, nx*nz) // ymax
	o.Face[4] = make([]int, nx*ny) // zmin
	o.Face[5] = make([]int, nx*ny) // zmax
	// TODO: implement ids of faces
	return
}

// Eval2d evaluates function over grid
//  X -- [ny][nx] will hold the grid coordinates
//  Y -- [ny][nx] will hold the grid coordinates
//  F -- [ny][nx] will hold the results
func (o *Grid) Eval2d(f fun.Sv) (X, Y, F [][]float64, err error) {
	nx := o.Npts[0]
	ny := o.Npts[1]
	X = utl.Alloc(ny, nx)
	Y = utl.Alloc(ny, nx)
	F = utl.Alloc(ny, nx)
	v := make([]float64, 2)
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			v[0] = o.Xmin[0] + float64(i)*o.Size[0]
			v[1] = o.Xmin[1] + float64(j)*o.Size[1]
			r, e := f(v)
			if e != nil {
				err = e
				return
			}
			X[j][i] = v[0]
			Y[j][i] = v[1]
			F[j][i] = r
		}
	}
	return
}

// Eval3d evaluates function over grid
//  X -- [nz][ny][nx] will hold the grid coordinates
//  Y -- [nz][ny][nx] will hold the grid coordinates
//  Z -- [nz][ny][nx] will hold the grid coordinates
//  F -- [nz][ny][nx] will hold the results
func (o *Grid) Eval3d(f fun.Sv) (X, Y, Z, F [][][]float64, err error) {
	nx := o.Npts[0]
	ny := o.Npts[1]
	nz := o.Npts[2]
	X = utl.Deep3alloc(nz, ny, nx)
	Y = utl.Deep3alloc(nz, ny, nx)
	Z = utl.Deep3alloc(nz, ny, nx)
	F = utl.Deep3alloc(nz, ny, nx)
	v := make([]float64, 3)
	for k := 0; k < nz; k++ {
		for j := 0; j < ny; j++ {
			for i := 0; i < nx; i++ {
				v[0] = o.Xmin[0] + float64(i)*o.Size[0]
				v[1] = o.Xmin[1] + float64(j)*o.Size[1]
				v[2] = o.Xmin[2] + float64(k)*o.Size[2]
				r, e := f(v)
				if e != nil {
					err = e
					return
				}
				X[k][j][i] = v[0]
				Y[k][j][i] = v[1]
				Z[k][j][i] = v[2]
				F[k][j][i] = r
			}
		}
	}
	return
}
