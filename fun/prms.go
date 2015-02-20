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
	N      string  `json:"n"`      // name of parameter
	V      float64 `json:"v"`      // value of parameter
	U      string  `json:"u"`      // unit (not verified)
	Extra  string  `json:"extra"`  // extra data
	Inact  bool    `json:"inact"`  // parameter is inactive in optimisation
	SetDef bool    `json:"setdef"` // tells model to use a default value
	Fcn    Func    `json:"fcn"`    // a function y=f(t,x)
}

// Prms holds many parameters
type Prms []*Prm

// Find finds a parameter by name
//  Note: returns nil if not found
func (o Prms) Find(name string) *Prm {
	for _, p := range o {
		if p.N == name {
			return p
		}
	}
	return nil
}

// String outputs a nice formatted representation of a parameter
func (o *Prm) String() string {
	sknam, skval := "%s", "%s"
	ncom := 0
	if o.U != "" {
		ncom = 1
	}
	if g_largestname > 0 {
		sknam = io.Sf("%%-%ds", g_largestname+3)
	}
	if g_largestsval > 0 {
		skval = io.Sf("%%-%ds", g_largestsval+ncom)
	}
	l := ""
	if o.U != "" {
		l = io.Sf("{\"n\":"+sknam+" \"v\":"+skval+" \"u\":%q", "\""+o.N+"\",", io.Sf("%g", o.V)+",", o.U)
	} else {
		l = io.Sf("{\"n\":"+sknam+" \"v\":"+skval, "\""+o.N+"\",", io.Sf("%g", o.V))
	}
	if o.Extra != "" {
		l += io.Sf(", \"extra\":%q", o.Extra)
	}
	if o.Inact {
		l += ", \"inact\":true"
	}
	if o.SetDef {
		l += ", \"setdef\":true"
	}
	l += "}"
	return l
}

// String outputs all materials
func (o Prms) String() string {
	g_largestname, g_largestsval = 0, 0
	for _, prm := range o {
		g_largestname = imax(g_largestname, len(prm.N))
		g_largestsval = imax(g_largestsval, len(io.Sf("%g", prm.V)))
	}
	l := ""
	if G_openbrackets {
		l += io.Sf("%s    [\n", G_extraindent)
	}
	for j, prm := range o {
		if j > 0 {
			l += ",\n"
		}
		l += io.Sf("%s      %v", G_extraindent, prm)
	}
	if len(o) > 0 {
		l += io.Sf("\n")
	}
	l += io.Sf("%s    ]", G_extraindent)
	return l
}
