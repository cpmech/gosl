// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// Plot3dLine plots 3d line
func Plot3dLine(X, Y, Z []float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("X%d", n)
	sy := io.Sf("Y%d", n)
	sz := io.Sf("Z%d", n)
	genArray(&bufferPy, sx, X)
	genArray(&bufferPy, sy, Y)
	genArray(&bufferPy, sz, Z)
	io.Ff(&bufferPy, "p%d = ax%d.plot(%s,%s,%s", n, n, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Plot3dPoint plot 3d point
func Plot3dPoint(x, y, z float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	io.Ff(&bufferPy, "p%d = ax%d.scatter(%g,%g,%g", n, n, x, y, z)
	updateBufferAndClose(&bufferPy, args, false, true)
}

// Plot3dPoints plots 3d points
func Plot3dPoints(X, Y, Z []float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("X%d", n)
	sy := io.Sf("Y%d", n)
	sz := io.Sf("Z%d", n)
	genArray(&bufferPy, sx, X)
	genArray(&bufferPy, sy, Y)
	genArray(&bufferPy, sz, Z)
	io.Ff(&bufferPy, "p%d = ax%d.scatter(%s,%s,%s", n, n, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false, true)
}

// Wireframe draws wireframe
func Wireframe(X, Y, Z [][]float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("X%d", n)
	sy := io.Sf("Y%d", n)
	sz := io.Sf("Z%d", n)
	genMat(&bufferPy, sx, X)
	genMat(&bufferPy, sy, Y)
	genMat(&bufferPy, sz, Z)
	_, rs, cs := args3d(args)
	io.Ff(&bufferPy, "p%d = ax%d.plot_wireframe(%s,%s,%s,rstride=%d,cstride=%d", n, n, sx, sy, sz, rs, cs)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Surface draws surface
func Surface(X, Y, Z [][]float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("X%d", n)
	sy := io.Sf("Y%d", n)
	sz := io.Sf("Z%d", n)
	genMat(&bufferPy, sx, X)
	genMat(&bufferPy, sy, Y)
	genMat(&bufferPy, sz, Z)
	cmapIdx, rs, cs := args3d(args)
	io.Ff(&bufferPy, "p%d = ax%d.plot_surface(%s,%s,%s,cmap=getCmap(%d),rstride=%d,cstride=%d", n, n, sx, sy, sz, cmapIdx, rs, cs)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Camera sets camera in 3d graph. Sets the elevation and azimuth of the axes.
//   elev -- is the elevation angle in the z plane
//   azim -- is the azimuth angle in the x,y plane
func Camera(elev, azim float64, args *A) {
	io.Ff(&bufferPy, "plt.gca().view_init(elev=%g, azim=%g", elev, azim)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// AxDist sets distance in 3d graph. e.g. to avoid clipping due to tight bbox
func AxDist(dist float64) {
	io.Ff(&bufferPy, "plt.gca().dist = %g\n", dist)
}

// Text3d adds text to 3d plot
func Text3d(x, y, z float64, txt string, doInit bool, args *A) {
	n := get3daxes(doInit)
	io.Ff(&bufferPy, "t%d = ax%d.text(%g,%g,%g,r'%s'", n, n, x, y, z, txt)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Triad draws icon indicating x-y-z origin and direction
func Triad(length float64, doInit, labels bool, argsLines, argsText *A) {
	a := argsLines
	if a == nil {
		a = &A{C: "black", Lw: 1.2}
	}
	Plot3dLine([]float64{0, length}, []float64{0, 0}, []float64{0, 0}, doInit, a)
	Plot3dLine([]float64{0, 0}, []float64{0, length}, []float64{0, 0}, false, a)
	Plot3dLine([]float64{0, 0}, []float64{0, 0}, []float64{0, length}, false, a)
	if labels {
		b := argsText
		if b == nil {
			b = &A{C: "black", Fsz: 10, Ha: "center", Va: "center"}
		}
		g := 0.05 * length
		Text3d(length+g, 0, 0, "x", false, b)
		Text3d(0, length+g, 0, "y", false, b)
		Text3d(0, 0, length+g, "z", false, b)
	}
}

// Scale3d scales 3d axes
func Scale3d(xmin, xmax, ymin, ymax, zmin, zmax float64, equal bool) {
	dx := xmax - xmin
	dy := ymax - ymin
	dz := zmax - zmin
	xmid := (xmin + xmax) / 2.0
	ymid := (ymin + ymax) / 2.0
	zmid := (zmin + zmax) / 2.0
	if equal {
		maxRange := dx
		if dy > maxRange {
			maxRange = dy
		}
		if dz > maxRange {
			maxRange = dz
		}
		dx = maxRange
		dy = maxRange
		dz = maxRange
	}
	dx *= 0.5
	dy *= 0.5
	dz *= 0.5
	xleft, xright := xmid-dx, xmid+dx
	yleft, yright := ymid-dy, ymid+dy
	zleft, zright := zmid-dz, zmid+dz
	io.Ff(&bufferPy, "plt.gca().set_xlim(%g, %g)\n", xleft, xright)
	io.Ff(&bufferPy, "plt.gca().set_ylim(%g, %g)\n", yleft, yright)
	io.Ff(&bufferPy, "plt.gca().set_zlim(%g, %g)\n", zleft, zright)
}

// Default3dView sets default 3d view (camera and scale)
func Default3dView(xmin, xmax, ymin, ymax, zmin, zmax float64, equal bool) {
	elev, azim := 30.0, 20.0
	Camera(elev, azim, nil)
	AxDist(10.5)
	Scale3d(xmin, xmax, ymin, ymax, zmin, zmax, equal)
}

// Draw3dVector adds segment to figure
//   p -- starting point
//   v -- vector
//   sf -- scale factor
//   normed -- normalised
func Draw3dVector(p, v []float64, sf float64, normed, doInit bool, args *A) {
	if len(p) != 3 || len(v) != 3 {
		return
	}
	scale := sf
	if normed {
		norm := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
		if norm > 1e-10 {
			scale = sf / norm
		}
	}
	n := get3daxes(doInit)
	io.Ff(&bufferPy, "p%d = ax%d.plot([%g,%g],[%g,%g],[%g,%g]", n, n,
		p[0], p[0]+v[0]*scale,
		p[1], p[1]+v[1]*scale,
		p[2], p[2]+v[2]*scale)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// PlaneZ draws a plane that has a normal vector with non-zero z component.
// The plane may be perpendicular to z.
//   p -- point on plane
//   n -- normal vector
//   nu -- number of divisions along one direction on plane
//   nv -- number of divisions along the orther direction on plane
//   showPN -- show point and normal
func PlaneZ(p, n []float64, xmin, xmax, ymin, ymax float64, nu, nv int, showPN, doInit bool, args *A) {
	if len(p) != 3 || len(n) != 3 {
		return
	}
	if math.Abs(n[2]) < 1e-10 {
		return
	}
	d := -n[0]*p[0] - n[1]*p[1] - n[2]*p[2]
	X, Y, Z := utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nu, nv, func(x, y float64) float64 {
		return (-d - n[0]*x - n[1]*y) / n[2]
	})
	Wireframe(X, Y, Z, doInit, args)
	if showPN {
		a := &A{C: "k", Ec: "k", M: "."}
		Plot3dPoint(p[0], p[1], p[2], false, a)
		a.M, a.Ec = "", ""
		Draw3dVector(p, n, 1.0, true, false, a)
	}
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

func get3daxes(doInit bool) (n int) {
	n = bufferPy.Len()
	if doInit {
		io.Ff(&bufferPy, "ax%d = plt.gcf().add_subplot(111, projection='3d')\n", n)
		io.Ff(&bufferPy, "ax%d.set_xlabel('x');ax%d.set_ylabel('y');ax%d.set_zlabel('z')\n", n, n, n)
		io.Ff(&bufferPy, "addToEA(ax%d)\n", n)
	} else {
		io.Ff(&bufferPy, "ax%d = plt.gca()\n", n)
	}
	return
}
