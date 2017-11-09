// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// RectGrid implements a rectangular grid (uniform or nonuniform) of points in 2D or 3D
//
//  NOTE: use methods GenUniform, Set2d, or Set3d to generate coordinates
//
type RectGrid struct {
	ndim   int           // space dimension
	npts   []int         // [ndim] number of points along each direction; i.e. ndiv + 1
	size   int           // total number of points
	coords [][]float64   // [ndim][nx] coordinates along each direction
	min    []float64     // [ndim] left/lower-most point
	max    []float64     // [ndim] right/upper-most point
	length []float64     // [ndim] the lengths along each direction (whole box)
	x2d    [][]float64   // [ny][nx] 2D grid coordinates
	y2d    [][]float64   // [ny][nx] 2D grid coordinates
	x3d    [][][]float64 // [nz][ny][nx] 3D grid coordinates
	y3d    [][][]float64 // [nz][ny][nx] 3D grid coordinates
	z3d    [][][]float64 // [nz][ny][nx] 3D grid coordinates
	edge   [][]int       // ids of points on edges: [edge0, edge1, edge2, edge3]
	face   [][]int       // ids of points on faces: [face0, face1, face2, face3, face4, face5]
}

// GenUniform generates uniform grid
//  min  -- min x-y-z values, len==ndim: [xmin, ymin, zmin]
//  max  -- max x-y-z values, len==ndim: [xmax, ymax, zmax]
//  ndiv -- number of divisions along each direction len==ndim: [ndivx, ndivy, ndivz]
//  genMeshCoords -- generate "meshgrid" coordinates. Accessed with Mesh2d or Mesh3d
func (o *RectGrid) GenUniform(min, max []float64, ndiv []int, genMeshCoords bool) {
	o.ndim = len(min)
	if o.ndim != 2 && o.ndim != 3 {
		chk.Panic("ndim must be 2 or 3. ndim=%d is invalid\n", o.ndim)
	}
	if len(max) != o.ndim {
		chk.Panic("len(max) must be equal to len(min) == ndim. %d != %d\n", len(max), o.ndim)
	}
	if len(ndiv) != o.ndim {
		chk.Panic("len(ndiv) must be equal to len(min) == ndim. %d != %d\n", len(ndiv), o.ndim)
	}
	o.npts = make([]int, o.ndim)
	o.min = utl.GetCopy(min)
	o.max = utl.GetCopy(max)
	o.length = make([]float64, o.ndim)
	o.coords = make([][]float64, o.ndim)
	o.size = 1
	for i := 0; i < o.ndim; i++ {
		o.npts[i] = ndiv[i] + 1
		o.size *= o.npts[i]
		o.coords[i] = utl.LinSpace(o.min[i], o.max[i], o.npts[i])
		o.length[i] = o.max[i] - o.min[i]
	}
	if genMeshCoords {
		o.genMesh()
	}
	o.edge, o.face = o.boundaries(o.npts)
}

// Set2d sets coordinates along each direction
//  genMeshCoords -- generate "meshgrid" coordinates. Accessed with Mesh2d or Mesh3d
func (o *RectGrid) Set2d(X, Y []float64, genMeshCoords bool) {
	nx := len(X)
	ny := len(Y)
	o.ndim = 2
	o.npts = []int{nx, ny}
	o.size = nx * ny
	o.coords = [][]float64{utl.GetCopy(X), utl.GetCopy(Y)}
	o.limits()
	if genMeshCoords {
		o.genMesh()
	}
	o.edge, o.face = o.boundaries(o.npts)
}

// Set3d sets coordinates along each direction
//  genMeshCoords -- generate "meshgrid" coordinates. Accessed with Mesh2d or Mesh3d
func (o *RectGrid) Set3d(X, Y, Z []float64, genMeshCoords bool) {
	nx := len(X)
	ny := len(Y)
	nz := len(Z)
	o.ndim = 3
	o.npts = []int{nx, ny, nz}
	o.size = nx * ny * nz
	o.coords = [][]float64{utl.GetCopy(X), utl.GetCopy(Y), utl.GetCopy(Z)}
	o.limits()
	if genMeshCoords {
		o.genMesh()
	}
	o.edge, o.face = o.boundaries(o.npts)
}

// Ndim returns the space dimension
func (o *RectGrid) Ndim() int {
	return o.ndim
}

// Npts returns the number of points along each direction
//   dim -- direction in [0, ndim-1]
func (o *RectGrid) Npts(dim int) int {
	return o.npts[dim]
}

// Size returns the total number of points; e.g. Npts(0) * Npts(1) * Npts(2)
func (o *RectGrid) Size() int {
	return o.size
}

