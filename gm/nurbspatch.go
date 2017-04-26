// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

// PointExchangeData holds data for exchanging control points
type PointExchangeData struct {
	Id  int       `json:"i"` // id
	Tag int       `json:"t"` // tag
	Q   []float64 `json:"q"` // coordinates (size==4)
}

// NurbsExchangeData holds all data required to exchange NURBS; e.g. read/save files
type NurbsExchangeData struct {
	Id    int         `json:"i"` // id of Nurbs
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

// NewNurbsPatch returns new patch of NURBS
//   tolerance -- tolerance to assume that two control points are the same
func NewNurbsPatch(binsNdiv int, tolerance float64, entities ...*Nurbs) (o *NurbsPatch, err error) {

	// allocate
	o = new(NurbsPatch)
	if len(entities) < 1 {
		return
	}

	// find limits
	xmin, xmax := entities[0].GetLimitsQ()
	for i := entities[0].gnd; i < 3; i++ {
		xmin[i] = math.Inf(+1)
		xmax[i] = math.Inf(-1)
	}
	for i := 1; i < len(entities); i++ {
		xmi, xma := entities[i].GetLimitsQ()
		for i := 0; i < 3; i++ {
			xmin[i] = utl.Min(xmin[i], xmi[i])
			xmax[i] = utl.Max(xmax[i], xma[i])
		}
	}

	// find max space dimension; e.g. we may have curves and surfaces at the same time
	ndim := 0
	for i := 0; i < 3; i++ {
		if xmax[i]-xmin[i] > XDELZERO {
			ndim++
		}
	}
	if ndim < 1 {
		err = chk.Err("ndim=%d is invalid", ndim)
		return
	}

	// allocate auxiliary structure to locate points
	err = o.Bins.Init(xmin[:ndim], xmax[:ndim], binsNdiv)
	if err != nil {
		return
	}

	// global index of control point
	nextId := 0

	// set exchange data
	o.Entities = make([]*Nurbs, len(entities))
	o.ExchangeData = make([]*NurbsExchangeData, len(entities))
	for e, entity := range entities {

		// set entity and compute number of control points
		o.Entities[e] = entity
		nctrls := entity.n[0] * entity.n[1] * entity.n[2]

		// new data structure
		o.ExchangeData[e] = &NurbsExchangeData{
			Id:    e,
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
					id := o.Bins.FindClosestAndAppend(&nextId, x, nil, tolerance)

					// set control point id
					o.ExchangeData[e].Ctrls[idxCtrl] = id
					idxCtrl++
				}
			}
		}
	}
	return
}
