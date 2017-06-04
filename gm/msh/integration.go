// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Integrator implements methods to perform numerical integration over a polyhedron/polygon
type Integrator struct {

	// input data
	Ctype  int         // cell type index
	Nverts int         // number of vertices = len(X)
	Ndim   int         // space dimension = len(X[0]) == len(P[0])
	Npts   int         // number of integration points = len(P)
	P      [][]float64 // (Gauss) integration points [npts][ndim]

	// slices related to integration points
	ShapeFcns [][]float64   // shape functions Sm @ all integ points [npts][nverts]
	RefGrads  [][][]float64 // reference gradients gm = dSm(r)/dr @ all integ points [npts][nverts][ndim]

	// mutable-data (scratchpad)
	JacobianMat [][]float64 // jacobian matrix Jr of the mapping reference to general coords [ndim][ndim]
	InvJacobMat [][]float64 // inverse of jacobian matrix [ndim][ndim]
	DetJacobian float64     // determinat of jacobian matrix
}

// NewIntegrator returns a new object to integrate over polyhedra/polygons (cells)
//   ctype -- index of cell type; e.g. TypeQuad4
//   P     -- integration points [npoints][ndim]. may be nil => default will be selected
//   pName -- use integration points from database instead of P or default ones. may be ""
func NewIntegrator(ctype int, P [][]float64, pName string) (o *Integrator, err error) {

	// check
	if ctype < 0 || ctype > TypeNumMax {
		err = chk.Err("ctype=%d is invalid; it must be in [0,%d]\n", ctype, TypeNumMax)
		return
	}

	// create new object
	o = new(Integrator)
	o.Ctype = ctype
	o.Nverts = NumVerts[ctype]
	o.Ndim = GeomNdim[ctype]

	// set integration points and related slices
	err = o.ResetP(P, pName)
	if err != nil {
		return
	}

	// allocate mutable data
	o.JacobianMat = la.MatAlloc(o.Ndim, o.Ndim)
	o.InvJacobMat = la.MatAlloc(o.Ndim, o.Ndim)
	return
}

// ResetP resets integration points
//   P     -- integration points [npoints][ndim]. may be nil => default will be selected
//   pName -- use integration points from database instead of P or default ones. may be ""
func (o *Integrator) ResetP(P [][]float64, pName string) (err error) {

	// set integration points
	if P != nil { // use given slice
		o.P = P
	} else if pName != "" { // find in database
		o.P, err = IntPointsFindSet(TypeIndexToKind[o.Ctype], pName)
		if err != nil {
			return
		}
	} else { // set default
		o.P = DefaultIntPoints[o.Ctype]
	}
	o.Npts = len(o.P)

	// allocate slices related to integration points
	if len(o.ShapeFcns) != o.Npts {
		o.ShapeFcns = la.MatAlloc(o.Npts, o.Nverts)
		o.RefGrads = utl.Deep3alloc(o.Npts, o.Nverts, o.Ndim)
	}

	// compute shape and reference gradient @ all integ points
	for ip, point := range o.P {
		Functions[o.Ctype](o.ShapeFcns[ip], o.RefGrads[ip], point, true)
	}
	return
}

// GetXip calculates coordinates Xip of integration points from X and P
//   Input:
//     X -- coordinates of vertices of cell (polyhedron/polygon) [nverts][ndim]
//   Output:
//     Xip -- general (non-reference) coordinate of integ points [npts][ndim]
func (o *Integrator) GetXip(X [][]float64) (Xip [][]float64) {
	Xip = make([][]float64, o.Npts)
	for i := 0; i < o.Npts; i++ {
		Xip[i] = make([]float64, o.Ndim)
		for j := 0; j < o.Ndim; j++ {
			for m := 0; m < o.Nverts; m++ {
				Xip[i][j] += o.ShapeFcns[i][m] * X[m][j]
			}
		}
	}
	return
}

