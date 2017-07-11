// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/io"
)

// Point holds the Cartesian coordinates of a point in 3D space
type Point struct {
	X, Y, Z float64
}

// Segment represents a directed segment from A to B
type Segment struct {
	A, B *Point
}

// Point methods /////////////////////////////////////////////////////////////////////////////////////

// NewCopy creates a new copy of Point
func (o *Point) NewCopy() *Point {
	return &Point{o.X, o.Y, o.Z}
}

// NewDisp creates a new copy of Point displaced by dx, dy, dz
func (o *Point) NewDisp(dx, dy, dz float64) *Point {
	return &Point{o.X + dx, o.Y + dy, o.Z + dz}
}

// String outputs Point
func (o *Point) String() string {
	return io.Sf("{%g, %g, %g}", o.X, o.Y, o.Z)
}

// DistPointPoint computes the unsigned distance from a to b
func DistPointPoint(a, b *Point) float64 {
	return math.Sqrt((a.X-b.X)*(a.X-b.X) +
		(a.Y-b.Y)*(a.Y-b.Y) +
		(a.Z-b.Z)*(a.Z-b.Z))
}

// Segment methods ///////////////////////////////////////////////////////////////////////////////////

// Len computes the length of Segment == Euclidian norm
func (o *Segment) Len() float64 {
	return DistPointPoint(o.A, o.B)
}

// New creates a new Segment scaled by m and starting from A
func (o *Segment) New(m float64) *Segment {
	return &Segment{o.A.NewCopy(), &Point{o.A.X + m*(o.B.X-o.A.X),
		o.A.Y + m*(o.B.Y-o.A.Y),
		o.A.Z + m*(o.B.Z-o.A.Z)}}
}

// Vector returns the vector representing Segment from A to B (scaled by m)
func (o *Segment) Vector(m float64) []float64 {
	return []float64{m * (o.B.X - o.A.X),
		m * (o.B.Y - o.A.Y),
		m * (o.B.Z - o.A.Z)}
}

// String outputs Segment
func (o *Segment) String() string {
	return io.Sf("{%v %v} len=%g", o.A, o.B, o.Len())
}

// NewSegment creates a new segment from a to b
func NewSegment(a, b *Point) *Segment {
	return &Segment{a, b}
}

// Vector ////////////////////////////////////////////////////////////////////////////////////////////

// VecDot returns the dot product between two vectors
func VecDot(u, v []float64) float64 {
	return u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
}

// VecNorm returns the length (Euclidian norm) of a vector
func VecNorm(u []float64) float64 {
	return math.Sqrt(u[0]*u[0] + u[1]*u[1] + u[2]*u[2])
}

// VecNew returns a new vector scaled by m
func VecNew(m float64, u []float64) []float64 {
	return []float64{m * u[0], m * u[1], m * u[2]}
}

// VecNewAdd returns a new vector by adding two other vectors
//  w := α*u + β*v
func VecNewAdd(α float64, u []float64, β float64, v []float64) []float64 {
	return []float64{α*u[0] + β*v[0], α*u[1] + β*v[1], α*u[2] + β*v[2]}
}

// distance functions ////////////////////////////////////////////////////////////////////////////////

// DistPointLine computes the distance from p to line passing through a -> b
func DistPointLine(p, a, b *Point, tol float64, verbose bool) float64 {
	ns := NewSegment(a, b)
	vs := NewSegment(p, a)
	nn := ns.Len()
	if nn < tol { // point-point distance
		if verbose {
			io.Pfred("basicgeom.go: DistPointLine: __WARNING__ point-point distance too small:\n p=%v a=%v b=%v\n", p, a, b)
		}
		return vs.Len()
	}
	n := ns.Vector(1.0 / nn)
	v := vs.Vector(1.0)
	s := VecDot(v, n)
	l := VecNewAdd(1, v, -s, n) // l := v - dot(v,n) * n
	return VecNorm(l)
}

// locate functions //////////////////////////////////////////////////////////////////////////////////

// PointsLims returns the limits of a set of points
func PointsLims(pp []*Point) (cmin, cmax []float64) {
	if len(pp) < 1 {
		return []float64{0, 0, 0}, []float64{0, 0, 0}
	}
	cmin = []float64{pp[0].X, pp[0].Y, pp[0].Z}
	cmax = []float64{pp[0].X, pp[0].Y, pp[0].Z}
	for i := 1; i < len(pp); i++ {
		if pp[i].X < cmin[0] {
			cmin[0] = pp[i].X
		}
		if pp[i].Y < cmin[1] {
			cmin[1] = pp[i].Y
		}
		if pp[i].Z < cmin[2] {
			cmin[2] = pp[i].Z
		}
		if pp[i].X > cmax[0] {
			cmax[0] = pp[i].X
		}
		if pp[i].Y > cmax[1] {
			cmax[1] = pp[i].Y
		}
		if pp[i].Z > cmax[2] {
			cmax[2] = pp[i].Z
		}
	}
	return
}

// IsPointIn returns whether p is inside box with cMin and cMax
func IsPointIn(p *Point, cMin, cMax []float64, tol float64) bool {
	if p.X < cMin[0]-tol || p.X > cMax[0]+tol {
		return false
	}
	if p.Y < cMin[1]-tol || p.Y > cMax[1]+tol {
		return false
	}
	if p.Z < cMin[2]-tol || p.Z > cMax[2]+tol {
		return false
	}
	return true
}

// IsPointInLine returns whether p is inside line passing through a and b
func IsPointInLine(p, a, b *Point, zero, told, tolin float64) bool {
	cmin, cmax := PointsLims([]*Point{a, b})
	d := DistPointLine(p, a, b, zero, false)
	if d < told && IsPointIn(p, cmin, cmax, tolin) {
		return true
	}
	return false
}

/*
// WireLength returns the length of a wireframe
func WireLength(pp []*Point) (l float64) {
    for xa, xb in zip(points[:-1], points[1:]) {
        l += DistPointPoint(xa, xb)
    }
    return l
}

// DistAlongWire returns the distance of p along oriented line
//  Input:
//   points -- list of points defining wireframe
//  Output:
//   d -- total distance from p to first point in wireframe,
//        negative value (-1.0) means point is outside line
func DistAlongWire(p, points, TolD=1e-10, TolIn=1e-10) {
    p, l, k = array(p, dtype=float), 0.0, 0
    for xa, xb in zip(points[:-1], points[1:]) {
        if IsPointInLine(p, xa, xb, TolD, TolIn) {
            return l + DistPointPoint(p, xa)
        }
        l += DistPointPoint(xa, xb)
    }
    return -1.0 // not found
}
*/
