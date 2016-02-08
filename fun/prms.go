// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/io"

// global auxiliary variables
var (
	g_largestname  int    // largest length of paramter name (to make a nice table)
	g_largestsval  int    // largest length of paramter value string representation (to make a nice table)
	G_extraindent  string // extra indentation
	G_openbrackets bool   // add initial brackets
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
	Adj    string  `json:"adj"`    // adjustable: search key
	Dep    string  `json:"dep"`    // depends on
	Extra  string  `json:"extra"`  // extra data
	Inact  bool    `json:"inact"`  // parameter is inactive in optimisation
	SetDef bool    `json:"setdef"` // tells model to use a default value
	Fcn    Func    `json:"fcn"`    // a function y=f(t,x)

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

// Connect connects parameter
func (o *Prms) Connect(V *float64, name string) (err string) {
	prm := o.Find(name)
	if prm == nil {
		return io.Sf("cannot find parameter named %q\n", name)
	}
	prm.Connect(V)
	return
}

func (o Prms) String() (l string) {
	for _, prm := range o {
		l += io.Sf("\nN=%q, ", prm.N)
		l += io.Sf("V=%v, ", prm.V)
		l += io.Sf("Min=%v, ", prm.Min)
		l += io.Sf("Max=%v, ", prm.Max)
		l += io.Sf("S=%v, ", prm.S)
		l += io.Sf("D=%q\n", prm.D)
		l += io.Sf("U=%v, ", prm.U)
		l += io.Sf("Adj=%q, ", prm.Adj)
		l += io.Sf("Dep=%q, ", prm.Dep)
		l += io.Sf("Extra=%q, ", prm.Extra)
		l += io.Sf("Inact=%v, ", prm.Inact)
		l += io.Sf("SetDef=%v, ", prm.SetDef)
		l += io.Sf("Fcn=%v\n", prm.Fcn)
	}
	return
}
