// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import "github.com/cpmech/gosl/io"

// Formatter defines an interface to format lines with markers and colors
type Formatter interface {
	GetArgs(start string) string
}

// FmtS implemetns a formatter using a slice of strings
type FmtS []string

// ShapeData holds data for drawing shapes
type FmtH struct {
	Fc     string
	Ec     string
	Lw     int
	Closed bool
}

// LineData holds data for ploting lines
type FmtL struct {
	L  string
	C  string
	Lw float64
	Ms float64 // negative values => default
	M  string
	Ls string
}

// Init initialises ShapeData with default values
func (o *FmtH) Init() {
	o.Fc = "#edf5ff"
	o.Ec = "black"
	o.Lw = 1
}

// GetArgs returns the arguments for plt
func (o FmtS) GetArgs(start string) string {
	l := "'"
	for _, s := range o {
		l += s
	}
	return l + "'"
}

// GetArgs returns arguments for Plot
func (o FmtL) GetArgs(start string) string {
	l := start
	if o.L != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("label='%s'", o.L)
	}
	if o.C != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("color='%s'", o.C)
	}
	if o.Lw > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("lw=%g", o.Lw)
	}
	if o.Ms > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("ms=%d", int(o.Ms))
	}
	if o.M != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("marker='%s'", o.M)
	}
	if o.Ls != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("ls='%s'", o.Ls)
	}
	return l
}
