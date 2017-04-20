// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import "github.com/cpmech/gosl/io"

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
	updateBufferAndClose(&bufferPy, args, false)
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
	updateBufferAndClose(&bufferPy, args, false)
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
	updateBufferAndClose(&bufferPy, args, false)
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
	updateBufferAndClose(&bufferPy, args, false)
}

// Camera sets camera in 3d graph. Sets the elevation and azimuth of the axes.
//   elev -- is the elevation angle in the z plane
//   azim -- is the azimuth angle in the x,y plane
func Camera(elev, azim float64, args *A) {
	io.Ff(&bufferPy, "plt.gca().view_init(elev=%g, azim=%g", elev, azim)
	updateBufferAndClose(&bufferPy, args, false)
}

// AxDist sets distance in 3d graph. e.g. to avoid clipping due to tight bbox
func AxDist(dist float64) {
	io.Ff(&bufferPy, "plt.gca().dist = %g\n", dist)
}

// Text3d adds text to 3d plot
func Text3d(x, y, z float64, txt string, doInit bool, args *A) {
	n := get3daxes(doInit)
	io.Ff(&bufferPy, "t%d = ax%d.text(%g,%g,%g,r'%s'", n, n, x, y, z, txt)
	updateBufferAndClose(&bufferPy, args, false)
}

// Triad draws icon indicating x-y-z origin and direction
func Triad(length float64, doInit bool, argsLines, argsText *A) {
	a := argsLines
	if a == nil {
		a = &A{C: "black", Lw: 1.2}
	}
	Plot3dLine([]float64{0, length}, []float64{0, 0}, []float64{0, 0}, doInit, a)
	Plot3dLine([]float64{0, 0}, []float64{0, length}, []float64{0, 0}, false, a)
	Plot3dLine([]float64{0, 0}, []float64{0, 0}, []float64{0, length}, false, a)
	b := argsText
	if b == nil {
		b = &A{C: "black", Fsz: 10, Ha: "center", Va: "center"}
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
