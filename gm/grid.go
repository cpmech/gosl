// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// GridBoundaries generates the IDs of nodes on the boundaries of a rectangular grid
//
//   The sets of boundary nodes are organised in the following order:
//
//              Edges                                         Faces
//
//                2                                       +----------------+
//          +-----------+                               ,'|              ,'|
//          |           |                             ,'  |  ___       ,'  |
//          |           |                           ,'    |,'5,' [0] ,'    |
//         3|           |1                        ,'      |~~~     ,'  ,   |
//          |           |                       +'===============+'  ,'|   |
//          |           |                       |   ,'|   |      |   |3|   |
//          +-----------+                       |   |2|   |      |   |,'   |
//                0                             |   |,'   +- - - | +- - - -+
//                                              |   '   ,'       |       ,'
//                                              |     ,' [1]  ___|     ,'
//                                              |   ,'      ,'4,'|   ,'
//                                              | ,'        ~~~  | ,'
//                                              +----------------+'
func GridBoundaries(npts []int) (edge, face [][]int) {
	ndim := len(npts)
	nx := npts[0]
	ny := npts[1]
	if ndim == 2 {
		edge = make([][]int, 4)   // bottom, right, top, left
		edge[0] = make([]int, nx) // bottom
		edge[1] = make([]int, ny) // right
		edge[2] = make([]int, nx) // top
		edge[3] = make([]int, ny) // left
		for i := 0; i < nx; i++ {
			edge[0][i] = i             // bottom
			edge[2][i] = i + nx*(ny-1) // top
		}
		for j := 0; j < ny; j++ {
			edge[1][j] = j*nx + nx - 1 // right
			edge[3][j] = j * nx        // left
		}
		return
	}
	nz := npts[2]
	face = make([][]int, 6)      // xmin,xmax, ymin,ymax, zmin,zmax
	face[0] = make([]int, ny*nz) // xmin
	face[1] = make([]int, ny*nz) // xmax
	face[2] = make([]int, nx*nz) // ymin
	face[3] = make([]int, nx*nz) // ymax
	face[4] = make([]int, nx*ny) // zmin
	face[5] = make([]int, nx*ny) // zmax
	p := 0
	for j := 0; j < npts[1]; j++ {
		for i := 0; i < npts[0]; i++ {
			face[4][p] = i + nx*j                  // zmin
			face[5][p] = i + nx*j + (nx*ny)*(nz-1) // zmax
			p++
		}
	}
	p = 0
	for k := 0; k < npts[2]; k++ {
		for i := 0; i < npts[0]; i++ {
			face[2][p] = i + (nx*ny)*k             // ymin
			face[3][p] = i + (nx*ny)*k + nx*(ny-1) // ymax
			p++
		}
	}
	p = 0
	for k := 0; k < npts[2]; k++ {
		for j := 0; j < npts[1]; j++ {
			face[0][p] = j*nx + (nx*ny)*k            // xmin
			face[1][p] = j*nx + (nx*ny)*k + (nx - 1) // xmax
			p++
		}
	}
	return
}

// GridBoundaryTag returns a list of nodes marked with given tag
//
//              Edges                                         Faces
//
//               21                                       +----------------+
//          +-----------+                               ,'|              ,'|
//          |           |                             ,'  |  ___       ,'  |
//          |           |                           ,'    |,'31'  10 ,'    |
//        10|           |11                       ,'      |~~~     ,'  ,,  |
//          |           |                       +'===============+'  ,' |  |
//          |           |                       |   ,'|   |      |   |21|  |
//          +-----------+                       |  |20|   |      |   |,'   |
//               20                             |  | ,'   +- - - | +- - - -+
//                                              |   '   ,'       |       ,'
//                                              |     ,' 11   ___|     ,'
//   NOTE: will return empty list if            |   ,'      ,'30'|   ,'
//         tag is not available                 | ,'        ~~~  | ,'
//                                              +----------------+'
func GridBoundaryTag(tag, ndim int, edge, face [][]int) []int {
	if ndim == 2 {
		switch tag {
		case 20:
			return edge[0]
		case 11:
			return edge[1]
		case 21:
			return edge[2]
		case 10:
			return edge[3]
		}
		return nil
	}
	switch tag {
	case 10:
		return face[0]
	case 11:
		return face[1]
	case 20:
		return face[2]
	case 21:
		return face[3]
	case 30:
		return face[4]
	case 31:
		return face[5]
	}
	return nil
}

// grid ////////////////////////////////////////////////////////////////////////////////////////////

