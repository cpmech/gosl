// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
    "code.google.com/p/gosl/utl"
)

// ShapeData holds data for drawing shapes
type ShapeData struct {
    FaceColor string
    EdgeColor string
    LineWidth int
    Closed    bool
}

// Init initialises ShapeData with default values
func (o *ShapeData) Init() {
    o.FaceColor = "#edf5ff"
    o.EdgeColor = "black"
    o.LineWidth = 1
}

// AutoScale rescales plot area
func AutoScale(P [][]float64) {
    if len(P) < 1 {
        return
    }
    xmin, ymin := P[0][0], P[0][1]
    xmax, ymax := xmin, ymin
    for _, p := range P {
        if p[0] < xmin { xmin = p[0] }
        if p[1] < ymin { ymin = p[1] }
        if p[0] > xmax { xmax = p[0] }
        if p[1] > ymax { ymax = p[1] }
    }
    utl.Ff(&bb, "axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// DrawPolyline draws a polyline
func DrawPolyline(P [][]float64, sd *ShapeData, args string) {
    if len(P) < 1 {
        return
    }
    n := bb.Len()
    utl.Ff(&bb, "dat%d = [[MPLPath.MOVETO, [%g, %g]]", n, P[0][0], P[0][1])
    for _, p := range P {
        utl.Ff(&bb, ", [MPLPath.LINETO, [%g, %g]]", p[0], p[1])
    }
    if sd.Closed {
        utl.Ff(&bb, ", [MPLPath.CLOSEPOLY, [0, 0]]")
    }
    utl.Ff(&bb, "]\n")
    utl.Ff(&bb, "commands%d, vertices%d = zip(*dat%d)\n", n, n, n)
    utl.Ff(&bb, "ph%d = MPLPath(vertices%d, commands%d)\n", n, n, n)
    utl.Ff(&bb, "pc%d = PathPatch(ph%d", n, n)
    utl.Ff(&bb, ", fc='%s', ec='%s', lw=%d", sd.FaceColor, sd.EdgeColor, sd.LineWidth)
    if len(args) > 0 {
        utl.Ff(&bb, ", %s)\n", args)
    } else {
        utl.Ff(&bb, ")\n")
    }
    utl.Ff(&bb, "gca().add_patch(pc%d)\n", n)
}

// LineData holds data for ploting lines
type LineData struct {
    Label      string
    Color      string
    LineWidth  float64
    MarkerSize float64 // negative values => default
    Marker     string
    LineStyle  string
}

// DrawLegend draws legend with given lines data. fs == fontsize
func DrawLegend(dat []LineData, fs int, loc string, frame bool, args string) {
    n := bb.Len()
    utl.Ff(&bb, "handles%d = [", n)
    for i, d := range dat {
        if i > 0 {
            utl.Ff(&bb, ",\n")
        }
        utl.Ff(&bb, "Line2D([], [], label='%s', color='%s', lw=%g, marker='%s', ls='%s'", d.Label, d.Color, d.LineWidth, d.Marker, d.LineStyle)
        if d.MarkerSize >= 0 {
            utl.Ff(&bb, ", ms=%g", d.MarkerSize)
        }
        utl.Ff(&bb, ")")
    }
    utl.Ff(&bb, "]\nlg%d=legend(handles=handles%d, fontsize=%d, loc='%s'", n, n, fs, loc)
    if len(args) > 0 {
        utl.Ff(&bb, ", %s", args)
    }
    utl.Ff(&bb, ")\n")
    if !frame {
        utl.Ff(&bb, "lg%d.get_frame().set_linewidth(0.0)\n", n)
    }
    utl.Ff(&bb, "ea.append(lg%d)\n", n)
}
