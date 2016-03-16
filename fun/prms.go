// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// global auxiliary variables
var (
	G_extraindent string // extra indentation
)

// Prm holds material parameter names and values
type Prm struct {

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
	Fcn   Func // a function y=f(t,x)
	Other *Prm // dependency: connected parameter

	// derived
	conn []*float64 // connected variables to V
}

// Connect connects parameter to variable
func (o *Prm) Connect(V *float64) {
	o.conn = append(o.conn, V)
	*V = o.V
}

// Set sets parameter, including connected variables
func (o *Prm) Set(V float64) {
	o.V = V
	for _, v := range o.conn {
		*v = V
	}
}

// Prms holds many parameters
type Prms []*Prm

// Find finds a parameter by name
//  Note: returns nil if not found
func (o *Prms) Find(name string) *Prm {
	for _, p := range *o {
		if p.N == name {
			return p
		}
	}
	return nil
}

// GetValues get parameter values
func (o *Prms) GetValues(names []string) (values []float64, found []bool) {
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

// Connect connects parameter
func (o *Prms) Connect(V *float64, name, caller string) (err string) {
	prm := o.Find(name)
	if prm == nil {
		return io.Sf("cannot find parameter named %q as requested by %q\n", name, caller)
	}
	prm.Connect(V)
	return
}

// ConnectSet connects set of parameters
func (o *Prms) ConnectSet(V []*float64, names []string, caller string) (err string) {
	chk.IntAssert(len(V), len(names))
	for i, name := range names {
		prm := o.Find(name)
		if prm == nil {
			return io.Sf("cannot find parameter named %q as requested by %q\n", name, caller)
		}
		prm.Connect(V[i])
	}
	return
}

func (o Prms) String() (l string) {
	for i, prm := range o {
		if i > 0 {
			l += ",\n"
		}
		l += io.Sf(G_extraindent + "{")
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
