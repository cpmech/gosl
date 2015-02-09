// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "github.com/cpmech/gosl/utl"

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

// GetValueOrPanic returns the parameter value if it exists; if it does not exists then panic.
// minval and maxval are used to check parameter range if usemin or usemax are true
func (o Prms) GetValueOrPanic(name string, minval, maxval float64, usemin, usemax bool) float64 {
	prm := o.Find(name)
	if prm == nil {
		utl.Panic("cannot get parameter named %s", name)
	}
	if usemin {
		if prm.V < minval {
			utl.Panic("%s parameter: wrong value: %g < %g", name, prm.V, minval)
		}
	}
	if usemax {
		if prm.V > maxval {
			utl.Panic("%s parameter: wrong value: %g > %g", name, prm.V, maxval)
		}
	}
	return prm.V
}

// String outputs a nice formatted representation of a parameter
func (o *Prm) String() string {
	sknam, skval := "%s", "%s"
	ncom := 0
	if o.U != "" {
		ncom = 1
	}
	if g_largestname > 0 {
		sknam = utl.Sf("%%-%ds", g_largestname+3)
	}
	if g_largestsval > 0 {
		skval = utl.Sf("%%-%ds", g_largestsval+ncom)
	}
	l := ""
	if o.U != "" {
		l = utl.Sf("{\"n\":"+sknam+" \"v\":"+skval+" \"u\":%q", "\""+o.N+"\",", utl.Sf("%g", o.V)+",", o.U)
	} else {
		l = utl.Sf("{\"n\":"+sknam+" \"v\":"+skval, "\""+o.N+"\",", utl.Sf("%g", o.V))
	}
	if o.Extra != "" {
		l += utl.Sf(", \"extra\":%q", o.Extra)
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
		g_largestsval = imax(g_largestsval, len(utl.Sf("%g", prm.V)))
	}
	l := ""
	if G_openbrackets {
		l += utl.Sf("%s    [\n", G_extraindent)
	}
	for j, prm := range o {
		if j > 0 {
			l += ",\n"
		}
		l += utl.Sf("%s      %v", G_extraindent, prm)
	}
	l += utl.Sf("\n%s    ]", G_extraindent)
	return l
}