// Grid implements a generic (uniform or nonuniform) 2D or 3D grid of points
// See function GridBoundaries() with a picture of how the boundaries are numbered
type Grid struct {
	Ndim int           // space dimension
	Npts []int         // [ndim] number of points along each direction; i.e. ndiv + 1
	N    int           // total number of points
	Edge [][]int       // ids of points on edges: [edge0, edge1, edge2, edge3]
	Face [][]int       // ids of points on faces: [face0, face1, face2, face3, face4, face5]
	X2d  [][]float64   // [ny][nx] 2D grid coordinates
	Y2d  [][]float64   // [ny][nx] 2D grid coordinates
	X3d  [][][]float64 // [nz][ny][nx] 3D grid coordinates
	Y3d  [][][]float64 // [nz][ny][nx] 3D grid coordinates
	Z3d  [][][]float64 // [nz][ny][nx] 3D grid coordinates
	Min  []float64     // [ndim] left/lower-most point
	Max  []float64     // [ndim] right/upper-most point
	Del  []float64     // [ndim] the lengths along each direction (whole box)
}

// NewGrid creates a new Grid structure
//   ndiv -- [ndim] number of divisions for xmax-xmin
func NewGrid(ndiv []int) (o *Grid, err error) {

	// new structure
	o = new(Grid)
	o.Ndim = len(ndiv)

	// number of points along each direction and total number of points
	o.Npts = make([]int, o.Ndim)
	o.N = 1
	for k := 0; k < o.Ndim; k++ {
		o.Npts[k] = ndiv[k] + 1
		o.N *= o.Npts[k]
	}

	// boundaries
	o.Edge, o.Face = GridBoundaries(o.Npts)
	return
}

// GetNodesWithTag returns a list of nodes marked with given tag
func (o *Grid) GetNodesWithTag(tag int) []int {
	return GridBoundaryTag(tag, o.Ndim, o.Edge, o.Face)
}

// SetCoords2d sets 2d coordinates: will allocate X2d[ny][nx] and Y3d[ny][nx]
func (o *Grid) SetCoords2d(X, Y []float64) (err error) {

	if o.Ndim != 2 {
		return chk.Err("grid must be 2D\n")
	}
	if len(X) != o.Npts[0] {
		return chk.Err("number of points along x is incorrect. %d != %d\n", len(X), o.Npts[0])
	}
	if len(Y) != o.Npts[1] {
		return chk.Err("number of points along y is incorrect. %d != %d\n", len(Y), o.Npts[1])
	}
	nx := o.Npts[0]
	ny := o.Npts[1]
	o.X2d = utl.Alloc(ny, nx)
	o.Y2d = utl.Alloc(ny, nx)
	o.Min = []float64{X[0], Y[0]}
	o.Max = []float64{X[0], Y[0]}
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			o.X2d[j][i] = X[i]
			o.Y2d[j][i] = Y[j]
		}
		o.Min[1] = utl.Min(o.Min[1], Y[j])
		o.Max[1] = utl.Max(o.Max[1], Y[j])
	}
	for i := 0; i < nx; i++ {
		o.Min[0] = utl.Min(o.Min[0], X[i])
		o.Max[0] = utl.Max(o.Max[0], X[i])
	}
	o.Del = make([]float64, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		o.Del[k] = o.Max[k] - o.Min[k]
	}
	return
}

// SetCoords3d sets 3d coordinates: will allocate X3d[nz][ny][nx], Y3d[nz][ny][nx], and Z3d[nz][ny][nx]
func (o *Grid) SetCoords3d(X, Y, Z []float64) (err error) {
	if o.Ndim != 3 {
		return chk.Err("grid must be 3D\n")
	}
	if len(X) != o.Npts[0] {
		return chk.Err("number of points along x is incorrect. %d != %d\n", len(X), o.Npts[0])
	}
	if len(Y) != o.Npts[1] {
		return chk.Err("number of points along y is incorrect. %d != %d\n", len(Y), o.Npts[1])
	}
	if len(Z) != o.Npts[2] {
		return chk.Err("number of points along z is incorrect. %d != %d\n", len(Z), o.Npts[2])
	}
	nx := o.Npts[0]
	ny := o.Npts[1]
	nz := o.Npts[2]
	o.X3d = utl.Deep3alloc(nz, ny, nx)
	o.Y3d = utl.Deep3alloc(nz, ny, nx)
	o.Z3d = utl.Deep3alloc(nz, ny, nx)
	o.Min = []float64{X[0], Y[0], Z[0]}
	o.Max = []float64{X[0], Y[0], Z[0]}
	for k := 0; k < nz; k++ {
		for j := 0; j < ny; j++ {
			for i := 0; i < nx; i++ {
				o.X3d[k][j][i] = X[i]
				o.Y3d[k][j][i] = Y[j]
				o.Z3d[k][j][i] = Z[k]
			}
		}
		o.Min[2] = utl.Min(o.Min[2], Z[k])
		o.Max[2] = utl.Max(o.Max[2], Z[k])
	}
	for j := 0; j < ny; j++ {
		o.Min[1] = utl.Min(o.Min[1], Y[j])
		o.Max[1] = utl.Max(o.Max[1], Y[j])
	}
	for i := 0; i < nx; i++ {
		o.Min[0] = utl.Min(o.Min[0], X[i])
		o.Max[0] = utl.Max(o.Max[0], X[i])
	}
	o.Del = make([]float64, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		o.Del[k] = o.Max[k] - o.Min[k]
	}
	return
}

