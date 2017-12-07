// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"sort"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/gm/msh"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// BoundaryConds holds data for prescribing a SET of boundary conditions
type BoundaryConds struct {
	grid *gm.Grid    // using grid
	mesh *msh.Mesh   // using mesh
	ndof int         // max number of degrees of freedom
	fcns [][]fun.Svs // [...][dof] function to compute BCs; f({x}, t)
	tags [][]int     // [...][3] tag used to set BC; max 3 sides (e.g. corner node)
	n2i  []int       // [nnodesTotal] maps node ID to position in fcns and tags; -1 means not set
}

// NewBoundaryCondsGrid returns a new structure using Grid
func NewBoundaryCondsGrid(grid *gm.Grid, ndof int) (o *BoundaryConds) {
	o = new(BoundaryConds)
	o.grid = grid
	o.ndof = ndof
	o.n2i = make([]int, grid.Size())
	utl.IntFill(o.n2i, -1)
	return
}

// NewBoundaryCondsMesh returns a new structure using Mesh
func NewBoundaryCondsMesh(mesh *msh.Mesh, ndof int) (o *BoundaryConds) {
	o = new(BoundaryConds)
	o.mesh = mesh
	o.ndof = ndof
	o.n2i = make([]int, len(mesh.Verts))
	utl.IntFill(o.n2i, -1)
	return
}

// Has tells whether node has prescribed boundary condition or not
func (o *BoundaryConds) Has(node int) bool {
	return o.n2i[node] >= 0
}

// Tags returns tags used to prescrie boundary condition at node
// NOTE: returns empty list of node does not have boundary condition
func (o *BoundaryConds) Tags(node int) []int {
	if o.n2i[node] < 0 {
		return nil
	}
	return o.tags[o.n2i[node]]
}

// NormalGrid computes the outward normal at node (summed at corners) [when using Grid only]
// NOTE: returns zero normal if node doesn't have prescribed condition
func (o *BoundaryConds) NormalGrid(N la.Vector, node int) {
	if o.n2i[node] < 0 {
		N.Fill(0)
		return
	}
	tt := o.tags[o.n2i[node]]
	o.grid.UnitNormal(N, tt[0], node)
	if len(tt) > 1 {
		Ntmp := la.NewVector(o.grid.Ndim())
		for k := 1; k < len(tt); k++ {
			o.grid.UnitNormal(Ntmp, tt[k], node)
			for i := 0; i < o.grid.Ndim(); i++ {
				N[i] += Ntmp[i]
			}
		}
	}
}

// AddUsingTag sets boundary condition using edge or face tag from grid or mesh
//   tag    -- edge or face tag
//   dof    -- index of "degree-of-freedom"; e.g. 0⇒horizontal displacement, 1⇒vertical displacement
//   cvalue -- constant value [optional]; or
//   fvalue -- function value [optional]
func (o *BoundaryConds) AddUsingTag(tag, dof int, cvalue float64, fvalue fun.Svs) {

	// use or create function
	f := fvalue
	if fvalue == nil {
		f = func(x la.Vector, t float64) float64 { return cvalue }
	}

	// grid nodes
	var nodes []int
	if o.grid != nil {
		nodes = o.grid.Boundary(tag)
	}

	// or mesh nodes
	if o.mesh != nil {
		nodes = o.mesh.Boundary(tag)
	}

	// check
	if nodes == nil {
		chk.Panic("cannot find nodes with tag=%d\n", tag)
	}

	// set
	for _, n := range nodes {
		if o.n2i[n] < 0 { // new
			o.n2i[n] = len(o.fcns)
			ff := make([]fun.Svs, o.ndof)
			ff[dof] = f
			o.fcns = append(o.fcns, ff)
			o.tags = append(o.tags, []int{tag})
		} else { // existent
			o.fcns[o.n2i[n]][dof] = f
			o.tags[o.n2i[n]] = append(o.tags[o.n2i[n]], tag)
		}
	}
}

// Nodes returns (unique/sorted) list of nodes with prescribed boundary conditions
func (o *BoundaryConds) Nodes() (list []int) {
	list = make([]int, len(o.fcns))
	for n, i := range o.n2i {
		if i >= 0 {
			list[i] = n
		}
	}
	sort.Ints(list)
	return
}

// Value returns the value of prescribed boundary condition @ {node,dof,time}
func (o *BoundaryConds) Value(node, dof int, t float64) (tags []int, val float64, available bool) {

	// check if available
	i := o.n2i[node]
	if i < 0 {
		return
	}
	if o.fcns[i][dof] == nil {
		return
	}

	// using grid
	if o.grid != nil {
		return o.tags[i], o.fcns[i][dof](o.grid.Node(node), t), true
	}

	// using mesh
	return o.tags[i], o.fcns[i][dof](o.mesh.Verts[node].X, t), true
}

// Print prints boundary conditions
func (o *BoundaryConds) Print() (l string) {
	var strNid string
	if o.grid != nil {
		_, strNid = utl.Digits(o.grid.Size())
	}
	if o.mesh != nil {
		_, strNid = utl.Digits(len(o.mesh.Verts))
	}
	_, strDof := utl.Digits(o.ndof)
	for _, n := range o.Nodes() {
		list := ""
		for dof := 0; dof < o.ndof; dof++ {
			tags, val, available := o.Value(n, dof, 0)
			if available {
				list += io.Sf("  dof="+strDof+" tags=%v value=%g", dof, tags, val)
			}
		}
		l += io.Sf("node = "+strNid+list+"\n", n)
	}
	return
}
