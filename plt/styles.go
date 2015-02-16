// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"github.com/cpmech/gosl/utl"
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

// LineData holds data for ploting lines
type LineData struct {
	Label      string
	Color      string
	LineWidth  float64
	MarkerSize float64 // negative values => default
	Marker     string
	LineStyle  string
}

// GetArgs returns arguments for Plot
func (o LineData) GetArgs() string {
	l := ""
	if o.Label != "" {
		if len(l) > 0 {
			l += ","
		}
		l += utl.Sf("label='%s'", o.Label)
	}
	if o.Color != "" {
		if len(l) > 0 {
			l += ","
		}
		l += utl.Sf("color='%s'", o.Color)
	}
	if o.LineWidth > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += utl.Sf("lw=%g", o.LineWidth)
	}
	if o.MarkerSize > 0 {
		if len(l) > 0 {
			l += ","
		}
		l += utl.Sf("ms=%d", int(o.MarkerSize))
	}
	if o.Marker != "" {
		if len(l) > 0 {
			l += ","
		}
		l += utl.Sf("marker='%s'", o.Marker)
	}
	if o.LineStyle != "" {
		if len(l) > 0 {
			l += ","
		}
		l += utl.Sf("ls='%s'", o.LineStyle)
	}
	return l
}