// IntegrateSv integrates scalar function of vector argument over Cell
//
//   Computes:
//
//           ⌠⌠⌠   →       ⌠⌠⌠   → →     →        n-1   →  →      →
//     res = │││ f(x) dΩ = │││ f(x(r))⋅J(r) dΩr ≈  Σ  f(xi(ri))⋅J(ri)⋅wi
//           ⌡⌡⌡           ⌡⌡⌡                    i=0
//              Ω             Ωr
//   where:
//            → →    m-1   m →     m
//            x(r) ≈  Σ   S (r) ⋅ x               J = det(Jmat)
//                   i=0
//   and:
//            m -- number of cell nodes
//            n -- number of integration points
//   Input:
//     X  -- coordinates of vertices of cell (polyhedron/polygon) [nverts][ndim]
//     f  -- integrand function
func (o *Integrator) IntegrateSv(X [][]float64, f fun.Sv) (res float64, err error) {
	xip := make([]float64, o.Ndim)
	var fx float64
	for ip, point := range o.P {
		for j := 0; j < o.Ndim; j++ {
			xip[j] = 0
			for m := 0; m < o.Nverts; m++ {
				xip[j] += o.ShapeFcns[ip][m] * X[m][j]
			}
		}
		err = o.EvalJacobian(X, ip)
		if err != nil {
			return
		}
		fx, err = f(xip)
		if err != nil {
			return
		}
		res += fx * o.DetJacobian * point[3]
	}
	return
}

// EvalJacobian computes the Jacobian of the mapping from general to reference space
// at integration point with index ip
//
//           →     _                           _
//          dx    |  ∂x0/∂r0  ∂x0/∂r1  ∂x0/∂r2  |                ∂xi
//   Jmat = —— =  |  ∂x1/∂r0  ∂x1/∂r1  ∂x1/∂r2  |     Jmat[ij] = ———
//           →    |_ ∂x2/∂r0  ∂x2/∂r1  ∂x2/∂r2 _|                ∂rj
//          dr
//
//   Input:
//     X  -- coordinates of vertices of cell (polyhedron/polygon) [nverts][ndim]
//     ip -- index of integration point
//   Computed (stored):
//     JacobianMat -- reference Jacobian matrix [ndim][ndim]
//     InvJacobMat -- inverse of Jmat [ndim][ndim]
//     DetJacobian -- determinat of the reference Jacobian matrix
func (o *Integrator) EvalJacobian(X [][]float64, ip int) (err error) {
	if ip < 0 || ip > o.Npts {
		chk.Err("index of integration point %d is invalid. ip must be in [0,%d]\n", ip, o.Npts)
		return
	}
	if o.Ndim == 1 {
		// TODO
		return
	}
	for i := 0; i < o.Ndim; i++ {
		for j := 0; j < o.Ndim; j++ {
			o.JacobianMat[i][j] = 0
			for m := 0; m < o.Nverts; m++ {
				o.JacobianMat[i][j] += X[m][i] * o.RefGrads[ip][m][j]
			}
		}
	}
	o.DetJacobian, err = la.MatInv(o.InvJacobMat, o.JacobianMat, 1e-14)
	return
}

// mesh integrator ////////////////////////////////////////////////////////////////////////////////

// MeshIntegrator implements methods to perform numerical integration over a mesh
type MeshIntegrator struct {
	M           *Mesh           // the mesh
	Ngoroutines int             // total number of go routines
	Integrators [][]*Integrator // all integrators [Ngoroutines][TypeNumMax]
}

// NewMeshIntegrator returns a new MeshIntegrator
func NewMeshIntegrator(mesh *Mesh, Ngoroutines int) (o *MeshIntegrator, err error) {

	// check
	if Ngoroutines < 1 {
		err = chk.Err("number of goroutines must be at least 1\n")
		return
	}

	// allocate integrators
	o = new(MeshIntegrator)
	o.M = mesh
	o.Integrators = make([][]*Integrator, Ngoroutines)
	for i := 0; i < Ngoroutines; i++ {
		o.Integrators[i] = make([]*Integrator, TypeNumMax)
		for j := 0; j < TypeNumMax; j++ {
			o.Integrators[i][j], err = NewIntegrator(j, nil, "")
			if err != nil {
				return
			}
		}
	}
	return
}

// IntegrateSv integrates scalar function of vector argument over mesh
//
//           ⌠⌠⌠   →
//     res = │││ f(x) dΩ
//           ⌡⌡⌡
//              Ω
//   Input:
//     goroutineId -- go routine id to use when performing optimisation (not to partition mesh)
func (o *MeshIntegrator) IntegrateSv(goroutineId int, f fun.Sv) (res float64, err error) {
	for _, c := range o.M.Cells {
		r, e := o.Integrators[goroutineId][c.TypeIndex].IntegrateSv(c.X, f)
		if e != nil {
			err = e
			return
		}
		res += r
	}
	return
}
