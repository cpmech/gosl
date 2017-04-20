// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// AxisRange3d sets x, y, and z ranges (i.e. limits)
func AxisRange3d(xmin, xmax, ymin, ymax, zmin, zmax float64) {
	io.Ff(&bufferPy, "plt.gca().set_xlim3d(%g,%g)\ngca().set_ylim3d(%g,%g)\ngca().set_zlim3d(%g,%g)\n", xmin, xmax, ymin, ymax, zmin, zmax)
}

// Plot3dLine plots 3d line
func Plot3dLine(X, Y, Z []float64, args *A) {
	createAxes3d()
	uid := genUid()
	sx := io.Sf("X%d", uid)
	sy := io.Sf("Y%d", uid)
	sz := io.Sf("Z%d", uid)
	genArray(&bufferPy, sx, X)
	genArray(&bufferPy, sy, Y)
	genArray(&bufferPy, sz, Z)
	io.Ff(&bufferPy, "p%d = AX3D.plot(%s,%s,%s", uid, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Plot3dPoint plot 3d point
func Plot3dPoint(x, y, z float64, args *A) {
	createAxes3d()
	io.Ff(&bufferPy, "p%d = AX3D.scatter(%g,%g,%g", genUid(), x, y, z)
	updateBufferAndClose(&bufferPy, args, false, true)
}

// Plot3dPoints plots 3d points
func Plot3dPoints(X, Y, Z []float64, args *A) {
	createAxes3d()
	uid := genUid()
	sx := io.Sf("X%d", uid)
	sy := io.Sf("Y%d", uid)
	sz := io.Sf("Z%d", uid)
	genArray(&bufferPy, sx, X)
	genArray(&bufferPy, sy, Y)
	genArray(&bufferPy, sz, Z)
	io.Ff(&bufferPy, "p%d = AX3D.scatter(%s,%s,%s", uid, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false, true)
}

// Wireframe draws wireframe
func Wireframe(X, Y, Z [][]float64, args *A) {
	createAxes3d()
	uid := genUid()
	sx := io.Sf("X%d", uid)
	sy := io.Sf("Y%d", uid)
	sz := io.Sf("Z%d", uid)
	genMat(&bufferPy, sx, X)
	genMat(&bufferPy, sy, Y)
	genMat(&bufferPy, sz, Z)
	_, rs, cs := args3d(args)
	io.Ff(&bufferPy, "p%d = AX3D.plot_wireframe(%s,%s,%s,rstride=%d,cstride=%d", uid, sx, sy, sz, rs, cs)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Surface draws surface
func Surface(X, Y, Z [][]float64, args *A) {
	createAxes3d()
	uid := genUid()
	sx := io.Sf("X%d", uid)
	sy := io.Sf("Y%d", uid)
	sz := io.Sf("Z%d", uid)
	genMat(&bufferPy, sx, X)
	genMat(&bufferPy, sy, Y)
	genMat(&bufferPy, sz, Z)
	cmapIdx, rs, cs := args3d(args)
	io.Ff(&bufferPy, "p%d = AX3D.plot_surface(%s,%s,%s,cmap=getCmap(%d),rstride=%d,cstride=%d", uid, sx, sy, sz, cmapIdx, rs, cs)
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
func Text3d(x, y, z float64, txt string, args *A) {
	createAxes3d()
	io.Ff(&bufferPy, "t%d = AX3D.text(%g,%g,%g,r'%s'", genUid(), x, y, z, txt)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Triad draws icon indicating x-y-z origin and direction
func Triad(length float64, labels bool, argsLines, argsText *A) {
	a := argsLines
	if a == nil {
		a = &A{C: "black", Lw: 1.2}
	}
	Plot3dLine([]float64{0, length}, []float64{0, 0}, []float64{0, 0}, a)
	Plot3dLine([]float64{0, 0}, []float64{0, length}, []float64{0, 0}, a)
	Plot3dLine([]float64{0, 0}, []float64{0, 0}, []float64{0, length}, a)
	if labels {
		b := argsText
		if b == nil {
			b = &A{C: "black", Fsz: 10, Ha: "center", Va: "center"}
		}
		g := 0.05 * length
		Text3d(length+g, 0, 0, "x", b)
		Text3d(0, length+g, 0, "y", b)
		Text3d(0, 0, length+g, "z", b)
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
func Draw3dVector(p, v []float64, sf float64, normed bool, args *A) {
	scale := sf
	if normed {
		norm := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
		if norm > 1e-10 {
			scale = sf / norm
		}
	}
	createAxes3d()
	io.Ff(&bufferPy, "p%d = AX3D.plot([%g,%g],[%g,%g],[%g,%g]", genUid(),
		p[0], p[0]+v[0]*scale,
		p[1], p[1]+v[1]*scale,
		p[2], p[2]+v[2]*scale)
	updateBufferAndClose(&bufferPy, args, false, false)
}

// Diag3d draws diagonal of 3d space
func Diag3d(scale float64, args *A) {
	createAxes3d()
	a := args
	if a == nil {
		a = &A{C: "k"}
	}
	io.Ff(&bufferPy, "p%d = AX3D.plot([0,%g],[0,%g],[0,%g]", genUid(), scale, scale, scale)
	updateBufferAndClose(&bufferPy, a, false, false)
}

// 3d shapes using meshgrid ///////////////////////////////////////////////////////////////////////

// addSurfAndOrWire adds surface and or wireframe
func addSurfAndOrWire(X, Y, Z [][]float64, args *A) {
	if args == nil {
		Wireframe(X, Y, Z, nil)
		return
	}
	if !args.Surf && !args.Wire {
		Wireframe(X, Y, Z, args)
	}
	if args.Surf {
		Surface(X, Y, Z, args)
	}
	if args.Wire {
		Wireframe(X, Y, Z, args)
	}
}

// PlaneZ draws a plane that has a normal vector with non-zero z component.
// The plane may be perpendicular to z.
//  Input:
//     p -- point on plane
//     n -- normal vector
//     nu -- number of divisions along one direction on plane
//     nv -- number of divisions along the orther direction on plane
//     showPN -- show point and normal
//   Output:
//     X, Y, Z -- the coordinages of all points as in a meshgrid
func PlaneZ(p, n []float64, xmin, xmax, ymin, ymax float64, nu, nv int, showPN bool, args *A) (X, Y, Z [][]float64) {
	if math.Abs(n[2]) < 1e-10 {
		return
	}
	d := -n[0]*p[0] - n[1]*p[1] - n[2]*p[2]
	X, Y, Z = utl.MeshGrid2dF(xmin, xmax, ymin, ymax, nu, nv, func(x, y float64) float64 {
		return (-d - n[0]*x - n[1]*y) / n[2]
	})
	addSurfAndOrWire(X, Y, Z, args)
	if showPN {
		a := &A{C: "k", Ec: "k", M: "."}
		Plot3dPoint(p[0], p[1], p[2], a)
		a.M, a.Ec = "", ""
		Draw3dVector(p, n, 1.0, true, a)
	}
	return
}

// Hemisphere draws Hemisphere
//   Input:
//     c -- centre coordinates. may be nil
//     r -- radius
//     alphaMin -- minimum circumference angle (degrees)
//     alphaMax -- minimum circumference angle (degrees)
//     nu -- number of divisions along one direction on plane
//     nv -- number of divisions along the orther direction on plane
//     cup -- upside-down; like a cup
//     surface -- generate surface
//     wireframe -- generate wireframe
//   Output:
//     X, Y, Z -- the coordinages of all points as in a meshgrid
func Hemisphere(c []float64, r, alphaMin, alphaMax float64, nu, nv int, cup bool, args *A) (X, Y, Z [][]float64) {
	if c == nil {
		c = []float64{0, 0, 0}
	}
	amin := alphaMin * math.Pi / 180.0
	amax := alphaMax * math.Pi / 180.0
	da := (amax - amin) / float64(nu)
	db := (math.Pi / 2.0) / float64(nv)
	X = make([][]float64, nu+1)
	Y = make([][]float64, nu+1)
	Z = make([][]float64, nu+1)
	for i := 0; i < nu+1; i++ {
		X[i] = make([]float64, nv+1)
		Y[i] = make([]float64, nv+1)
		Z[i] = make([]float64, nv+1)
		a := amin + float64(i)*da
		for j := 0; j < nv+1; j++ {
			b := float64(j) * db
			if cup {
				X[i][j] = c[0] + r*math.Cos(a)*math.Sin(b)
				Y[i][j] = c[1] + r*math.Sin(a)*math.Sin(b)
				Z[i][j] = c[2] - r*math.Cos(b)
			} else {
				X[i][j] = c[0] + r*math.Cos(a)*math.Sin(b)
				Y[i][j] = c[1] + r*math.Sin(a)*math.Sin(b)
				Z[i][j] = c[2] + r*math.Cos(b)
			}
		}
	}
	addSurfAndOrWire(X, Y, Z, args)
	return
}

// CylinderZ draws cylinder aligned with the z axis
//  Input:
//     alphaDeg -- half opening angle in degrees
//     height -- height of cone
//     nu -- number of divisions along the height of cone; e.g. 11
//     nv -- number of divisions along circumference of cone; e.g. 21
//   Output:
//     X, Y, Z -- the coordinages of all points as in a meshgrid
func CylinderZ(radius, height float64, nu, nv int, args *A) (X, Y, Z [][]float64) {
	X = make([][]float64, nu)
	Y = make([][]float64, nu)
	Z = make([][]float64, nu)
	for i := 0; i < nu; i++ {
		X[i] = make([]float64, nv+1)
		Y[i] = make([]float64, nv+1)
		Z[i] = make([]float64, nv+1)
		for j := 0; j < nv+1; j++ {
			h := height * float64(i) / float64(nu-1)
			θ := 2.0 * math.Pi * float64(j) / float64(nv)
			X[i][j] = radius * math.Cos(θ)
			Y[i][j] = radius * math.Sin(θ)
			Z[i][j] = h
		}
	}
	addSurfAndOrWire(X, Y, Z, args)
	return
}

// ConeZ draws cone aligned with the z axis
//  Input:
//     alphaDeg -- half opening angle in degrees
//     height -- height of cone
//     nu -- number of divisions along the height of cone; e.g. 11
//     nv -- number of divisions along circumference of cone; e.g. 21
//   Output:
//     X, Y, Z -- the coordinages of all points as in a meshgrid
func ConeZ(alphaDeg float64, height float64, nu, nv int, args *A) (X, Y, Z [][]float64) {
	r := math.Tan(alphaDeg*math.Pi/180.0) * height
	X = make([][]float64, nu)
	Y = make([][]float64, nu)
	Z = make([][]float64, nu)
	for i := 0; i < nu; i++ {
		X[i] = make([]float64, nv+1)
		Y[i] = make([]float64, nv+1)
		Z[i] = make([]float64, nv+1)
		for j := 0; j < nv+1; j++ {
			h := height * float64(i) / float64(nu-1)
			θ := 2.0 * math.Pi * float64(j) / float64(nv)
			X[i][j] = h * r * math.Cos(θ)
			Y[i][j] = h * r * math.Sin(θ)
			Z[i][j] = h
		}
	}
	addSurfAndOrWire(X, Y, Z, args)
	return
}

// ConeDiag draws cone parallel to the diagonal of the 3d space
//  Input:
//     alphaDeg -- half opening angle in degrees
//     height -- height of cone; i.e. length along space diagonal
//     nu -- number of divisions along the height of cone; e.g. 11
//     nv -- number of divisions along circumference of cone; e.g. 21
//   Output:
//     X, Y, Z -- the coordinages of all points as in a meshgrid
func ConeDiag(alphaDeg float64, height float64, nu, nv int, args *A) (X, Y, Z [][]float64) {
	r := math.Tan(alphaDeg*math.Pi/180.0) * height
	SQ2, SQ3, SQ6 := math.Sqrt2, math.Sqrt(3.0), math.Sqrt(6.0)
	X = make([][]float64, nu)
	Y = make([][]float64, nu)
	Z = make([][]float64, nu)
	for i := 0; i < nu; i++ {
		X[i] = make([]float64, nv+1)
		Y[i] = make([]float64, nv+1)
		Z[i] = make([]float64, nv+1)
		for j := 0; j < nv+1; j++ {
			h := height * float64(i) / float64(nu-1)
			θ := 2.0 * math.Pi * float64(j) / float64(nv)
			a := h * r * math.Cos(θ)
			b := h * r * math.Sin(θ)
			X[i][j] = (SQ2*h - b + SQ3*a) / SQ6
			Y[i][j] = (SQ2*h - b - SQ3*a) / SQ6
			Z[i][j] = (h + SQ2*b) / SQ3
		}
	}
	addSurfAndOrWire(X, Y, Z, args)
	return
}

// auxiliary ///////////////////////////////////////////////////////////////////////////////////////

// CalcDiagAngle computes the angle between a point and the diagonal of the 3d space
//   p -- point coordinates
//   returns the angle in radians
func CalcDiagAngle(p []float64) (alphaRad float64) {
	den := p[0] + p[1] + p[2]
	if den < 1e-10 {
		return 0.0
	}
	return math.Sqrt(math.Pow(p[0]-p[1], 2.0)+math.Pow(p[1]-p[2], 2.0)+math.Pow(p[2]-p[0], 2.0)) / den
}

// createAxes3d creates Python Axes3D if not yet created
func createAxes3d() {
	if !axes3dCreated {
		io.Ff(&bufferPy, "AX3D = plt.gcf().add_subplot(111, projection='3d')\n")
		io.Ff(&bufferPy, "AX3D.set_xlabel('x');AX3D.set_ylabel('y');AX3D.set_zlabel('z')\n")
		io.Ff(&bufferPy, "addToEA(AX3D)\n")
		axes3dCreated = true
	}
}
