// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/gm/msh"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// EssentialBcs holds data for prescribing a SET of essential (Dirichlet) boundary conditions
type EssentialBcs struct {
	all     [][]dbf.T    // [node][dof] function to compute BCs
	grid    *gm.Grid     // using grid
	mesh    *msh.Mesh    // using mesh
	maxNdof int          // max number of "degrees-of-freedom" per node
	nodes   map[int]bool // list of nodes with prescribed boundary conditions
}

// NewEssentialBcsGrid returns a new EssentialBcs structure using Grid
//  grid    -- grid
//  maxNdof -- max number of "degrees-of-freedom" per node
func NewEssentialBcsGrid(grid *gm.Grid, maxNdof int) (o *EssentialBcs) {
	o = new(EssentialBcs)
	o.all = make([][]dbf.T, grid.Size())
	o.grid = grid
	o.maxNdof = maxNdof
	o.nodes = make(map[int]bool)
	return
}

// NewEssentialBcsMesh returns a new EssentialBcs structure using Mesh
//  mesh    -- mesh
//  maxNdof -- max number of "degrees-of-freedom" per node
func NewEssentialBcsMesh(mesh *msh.Mesh, maxNdof int) (o *EssentialBcs) {
	o = new(EssentialBcs)
	o.all = make([][]dbf.T, len(mesh.Verts))
	o.mesh = mesh
	o.maxNdof = maxNdof
	o.nodes = make(map[int]bool)
	return
}

// AddUsingTag sets boundary condition using edge or face tag from grid or mesh
//   tag    -- edge or face tag
//   dof    -- index of "degree-of-freedom"; e.g. 0⇒horizontal displacement, 1⇒vertical displacement
//   cvalue -- constant value [optional]; or
//   fvalue -- function value [optional]
func (o *EssentialBcs) AddUsingTag(tag, dof int, cvalue float64, fvalue dbf.T) {

	// check
	if dof > o.maxNdof-1 {
		chk.Panic("cannot set dof=%d because maxNdof=%d\n", dof, o.maxNdof)
	}

	// function
	f := fvalue
	if fvalue == nil {
		f = &dbf.Cte{cvalue}
	}

	// using grid
	var nodes []int
	if o.grid != nil {
		nodes = o.grid.Boundary(tag)
	} else {

		// using mesh
		if o.mesh == nil {
			chk.Panic("mesh is required if not using grid\n")
		}
		nodes = o.mesh.Boundary(tag)
	}

	// check
	if nodes == nil {
		chk.Panic("cannot find nodes with tag=%d\n", tag)
	}

	// set
	for _, n := range nodes {
		if o.all[n] == nil {
			o.all[n] = make([]dbf.T, o.maxNdof)
		}
		o.all[n][dof] = f
		o.nodes[n] = true
	}
}

// Nodes returns (unique/sorted) list of nodes with prescribed boundary conditions
func (o *EssentialBcs) Nodes() []int {
	return utl.IntBoolMapSort(o.nodes)
}

// Value returns the value of prescribed boundary condition @ {node,dof,time}
func (o *EssentialBcs) Value(node, dof int, t float64) (val float64, available bool) {

	// check if available
	bc := o.all[node]
	if bc == nil {
		return
	}
	if bc[dof] == nil {
		return
	}

	// using grid
	if o.grid != nil {
		return bc[dof].F(t, o.grid.Node(node)), true
	}

	// using mesh
	return bc[dof].F(t, o.mesh.Verts[node].X), true
}

// Print prints boundary conditions
func (o *EssentialBcs) Print() (l string) {
	var strNid string
	if o.grid != nil {
		_, strNid = utl.Digits(o.grid.Size())
	}
	if o.mesh != nil {
		_, strNid = utl.Digits(len(o.mesh.Verts))
	}
	_, strDof := utl.Digits(o.maxNdof)
	for _, n := range o.Nodes() {
		list := ""
		for dof := 0; dof < o.maxNdof; dof++ {
			val, available := o.Value(n, dof, 0)
			if available {
				list += io.Sf("  dof="+strDof+" value=%g", dof, val)
			}
		}
		l += io.Sf("node = "+strNid+list+"\n", n)
	}
	return
}
