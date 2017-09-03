// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// EssentialBc holds data for prescribing ONE essential (Dirichlet) boundary conditions
type EssentialBc struct {
	Tag    int          // tag of node, edge or face to which this BCS is prescribed
	Key    string       // name of bcs; e.g. "ux" (horizontal displacement) or "pl" (liquid pressure)
	Cvalue float64      // constant value of boundary condition (constant in time)
	Fvalue dbf.T        // function f(t,{x}) [may be nil ⇒ use constant value]
	Nodes  map[int]bool // list of nodes with this BCS
}

// GetNodesSorted returns sorted list of nodes with this boundary condition
func (o *EssentialBc) GetNodesSorted() []int {
	return utl.IntBoolMapSort(o.Nodes)
}

// Value returns the constant value (cvalue) or function value @ t (time)
func (o *EssentialBc) Value(t float64) float64 {
	if o.Fvalue == nil {
		return o.Cvalue
	}
	return o.Fvalue.F(t, nil)
}

// uuidT defines a type for unique identifier
type uuidT struct {
	tag int
	key string
}

// EssentialBcs holds data for prescribing a SET of essential (Dirichlet) boundary conditions
type EssentialBcs struct {
	All    []*EssentialBc
	Finder map[uuidT]int
}

// NewEssentialBcs returns a new EssentialBcs structure
func NewEssentialBcs() (o *EssentialBcs) {
	o = new(EssentialBcs)
	o.All = make([]*EssentialBc, 0, 1000)
	o.Finder = make(map[uuidT]int)
	return
}

// SetInGrid sets boundary condition considering Grid data
func (o *EssentialBcs) SetInGrid(g *gm.Grid, tag int, key string, cvalue float64, fvalue dbf.T) (err error) {
	nodes := g.Boundary(tag)
	if nodes == nil {
		return chk.Err("cannot find nodes with tag=%d in grid\n", tag)
	}
	uuid := uuidT{tag, key}
	if idx, ok := o.Finder[uuid]; ok {
		bc := o.All[idx]
		for _, n := range nodes {
			bc.Nodes[n] = true
		}
	} else {
		o.Finder[uuid] = len(o.All)
		nmap := make(map[int]bool)
		for _, n := range nodes {
			nmap[n] = true
		}
		o.All = append(o.All, &EssentialBc{tag, key, cvalue, fvalue, nmap})
	}
	return
}

// GetNodesList returns (unique/sorted) list of all nodes
func (o *EssentialBcs) GetNodesList() (allNodes []int) {
	var list []int
	for _, bc := range o.All {
		list = append(list, bc.GetNodesSorted()...)
	}
	return utl.IntUnique(list)
}

// Print prints boundary conditions
func (o *EssentialBcs) Print() (l string) {
	for _, bc := range o.All {
		l += "------------------------------------------------------------------------------------\n"
		l += io.Sf("tag    = %v\n", bc.Tag)
		l += io.Sf("key    = %v\n", bc.Key)
		if bc.Fvalue != nil {
			l += io.Sf("fvalue = <given function>\n")
		} else {
			l += io.Sf("cvalue = %v\n", bc.Cvalue)
		}
		l += io.Sf("nodes  = %v\n", bc.GetNodesSorted())
	}
	return
}
