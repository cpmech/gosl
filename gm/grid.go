// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
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
			edge[2][i] = i + (ny-1)*nx // top
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
	chk.Panic("TODO: implement ids of faces\n")
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

// UniformGrid implements a 2D or 3D grid of points (based on Bins)
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

// GridBoundaryTag returns a list of nodes marked with given tag
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