// Min returns mininum x-y-z value
//   dim -- direction in [0, ndim-1]
func (o *RectGrid) Min(dim int) float64 {
	return o.min[dim]
}

// Max returns maximum x-y-z value
//   dim -- direction in [0, ndim-1]
func (o *RectGrid) Max(dim int) float64 {
	return o.max[dim]
}

// Length returns the lengths along each direction (whole box)
//   dim -- direction in [0, ndim-1]
func (o *RectGrid) Length(dim int) float64 {
	return o.length[dim]
}

// Mesh2d returns 2D "meshgrid"
func (o *RectGrid) Mesh2d() (X, Y [][]float64) {
	if !utl.Deep2checkSize(o.npts[0], o.npts[1], o.x2d) {
		o.genMesh()
	}
	return o.x2d, o.y2d
}

// Mesh3d returns 3D "meshgrid"
func (o *RectGrid) Mesh3d() (X, Y, Z [][][]float64) {
	if !utl.Deep3checkSize(o.npts[0], o.npts[1], o.npts[2], o.x3d) {
		o.genMesh()
	}
	return o.x3d, o.y3d, o.z3d
}

// Node returns in 'x' the coordinates of a node 'n'
func (o *RectGrid) Node(x la.Vector, n int) {
	nx := o.npts[0]
	if o.ndim == 2 {
		i := n % nx
		j := n / nx
		x[0] = o.coords[0][i]
		x[1] = o.coords[1][j]
		return
	}
	ny := o.npts[1]
	k := n / (nx * ny)
	m := n % (nx * ny)
	i := m % nx
	j := m / nx
	x[0] = o.coords[0][i]
	x[1] = o.coords[1][j]
	x[2] = o.coords[2][k]
}

// GetNode returns coordinates of node 'n'
func (o *RectGrid) GetNode(n int) (x la.Vector) {
	x = la.NewVector(o.ndim)
	o.Node(x, n)
	return
}

// Edge returns the ids of points on edges: [edge0, edge1, edge2, edge3]
//            2
//      +-----------+
//      |           |
//      |           |
//     3|           |1
//      |           |
//      |           |
//      +-----------+
//            0
func (o *RectGrid) Edge(iEdge int) []int {
	return o.edge[iEdge]
}

// EdgeT returns a list of nodes marked with given tag
//           21
//      +-----------+
//      |           |
//      |           |
//    10|           |11
//      |           |
//      |           |
//      +-----------+
//           20
//
//   NOTE: will return empty list if tag is not available
//
func (o *RectGrid) EdgeT(tag int) []int {
	switch tag {
	case 20:
		return o.edge[0]
	case 11:
		return o.edge[1]
	case 21:
		return o.edge[2]
	case 10:
		return o.edge[3]
	}
	return nil
}

// Face returns the ids of points on faces: [face0, face1, face2, face3, face4, face5]
//               +----------------+
//             ,'|              ,'|
//           ,'  |  ___       ,'  |
//         ,'    |,'5,' [0] ,'    |
//       ,'      |~~~     ,'  ,   |
//     +'===============+'  ,'|   |
//     |   ,'|   |      |   |3|   |
//     |   |2|   |      |   |,'   |
//     |   |,'   +- - - | +- - - -+
//     |   '   ,'       |       ,'
//     |     ,' [1]  ___|     ,'
//     |   ,'      ,'4,'|   ,'
//     | ,'        ~~~  | ,'
//     +----------------+'
func (o *RectGrid) Face(iFace int) []int {
	return o.face[iFace]
}

// FaceT returns a list of nodes marked with given tag
//               +----------------+
//             ,'|              ,'|
//           ,'  |  ___       ,'  |
//         ,'    |,'31'  10 ,'    |
//       ,'      |~~~     ,'  ,,  |
//     +'===============+'  ,' |  |
//     |   ,'|   |      |   |21|  |
//     |  |20|   |      |   |,'   |
//     |  | ,'   +- - - | +- - - -+
//     |   '   ,'       |       ,'
//     |     ,' 11   ___|     ,'
//     |   ,'      ,'30'|   ,'
//     | ,'        ~~~  | ,'
//     +----------------+'
//
//   NOTE: will return empty list if tag is not available
//
func (o *RectGrid) FaceT(tag int) []int {
	switch tag {
	case 100:
		return o.face[0]
	case 101:
		return o.face[1]
	case 200:
		return o.face[2]
	case 201:
		return o.face[3]
	case 300:
		return o.face[4]
	case 301:
		return o.face[5]
	}
	return nil
}

// Boundary returns list of edge or face nodes on boundary
//   NOTE: will return empty list if tag is not available
func (o *RectGrid) Boundary(tag int) []int {
	if tag > 50 {
		if o.ndim == 2 {
			return nil
		}
		return o.FaceT(tag)
	}
	return o.EdgeT(tag)
}

