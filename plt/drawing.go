// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import "github.com/cpmech/gosl/io"

// AutoScale rescales plot area
func AutoScale(P [][]float64) {
	if len(P) < 1 {
		return
	}
	xmin, ymin := P[0][0], P[0][1]
	xmax, ymax := xmin, ymin
	for _, p := range P {
		if p[0] < xmin {
			xmin = p[0]
		}
		if p[1] < ymin {
			ymin = p[1]
		}
		if p[0] > xmax {
			xmax = p[0]
		}
		if p[1] > ymax {
			ymax = p[1]
		}
	}
	io.Ff(&bb, "axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// DrawPolyline draws a polyline
func DrawPolyline(P [][]float64, sd *FmtH, args string) {
	if len(P) < 1 {
		return
	}
	n := bb.Len()
	io.Ff(&bb, "dat%d = [[MPLPath.MOVETO, [%g, %g]]", n, P[0][0], P[0][1])
	for _, p := range P {
		io.Ff(&bb, ", [MPLPath.LINETO, [%g, %g]]", p[0], p[1])
	}
	if sd.Closed {
		io.Ff(&bb, ", [MPLPath.CLOSEPOLY, [0, 0]]")
	}
	io.Ff(&bb, "]\n")
	io.Ff(&bb, "commands%d, vertices%d = zip(*dat%d)\n", n, n, n)
	io.Ff(&bb, "ph%d = MPLPath(vertices%d, commands%d)\n", n, n, n)
	io.Ff(&bb, "pc%d = PathPatch(ph%d", n, n)
	io.Ff(&bb, ", fc='%s', ec='%s', lw=%d", sd.Fc, sd.Ec, sd.Lw)
	if len(args) > 0 {
		io.Ff(&bb, ", %s)\n", args)
	} else {
		io.Ff(&bb, ")\n")
	}
	io.Ff(&bb, "gca().add_patch(pc%d)\n", n)
}

// DrawLegend draws legend with given lines data. fs == fontsize
func DrawLegend(dat []FmtL, fs int, loc string, frame bool, args string) {
	n := bb.Len()
	io.Ff(&bb, "handles%d = [", n)
	for i, d := range dat {
		if i > 0 {
			io.Ff(&bb, ",\n")
		}
		io.Ff(&bb, "Line2D([], [], label='%s', color='%s', lw=%g, marker='%s', ls='%s'", d.L, d.C, d.Lw, d.M, d.Ls)
		if d.Ms >= 0 {
			io.Ff(&bb, ", ms=%g", d.Ms)
		}
		io.Ff(&bb, ")")
	}
	io.Ff(&bb, "]\nlg%d=legend(handles=handles%d, fontsize=%d, loc='%s'", n, n, fs, loc)
	if len(args) > 0 {
		io.Ff(&bb, ", %s", args)
	}
	io.Ff(&bb, ")\n")
	if !frame {
		io.Ff(&bb, "lg%d.get_frame().set_linewidth(0.0)\n", n)
	}
	io.Ff(&bb, "ea.append(lg%d)\n", n)
}
