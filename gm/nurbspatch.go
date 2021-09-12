// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"encoding/json"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// PointExchangeData holds data for exchanging control points
type PointExchangeData struct {
	ID  int       `json:"i"` // id
	Tag int       `json:"t"` // tag
	X   []float64 `json:"x"` // coordinates (size==4)
}

// NurbsExchangeData holds all data required to exchange NURBS; e.g. read/save files
type NurbsExchangeData struct {
	ID    int         `json:"i"` // id of Nurbs
	Gnd   int         `json:"g"` // 1: curve, 2:surface, 3:volume (geometry dimension)
	Ords  []int       `json:"o"` // order along each x-y-z direction [gnd]
	Knots [][]float64 `json:"k"` // knots along each x-y-z direction [gnd][m]
	Ctrls []int       `json:"c"` // global ids of control points
}

// NurbsExchangeDataSet defines a set of nurbs exchange data
type NurbsExchangeDataSet []*NurbsExchangeData

// NurbsPatch defines a patch of many NURBS'
type NurbsPatch struct {

	// input/output data
	ControlPoints []*PointExchangeData `json:"points"` // input/output control points
	ExchangeData  NurbsExchangeDataSet `json:"patch"`  // input/output nurbs data

	// Nurbs structures
	Entities []*Nurbs `json:"-"` // pointers to NURBS

	// auxiliary
	Bins Bins `json:"-"` // auxiliary structure to locate points
}

// initialization functions ///////////////////////////////////////////////////////////////////////

// NewNurbsPatch returns new patch of NURBS
//   tolerance -- tolerance to assume that two control points are the same
func NewNurbsPatch(binsNdiv int, tolerance float64, entities ...*Nurbs) (o *NurbsPatch) {
	o = new(NurbsPatch)
	o.Entities = entities
	o.ResetFromEntities(binsNdiv, tolerance)
	return
}

// ResetFromEntities will reset all exchange data with information from Entities slice
func (o *NurbsPatch) ResetFromEntities(binsNdiv int, tolerance float64) {

	// limits and ndim
	xmin, xmax, ndim := o.LimitsAndNdim()

	// reset bins
	o.Bins.Clear()
	o.Bins.Init(xmin[:ndim], xmax[:ndim], utl.IntVals(ndim, binsNdiv))

	// global index of control point
	nextID := 0

	// set exchange data
	o.ExchangeData = make([]*NurbsExchangeData, len(o.Entities))
	for e, entity := range o.Entities {

		// set entity and compute number of control points
		nctrls := entity.n[0] * entity.n[1] * entity.n[2]

		// new data structure
		o.ExchangeData[e] = &NurbsExchangeData{
			ID:    e,
			Gnd:   entity.gnd,
			Ords:  make([]int, entity.gnd),
			Knots: make([][]float64, entity.gnd),
			Ctrls: make([]int, nctrls),
		}

		// set orders and knots
		for d := 0; d < entity.gnd; d++ {
			nknots := len(entity.b[d].T)
			o.ExchangeData[e].Ords[d] = entity.p[d]
			o.ExchangeData[e].Knots[d] = make([]float64, nknots)
			for k, u := range entity.b[d].T {
				o.ExchangeData[e].Knots[d][k] = u
			}
		}

		// set control points
		idxCtrl := 0
		for k := 0; k < entity.n[2]; k++ {
			for j := 0; j < entity.n[1]; j++ {
				for i := 0; i < entity.n[0]; i++ {

					// get point Id
					x := entity.GetQ(i, j, k)
					id, existent := o.Bins.FindClosestAndAppend(&nextID, x, nil, tolerance, o.diffPoints)

					// set control point id
					o.ExchangeData[e].Ctrls[idxCtrl] = id
					idxCtrl++

					// set list of points
					if !existent {
						o.ControlPoints = append(o.ControlPoints, &PointExchangeData{ID: id, X: x})
					}
				}
			}
		}
	}
}

