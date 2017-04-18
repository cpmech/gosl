// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import "github.com/cpmech/gosl/io"

// S holds "style" data for drawing polygons or other shapes
type S struct {
	Fc     string  // face color
	Ec     string  // edge color
	Lw     float64 // linewidth
	Closed bool    // closed shape
}

// A holds "arguments" to configure plots
type A struct {
	C     string  // color
	M     string  // marker
	Ls    string  // linestyle
	Lw    float64 // linewidth; -1 => default
	Ms    int     // marker size; -1 => default
	L     string  // label
	Me    int     // mark-every; -1 => default
	Z     int     // z-order
	Mec   string  // marker edge color
	Mew   float64 // marker edge width
	Void  bool    // void marker => markeredgecolor='C', markerfacecolor='none'
	Clip  bool    // turn clip => clip_on=True
	Extra string  // extra arguments
}

// Init initialises S (style) with default values
func (o *S) Init() {
	o.Fc = "#edf5ff"
	o.Ec = "black"
	o.Lw = 1
}

// String returns a string representation of arguments
func (o A) String(start string) string {
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
	if o.Z > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("zorder=%d", o.Z)
	}
	if o.Mec != "" {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("markeredgecolor='%s'", o.Mec)
	}
	if o.Mew > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += io.Sf("mew=%g", o.Mew)
	}
	if o.Void {
		if len(l) > 0 {
			l += ","
		}
		l += "markerfacecolor='none'"
		if o.Mec == "" {
			l += io.Sf(",markeredgecolor='%s'", o.C)
		}
	}
	if o.Clip {
		if len(l) > 0 {
			l += ","
		}
		l += "clip_on=1"
	} else {
		if len(l) > 0 {
			l += ","
		}
		l += "clip_on=0"
	}
	if o.Extra != "" {
		l += o.Extra
	}
	return l
}
