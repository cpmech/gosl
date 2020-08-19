// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rw

import (
	"strings"

	"gosl/chk"
	"gosl/io"
)

// CartesianPoint holds Cartesian point data
type CartesianPoint struct {
	Name        string
	Coordinates []float64
}

// BsplineCurveWithKnots holds data representing B-splines
type BsplineCurveWithKnots struct {
	Name               string
	Degree             int
	ControlPointsList  []string
	CurveForm          string
	ClosedCurve        bool
	SelfIntersect      bool
	KnotMultiplicities []int
	Knots              []float64
	KnotSpec           string
}

// SurfaceCurve represents a surface curve
type SurfaceCurve struct {
	Name                 string
	Curve3d              string
	AssociatedGeometry   []string
	MasterRepresentation string
}

// StepFile holds information to read/write STEP files
type StepFile struct {
	Points        map[string]*CartesianPoint
	BsplineCurves map[string]*BsplineCurveWithKnots
	SurfaceCurves map[string]*SurfaceCurve
}

// ParseDataFile parses data section from file
func (o *StepFile) ParseDataFile(filename string) {
	buf := io.ReadFile(filename)
	o.ParseData(string(buf))
}

// ParseData parses data section of Step file
func (o *StepFile) ParseData(dat string) {

	// remove newlines and split commands
	sdat := strings.Replace(dat, "\n", "", -1)
	res := strings.Split(sdat, ";")

	// allocate maps
	o.Points = make(map[string]*CartesianPoint)
	o.BsplineCurves = make(map[string]*BsplineCurveWithKnots)
	o.SurfaceCurves = make(map[string]*SurfaceCurve)

	// for each line
	readingData := false
	for _, lin := range res {

		// activate data reading
		lin = strings.TrimSpace(strings.ToLower(lin))
		if readingData {
			if strings.HasPrefix(lin, "endsec") {
				return
			}
		} else {
			if strings.HasPrefix(lin, "data") {
				readingData = true
			}
			continue
		}

		// left- and right-hand-sides => key = function
		lhsAndRHS := strings.Split(lin, "=")
		if len(lhsAndRHS) != 2 {
			continue
		}

		// key
		lhs := strings.TrimSpace(lhsAndRHS[0])

		// function call
		rhs := strings.ToLower(strings.TrimSpace(lhsAndRHS[1]))

		// extract entities
		switch {

		// Cartesian point
		case strings.HasPrefix(rhs, "cartesian_point"):
			sargs := strings.TrimPrefix(rhs, "cartesian_point")
			args := io.SplitWithinParentheses(sargs)
			n := len(args)
			var name, strFloats string
			switch n {
			case 1:
				strFloats = args[0]
			case 2:
				name = args[0]
				strFloats = args[1]
			default:
				chk.Panic("cartesian_point has the wrong number of arguments. n=%d is invalid\n", n)
			}
			p := CartesianPoint{
				Name:        name,
				Coordinates: io.SplitFloats(strFloats),
			}
			o.Points[lhs] = &p

		// B-spline curve
		case strings.HasPrefix(rhs, "b_spline_curve_with_knots"):
			sargs := strings.TrimPrefix(rhs, "b_spline_curve_with_knots")
			args := io.SplitWithinParentheses(sargs)
			n := len(args)
			if n != 9 {
				chk.Panic("b_spline_curve_with_knots has the wrong number of arguments. %d != %d\n", n, 9)
			}
			b := BsplineCurveWithKnots{
				Name:               args[0],
				Degree:             io.Atoi(args[1]),
				ControlPointsList:  strings.Fields(args[2]),
				CurveForm:          args[3],
				ClosedCurve:        atob(args[4]),
				SelfIntersect:      atob(args[5]),
				KnotMultiplicities: io.SplitInts(args[6]),
				Knots:              io.SplitFloats(args[7]),
				KnotSpec:           args[8],
			}
			chk.IntAssert(len(b.KnotMultiplicities), len(b.Knots))
			o.BsplineCurves[lhs] = &b

		// surface curve
		case strings.HasPrefix(rhs, "surface_curve"):
			sargs := strings.TrimPrefix(rhs, "surface_curve")
			args := io.SplitWithinParentheses(sargs)
			n := len(args)
			if n != 4 {
				chk.Panic("surface_curve has the wrong number of arguments. %d != %d\n", n, 4)
			}
			b := SurfaceCurve{
				Name:                 args[0],
				Curve3d:              args[1],
				AssociatedGeometry:   io.SplitWithinParentheses(args[2]),
				MasterRepresentation: args[3],
			}
			o.SurfaceCurves[lhs] = &b
		}
	}
}