// ResetFromExchangeData will reset all Entities with information from ExchangeData (and ControlPoints)
func (o *NurbsPatch) ResetFromExchangeData(binsNdiv int, tolerance float64) {

	// check
	if len(o.ExchangeData) < 1 || len(o.ControlPoints) < 1 {
		chk.Panic("there are no ExchangeData or ControlPoints\n")
	}

	// collect vertices
	verts := make([][]float64, len(o.ControlPoints))
	for i, cp := range o.ControlPoints {
		verts[i] = cp.X
	}

	// allocate nurbs
	o.Entities = make([]*Nurbs, len(o.ExchangeData))
	for i, ed := range o.ExchangeData {
		o.Entities[i] = NewNurbs(ed.Gnd, ed.Ords, ed.Knots)
		o.Entities[i].SetControl(verts, ed.Ctrls)
	}

	// limits and ndim
	xmin, xmax, ndim := o.LimitsAndNdim()

	// reset bins
	o.Bins.Clear()
	o.Bins.Init(xmin[:ndim], xmax[:ndim], utl.IntVals(ndim, binsNdiv))

	// global index of control point
	nextID := 0

	// set bins
	for _, entity := range o.Entities {
		for k := 0; k < entity.n[2]; k++ {
			for j := 0; j < entity.n[1]; j++ {
				for i := 0; i < entity.n[0]; i++ {
					x := entity.GetQ(i, j, k)
					o.Bins.FindClosestAndAppend(&nextID, x, nil, tolerance, o.diffPoints)
				}
			}
		}
	}
}

// read/write function ////////////////////////////////////////////////////////////////////////////

// Write writes ExchangeData to json file
func (o NurbsPatch) Write(dirout, fnkey string) {
	b, err := json.Marshal(o)
	if err != nil {
		chk.Panic("%v\n", err)
	}
	io.WriteBytesToFileVD(dirout, fnkey+".json", b)
}

// NewNurbsPatchFromFile allocates a NurbsPatch with data from file
func NewNurbsPatchFromFile(filename string, binsNdiv int, tolerance float64) (o *NurbsPatch) {

	// read exchange data
	b := io.ReadFile(filename)
	o = new(NurbsPatch)
	err := json.Unmarshal(b, o)
	if err != nil {
		chk.Panic("%v\n", err)
	}

	// allocate nurbs
	o.ResetFromExchangeData(binsNdiv, tolerance)
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// LimitsAndNdim computes the limits of the patch and max dimension by looping over all Entities
func (o NurbsPatch) LimitsAndNdim() (xmin, xmax []float64, ndim int) {

	// check
	if len(o.Entities) < 1 {
		chk.Panic("there are no Entities\n")
	}

	// find limits
	xmin, xmax = o.Entities[0].GetLimitsQ()
	for i := o.Entities[0].gnd; i < 3; i++ {
		xmin[i] = math.Inf(+1)
		xmax[i] = math.Inf(-1)
	}
	for i := 1; i < len(o.Entities); i++ {
		xmi, xma := o.Entities[i].GetLimitsQ()
		for i := 0; i < 3; i++ {
			xmin[i] = utl.Min(xmin[i], xmi[i])
			xmax[i] = utl.Max(xmax[i], xma[i])
		}
	}

	// find max space dimension; e.g. we may have curves and surfaces at the same time
	ndim = 0
	for i := 0; i < 3; i++ {
		if xmax[i]-xmin[i] > XDELZERO {
			ndim++
		}
	}
	if ndim < 1 {
		chk.Panic("ndim=%d is invalid\n", ndim)
	}
	return
}

// auxiliary internal /////////////////////////////////////////////////////////////////////////////

// diffPoints returns true if two control points are different, considering the weights
func (o NurbsPatch) diffPoints(idOld int, xNew []float64) bool {
	if math.Abs(xNew[3]-o.ControlPoints[idOld].X[3]) > 0 {
		return true
	}
	return false
}