// Draw draws grid
//  coordsGenerate -- tells whether the meshgrid coordinates are generated or not.
//                    if the size of the current arrays are the same, no new coordinates
//                    will be generated. Thus, make sure to call the "set" commands with the
//                    genMeshCoordinates flag equal to true if when updating the grid.
func (o *RectGrid) Draw(withTxt bool, argsGrid, argsTxt *plt.A) (coordsGenerated bool) {

	// check dimension and number of points along each direction
	if o.ndim != 2 && o.ndim != 3 {
		return
	}
	nx := o.npts[0]
	ny := o.npts[1]
	if nx < 1 || ny < 1 {
		return
	}
	nz := 1
	if o.ndim == 3 {
		nz = o.npts[2]
		if nz < 1 {
			return
		}
	}

	// generate coords
	if o.ndim == 2 {
		if !utl.Deep2checkSize(nx, ny, o.x2d) {
			coordsGenerated = true
			o.genMesh()
		}
	} else {
		if !utl.Deep3checkSize(nx, ny, nz, o.x3d) {
			coordsGenerated = true
			o.genMesh()
		}
	}

	// configuration
	if argsGrid == nil {
		argsGrid = &plt.A{C: "#427ce5", Lw: 0.8, NoClip: true}
	}

	// draw grid
	if o.ndim == 2 {
		plt.Grid2d(o.x2d, o.y2d, false, argsGrid, nil)
	} else {
		plt.Grid3dZlevels(o.x3d[0], o.y3d[0], o.coords[2], argsGrid)
	}

	// grid txt
	if withTxt {

		// configuration
		if argsTxt == nil {
			argsTxt = &plt.A{C: plt.C(2, 0), Fsz: 7}
		}

		// add text
		if o.ndim == 2 {
			for j := 0; j < o.npts[1]; j++ {
				for i := 0; i < o.npts[0]; i++ {
					idx := i + j*o.npts[0]
					txt := io.Sf("%d", idx)
					plt.Text(o.x2d[j][i], o.y2d[j][i], txt, argsTxt)
				}
			}
		} else {
			for k := 0; k < o.npts[2]; k++ {
				for j := 0; j < o.npts[1]; j++ {
					for i := 0; i < o.npts[0]; i++ {
						idx := i + j*o.npts[0] + k*o.npts[0]*o.npts[1]
						txt := io.Sf("%d", idx)
						plt.Text3d(o.x3d[k][j][i], o.y3d[k][j][i], o.z3d[k][j][i], txt, argsTxt)
					}
				}
			}
		}
	}
	return
}

// limits computes min, max and length
func (o *RectGrid) limits() {
	o.min = make([]float64, o.ndim)
	o.max = make([]float64, o.ndim)
	o.length = make([]float64, o.ndim)
	for i := 0; i < o.ndim; i++ {
		if len(o.coords[i]) < 1 {
			continue
		}
		o.min[i] = o.coords[i][0]
		o.max[i] = o.coords[i][0]
		for j := 0; j < len(o.coords[i]); j++ {
			o.min[i] = utl.Min(o.min[i], o.coords[i][j])
			o.max[i] = utl.Max(o.max[i], o.coords[i][j])
		}
		o.length[i] = o.max[i] - o.min[i]
	}
}

// genMesh generates meshgrid coordinates
func (o *RectGrid) genMesh() {
	if o.ndim == 2 {
		o.x2d = utl.Alloc(o.npts[1], o.npts[0])
		o.y2d = utl.Alloc(o.npts[1], o.npts[0])
		for j := 0; j < o.npts[1]; j++ {
			for i := 0; i < o.npts[0]; i++ {
				o.x2d[j][i] = o.coords[0][i]
				o.y2d[j][i] = o.coords[1][j]
			}
		}
		return
	}
	o.x3d = utl.Deep3alloc(o.npts[2], o.npts[1], o.npts[0])
	o.y3d = utl.Deep3alloc(o.npts[2], o.npts[1], o.npts[0])
	o.z3d = utl.Deep3alloc(o.npts[2], o.npts[1], o.npts[0])
	for k := 0; k < o.npts[2]; k++ {
		for j := 0; j < o.npts[1]; j++ {
			for i := 0; i < o.npts[0]; i++ {
				o.x3d[k][j][i] = o.coords[0][i]
				o.y3d[k][j][i] = o.coords[1][j]
				o.z3d[k][j][i] = o.coords[2][k]
			}
		}
	}
}

// boundaries generates the IDs of nodes on the boundaries of a rectangular grid
func (o *RectGrid) boundaries(npts []int) (edge, face [][]int) {
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
