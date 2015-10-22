// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import "github.com/cpmech/gosl/io"

// Sty holds data for drawing shapes
type Sty struct {
	Fc     string  // face color
	Ec     string  // edge color
	Lw     float64 // linewidth
	Closed bool    // closed shape
}

// Fmt holds data for ploting lines
type Fmt struct {
	C  string  // color
	M  string  // marker
	Ls string  // linestyle
	Lw float64 // linewidth; -1 => default
	Ms int     // marker size; -1 => default
	L  string  // label
	Me int     // mark-every; -1 => default
}

// Init initialises Fmt with default values
func (o *Sty) Init() {
	o.Fc = "#edf5ff"
	o.Ec = "black"
	o.Lw = 1
}

// GetArgs returns arguments for Plot
func (o Fmt) GetArgs(start string) string {
	l := start
	if o.C != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("color='%s'", o.C)
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
		l += io.Sf("ms=%d", o.Ms)
	}
	if o.L != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("label='%s'", o.L)
	}
	if o.Me > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("markevery=%d", o.Me)
	}
	return l
}
