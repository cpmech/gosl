// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// ExtractSurfaces returns a new NURBS representing a boundary of this NURBS
func (o *Nurbs) ExtractSurfaces() (surfs []*Nurbs) {
	if o.gnd == 1 {
		return
	}
	nsurf := o.gnd * 2
	surfs = make([]*Nurbs, nsurf)
	var ords [][]int
	var knots [][][]float64
	if o.gnd == 2 {
		ords = [][]int{
			{o.p[1]}, // perpendicular to x
			{o.p[0]}, // perpendicular to y
		}
		knots = [][][]float64{
			{o.b[1].T}, // perpendicular to x
			{o.b[0].T}, // perpendicular to y
		}
	} else {
		ords = [][]int{
			{o.p[1], o.p[2]}, // perpendicular to x
			{o.p[2], o.p[0]}, // perpendicular to y
			{o.p[0], o.p[1]}, // perpendicular to z
		}
		knots = [][][]float64{
			{o.b[1].T, o.b[2].T}, // perpendicular to x
			{o.b[2].T, o.b[0].T}, // perpendicular to y
			{o.b[0].T, o.b[1].T}, // perpendicular to z
		}
	}
	for i := 0; i < o.gnd; i++ {
		a, b := i*o.gnd, i*o.gnd+1
		surfs[a] = NewNurbs(o.gnd-1, ords[i], knots[i]) // surface perpendicular to i
		surfs[b] = NewNurbs(o.gnd-1, ords[i], knots[i]) // opposite surface perpendicular to i
		if o.gnd == 2 {                                 // boundary is curve
			j := (i + 1) % o.gnd // direction perpendicular to i
			surfs[a].Q = o.CloneCtrlsAlongCurve(j, 0)
			surfs[b].Q = o.CloneCtrlsAlongCurve(j, o.n[i]-1)
		} else { // boundary is surface
			j := (i + 1) % o.gnd // direction perpendicular to i
			k := (i + 2) % o.gnd // other direction perpendicular to i
			surfs[a].Q = o.CloneCtrlsAlongSurface(j, k, 0)
			surfs[b].Q = o.CloneCtrlsAlongSurface(j, k, o.n[i]-1)
		}
	}
	return
}

// CloneCtrlsAlongCurve returns a copy of control points @ 2D boundary
func (o *Nurbs) CloneCtrlsAlongCurve(iAlong, jAt int) (Qnew [][][][]float64) {
	Qnew = utl.Deep4alloc(o.n[iAlong], 1, 1, 4)
	var i, j int
	for m := 0; m < o.n[iAlong]; m++ {
		i, j = m, jAt
		if iAlong == 1 {
			i, j = jAt, m
		}
		for e := 0; e < 4; e++ {
			Qnew[m][0][0][e] = o.Q[i][j][0][e]
		}
	}
	return
}

// CloneCtrlsAlongSurface returns a copy of control points @ 3D boundary
func (o *Nurbs) CloneCtrlsAlongSurface(iAlong, jAlong, kat int) (Qnew [][][][]float64) {
	Qnew = utl.Deep4alloc(o.n[iAlong], o.n[jAlong], 1, 4)
	var i, j, k int
	for m := 0; m < o.n[iAlong]; m++ {
		for n := 0; n < o.n[jAlong]; n++ {
			switch {
			case iAlong == 0 && jAlong == 1:
				i, j, k = m, n, kat
			case iAlong == 1 && jAlong == 2:
				i, j, k = kat, m, n
			case iAlong == 2 && jAlong == 0:
				i, j, k = n, kat, m
			default:
				chk.Panic("clone Q surface is specified by 'along' indices in (0,1) or (1,2) or (2,0). (%d,%d) is incorrect", iAlong, jAlong)
			}
			for e := 0; e < 4; e++ {
				Qnew[m][n][0][e] = o.Q[i][j][k][e]
			}
		}
	}
	return
}

// IndsAlongCurve returns the control points indices along curve
func (o *Nurbs) IndsAlongCurve(iAlong, iSpan0, jAt int) (L []int) {
	nb := o.p[iAlong] + 1 // number of basis along i
	L = make([]int, nb)
	var i, j int
	for m := 0; m < nb; m++ {
		if iAlong == 0 {
			i = iSpan0 - o.p[0] + m
			j = jAt
		} else {
			i = jAt
			j = iSpan0 - o.p[1] + m
		}
		L[m] = i + j*o.n[0]
	}
	return
}

// IndsAlongSurface return the control points indices along surface
func (o *Nurbs) IndsAlongSurface(iAlong, jAlong, iSpan0, jSpan0, kat int) (L []int) {
	nbu := o.p[iAlong] + 1 // number of basis functions along i
	nbv := o.p[jAlong] + 1 // number of basis functions along j
	L = make([]int, nbu*nbv)
	var c, i, j, k int
	for m := 0; m < nbu; m++ {
		for n := 0; n < nbv; n++ {
			switch {
			case iAlong == 0 && jAlong == 1:
				i = iSpan0 - o.p[0] + m
				j = jSpan0 - o.p[1] + n
				k = kat
			case iAlong == 1 && jAlong == 2:
				i = kat
				j = iSpan0 - o.p[1] + m
				k = jSpan0 - o.p[2] + n
			case iAlong == 2 && jAlong == 0:
				i = jSpan0 - o.p[0] + n
				j = kat
				k = iSpan0 - o.p[2] + m
			}
			L[c] = i + j*o.n[0] + k*o.n[1]*o.n[2]
			c++
		}
	}
	return
}

// ElemBryLocalInds returns the local (element) indices of control points @ boundaries
// (if element would have all surfaces @ boundaries)
func (o *Nurbs) ElemBryLocalInds() (I [][]int) {
	switch o.gnd {
	case 1:
		return
	case 2:
		I = make([][]int, 2*o.gnd)
		nx, ny := o.p[0]+1, o.p[1]+1
		I[3] = utl.IntRange3(0, nx*ny, nx)
		I[1] = utl.IntAddScalar(I[3], nx-1)
		I[0] = utl.IntRange(nx)
		I[2] = utl.IntAddScalar(I[0], (ny-1)*nx)
	case 3:
		I = make([][]int, 2*o.gnd)
		chk.Panic("3D NURBS: ElemBryLocalInds: TODO") // TODO
	}
	return
}
