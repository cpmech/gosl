// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// P holds numeric parameters defined by a name N and a value V.
//
// P is convenient to store the range of allowed values in Min and Max,
// and other information such as standard deviation S, probability distribution type D,
// among others.
//
// Dependent variables may be connected to P using Connect so when Set is called,
// the dependendt variable is updated as well.
//
// Other parameters can be linked to this one via the Other data member
// and Func may be useful to compute y=f(t,x)
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
	Inact  bool    `json:"inact"`  // parameter is inactive in optimization
	SetDef bool    `json:"setdef"` // tells model to use a default value

	// connected parameter
	Other *P // dependency: connected parameter

	// function
	Func func(t float64, x []float64) // a function y=f(t,x)

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
//
//   A set of Params can be initialized as follows:
//
//     var params Params
//     params = []*P{
//         {N: "klx", V: 1.0},
//         {N: "kly", V: 2.0},
//         {N: "klz", V: 3.0},
//     }
//
//   Alternatively, see NewParams function
//
type Params []*P

// NewParams returns a set of parameters
//
//   This is an alternative to initializing Params by setting slice items
//
//   A set of Params can be initialized as follows:
//
//     params := NewParams(
//         &P{N: "P1", V: 1},
//         &P{N: "P2", V: 2},
//         &P{N: "P3", V: 3},
//     )
//
//   Alternatively, you may set slice components directly (see Params definition)
//
func NewParams(pp ...interface{}) (o Params) {
	o = make([]*P, len(pp))
	for i, p := range pp {
		o[i] = p.(*P)
	}
	return
}

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

// GetValue reads parameter or Panic
// Will panic if name does not exist in parameters set
func (o *Params) GetValue(name string) float64 {
	p := o.Find(name)
	if p == nil {
		chk.Panic("cannot find parameter named %q\n", name)
	}
	return p.V
}

// GetValueOrDefault reads parameter or returns default value
// Will return defaultValue if name does not exist in parameters set
func (o *Params) GetValueOrDefault(name string, defaultValue float64) float64 {
	p := o.Find(name)
	if p == nil {
		return defaultValue
	}
	return p.V
}

// GetIntOrDefault reads parameter or returns default value
// Will return defaultInt if name does not exist in parameters set
func (o *Params) GetIntOrDefault(name string, defaultInt int) int {
	return int(o.GetValueOrDefault(name, float64(defaultInt)))
}

// GetBool reads Boolean parameter or Panic
// Returns true if P[name] > 0; otherwise returns false
// Will panic if name does not exist in parameters set
func (o *Params) GetBool(name string) bool {
	p := o.Find(name)
	if p == nil {
		chk.Panic("cannot find Boolean parameter named %q\n", name)
	}
	if p.V > 0 {
		return true
	}
	return false
}

// GetBoolOrDefault reads Boolean parameter or returns default value
// Returns true if P[name] > 0; otherwise returns false
// Will return defaultValue if name does not exist in parameters set
func (o *Params) GetBoolOrDefault(name string, defaultValue bool) bool {
	p := o.Find(name)
	if p == nil {
		return defaultValue
	}
	if p.V > 0 {
		return true
	}
	return false
}

// SetValue sets parameter or Panic
// Will panic if name does not exist in parameters set
func (o *Params) SetValue(name string, value float64) {
	p := o.Find(name)
	if p == nil {
		chk.Panic("cannot find parameter named %q\n", name)
	}
	p.V = value
}

// SetBool sets Boolean parameter or Panic
// Sets +1==true if value > 0; otherwise sets -1==false
// Will panic if name does not exist in parameters set
func (o *Params) SetBool(name string, value float64) {
	p := o.Find(name)
	if p == nil {
		chk.Panic("cannot find Boolean parameter named %q\n", name)
	}
	if value > 0 {
		p.V = +1.0
		return
	}
	p.V = -1.0
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

// CheckAndSetVariables get parameter values and check limits defined in Min and Max
// Will panic if values are outside corresponding Min/Max range.
// Will also panic if a parameter name is not found.
func (o *Params) CheckAndSetVariables(names []string, variables []*float64) {
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
		io.Pforan("name=%v  prm = %v\n", name, prm)
		if prm == nil {
			errorMessage += io.Sf("cannot find parameter named %q as requested by %q\n", name, caller)
		} else {
			prm.Connect(V[i])
		}
	}
	return
}

// ConnectSetOpt connects set of parameters with some being optional
func (o *Params) ConnectSetOpt(V []*float64, names []string, optional []bool, caller string) (errorMessage string) {
	chk.IntAssert(len(V), len(names))
	chk.IntAssert(len(V), len(optional))
	for i, name := range names {
		prm := o.Find(name)
		if prm == nil {
			if !optional[i] {
				errorMessage += io.Sf("cannot find parameter named %q as requested by %q\n", name, caller)
			}
		} else {
			prm.Connect(V[i])
		}
	}
	return
}

// String returns a summary of parameters
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
