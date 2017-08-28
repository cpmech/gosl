// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// P holds material parameter names and values
//
// The connected variables to V data holds pointers to other scalars that need to be updated when
// the paramter is changed. For instance, when running simulations with variable parameters.
//
type P struct {

	// input
	N      string  `json:"n"`      // name of parameter
	V      float64 `json:"v"`      // value of parameter
	Min    float64 `json:"min"`    // min value
	Max    float64 `json:"max"`    // max value
	S      float64 `json:"s"`      // standard deviation
	D      string  `json:"d"`      // probability distribution type
	U      string  `json:"u"`      // unit (not verified)
	Adj    int     `json:"adj"`    // adjustable: unique ID (greater than zero)
	Dep    int     `json:"dep"`    // depends on "adj"
	Extra  string  `json:"extra"`  // extra data
	Inact  bool    `json:"inact"`  // parameter is inactive in optimisation
	SetDef bool    `json:"setdef"` // tells model to use a default value

	// auxiliary
	Fcn   T  // a function y=f(t,x)
	Other *P // dependency: connected parameter

	// derived
	conn []*float64 // connected variables to V
}

// Connect connects parameter to variable
func (o *P) Connect(V *float64) {
	o.conn = append(o.conn, V)
	*V = o.V
}

// Set sets parameter, including connected variables
func (o *P) Set(V float64) {
	o.V = V
	for _, v := range o.conn {
		*v = V
	}
}

// Params holds many parameters
type Params []*P

// Find finds a parameter by name
//  Note: returns nil if not found
func (o *Params) Find(name string) *P {
	for _, p := range *o {
		if p.N == name {
			return p
		}
	}
	return nil
}

// CheckLimits check limits of variables given in Min/Max
// Will panic if values are outside corresponding Min/Max range.
func (o *Params) CheckLimits() {
	for _, p := range *o {
		if p.V < p.Min {
			chk.Panic("parameter %q has value smaller than minimum. %v < %v is not aceptable", p.N, p.V, p.Min)
		}
		if p.V > p.Max {
			chk.Panic("parameter %q has value greater than maximum. %v > %v is not aceptable", p.N, p.V, p.Max)
		}
	}
}

// GetValues get parameter values
func (o *Params) GetValues(names []string) (values []float64, found []bool) {
	n := len(names)
	values = make([]float64, n)
	found = make([]bool, n)
	for i, name := range names {
		prm := o.Find(name)
		if prm != nil {
			values[i] = prm.V
			found[i] = true
		}
	}
	return
}

// CheckAndGetValues check min/max limits and return values.
// Will panic if values are outside corresponding min/max range.
// Will also panic if a parameter name is not found.
func (o *Params) CheckAndGetValues(names []string) (values []float64) {
	n := len(names)
	values = make([]float64, n)
	for i, name := range names {
		prm := o.Find(name)
		if prm == nil {
			chk.Panic("cannot find parameter named %q", name)
		}
		if prm.V < prm.Min {
			chk.Panic("parameter %q has value smaller than minimum. %v < %v is not aceptable", name, prm.V, prm.Min)
		}
		if prm.V > prm.Max {
			chk.Panic("parameter %q has value greater than maximum. %v > %v is not aceptable", name, prm.V, prm.Max)
		}
		values[i] = prm.V
	}
	return
}

// CheckAndSetVars get parameter values and check limits defined in Min and Max
// Will panic if values are outside corresponding Min/Max range.
// Will also panic if a parameter name is not found.
func (o *Params) CheckAndSetVars(names []string, variables []*float64) {
	n := len(names)
	if len(variables) != n {
		chk.Panic("array of variables must have the same size as the slice of names. %d != %d", len(variables), n)
	}
	for i, name := range names {
		prm := o.Find(name)
		if prm == nil {
			chk.Panic("cannot find parameter named %q", name)
		}
		if prm.V < prm.Min {
			chk.Panic("parameter %q has value smaller than minimum. %v < %v is not aceptable", name, prm.V, prm.Min)
		}
		if prm.V > prm.Max {
			chk.Panic("parameter %q has value greater than maximum. %v > %v is not aceptable", name, prm.V, prm.Max)
		}
		if variables[i] == nil {
			chk.Panic("array of variables must not have nil entries")
		}
		*variables[i] = prm.V
	}
	return
}

// Connect connects parameter
func (o *Params) Connect(V *float64, name, caller string) (errorMessage string) {
	prm := o.Find(name)
	if prm == nil {
		return io.Sf("cannot find parameter named %q as requested by %q\n", name, caller)
	}
	prm.Connect(V)
	return
}

// ConnectSet connects set of parameters
func (o *Params) ConnectSet(V []*float64, names []string, caller string) (errorMessage string) {
	chk.IntAssert(len(V), len(names))
	for i, name := range names {
		prm := o.Find(name)
		if prm == nil {
			errorMessage += io.Sf("cannot find parameter named %q as requested by %q\n", name, caller)
		} else {
			prm.Connect(V[i])
		}
	}
	return
}

func (o Params) String() (l string) {
	for i, prm := range o {
		if i > 0 {
			l += ",\n"
		}
		l += io.Sf("{")
		l += io.Sf(`"n":%q, `, prm.N)
		l += io.Sf(`"v":%v, `, prm.V)
		l += io.Sf(`"min":%v, `, prm.Min)
		l += io.Sf(`"max":%v, `, prm.Max)
		l += io.Sf(`"s":%v, `, prm.S)
		l += io.Sf(`"d":%q, `, prm.D)
		l += io.Sf(`"u":%q, `, prm.U)
		l += io.Sf(`"adj":%v, `, prm.Adj)
		l += io.Sf(`"dep":%v, `, prm.Dep)
		l += io.Sf(`"extra":%q, `, prm.Extra)
		l += io.Sf(`"inact":%v, `, prm.Inact)
		l += io.Sf(`"setdef":%v`, prm.SetDef)
		l += io.Sf("}")
	}
	return
}