// Draw draws grid
func (o *Grid) Draw(withTxt bool, argsGrid, argsTxt *plt.A) {

	// configuration
	if argsGrid == nil {
		argsGrid = &plt.A{C: "#427ce5", Lw: 0.8, NoClip: true}
	}

	// draw grid
	if o.Ndim == 2 {
		plt.Grid2d(o.X2d, o.Y2d, false, argsGrid, nil)
	} else {
		Zlevels := make([]float64, o.Npts[2])
		for k := 0; k < o.Npts[2]; k++ {
			Zlevels[k] = o.Z3d[k][0][0]
		}
		plt.Grid3d(o.X3d[0], o.Y3d[0], Zlevels, argsGrid)
	}

	// grid txt
	if withTxt {

		// configuration
		if argsTxt == nil {
			argsTxt = &plt.A{C: "orange", Fsz: 8}
		}

		// add text
		if o.Ndim == 2 {
			for j := 0; j < o.Npts[1]; j++ {
				for i := 0; i < o.Npts[0]; i++ {
					idx := i + j*o.Npts[0]
					txt := io.Sf("%d", idx)
					plt.Text(o.X2d[j][i], o.Y2d[j][i], txt, argsTxt)
				}
			}
		} else {
			for k := 0; k < o.Npts[2]; k++ {
				for j := 0; j < o.Npts[1]; j++ {
					for i := 0; i < o.Npts[0]; i++ {
						idx := i + j*o.Npts[0] + k*o.Npts[0]*o.Npts[1]
						txt := io.Sf("%d", idx)
						plt.Text3d(o.X3d[k][j][i], o.Y3d[k][j][i], o.Z3d[k][j][i], txt, argsTxt)
					}
				}
			}
		}
	}
}

// uniform grid ///////////////////////////////////////////////////////////////////////////////////

// UniformGrid implements a uniform 2D or 3D grid of points (based on Bins)
// See function GridBoundaries() with a picture of how the boundaries are numbered
type UniformGrid struct {
	Bins // derived

	N    int     // total number of points
	Edge [][]int // ids of points on edges: [edge0, edge1, edge2, edge3]
	Face [][]int // ids of points on faces: [face0, face1, face2, face3, face4, face5]
}

// NewUniformGrid creates a new UniformGrid structure
//   xmin -- [ndim] min/initial coordinates of the whole space (box/cube)
//   xmax -- [ndim] max/final coordinates of the whole space (box/cube)
//   ndiv -- [ndim] number of divisions for xmax-xmin
func NewUniformGrid(xmin, xmax []float64, ndiv []int) (o *UniformGrid, err error) {
	o = new(UniformGrid)
	err = o.Bins.Init(xmin, xmax, ndiv)
	o.N = 1
	for k := 0; k < o.Ndim; k++ {
		o.N *= o.Npts[k]
	}
	o.Edge, o.Face = GridBoundaries(o.Npts)
	return
}

// GetNodesWithTag returns a list of nodes marked with given tag
func (o *UniformGrid) GetNodesWithTag(tag int) []int {
	return GridBoundaryTag(tag, o.Ndim, o.Edge, o.Face)
}

// Eval2d evaluates function over grid
//  X -- [ny][nx] will hold the grid coordinates
//  Y -- [ny][nx] will hold the grid coordinates
//  F -- [ny][nx] will hold the results
func (o *UniformGrid) Eval2d(f fun.Sv) (X, Y, F [][]float64, err error) {
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
func (o *UniformGrid) Eval3d(f fun.Sv) (X, Y, Z, F [][][]float64, err error) {
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
