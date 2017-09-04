// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pde

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/fun/dbf"
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/la"
)

// Operator defines the interface for differential operators such as the Laplacian and so on
type Operator interface {
	InitWithGrid(gtype string, xmin, xmax []float64, ndiv []int) (*gm.Grid, error)
	Assemble(e *la.Equations) error
	SourceTerm(e *la.Equations) error
}

// operatorMaker defines a function that makes (allocates) Operators
type operatorMaker func(params dbf.Params, source fun.Svs) (Operator, error)

// operatorDB implemetns a database of Operators
var operatorDB = make(map[string]operatorMaker)

// NewOperator finds a Operator in database or panic
func NewOperator(kind string, params dbf.Params, source fun.Svs) (Operator, error) {
	if maker, ok := operatorDB[kind]; ok {
		return maker(params, source)
	}
	return nil, chk.Err("cannot find Operator named %q in database", kind)
}
