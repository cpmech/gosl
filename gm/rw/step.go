// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rw

import (
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

type Cartesian_point struct {
	Name        string
	Coordinates []float64
}

type B_spline_curve_with_knots struct {
	Name                string
	Degree              int
	Control_points_list []string
	Curve_form          string
	Closed_curve        bool
	Self_intersect      bool
	Knot_multiplicities []int
	Knots               []float64
	Knot_spec           string
}

type Surface_curve struct {
	Name                  string
	Curve_3d              string
	Associated_geometry   []string
	Master_representation string
}

type STEP struct {
	Points   map[string]*Cartesian_point
	BScurves map[string]*B_spline_curve_with_knots
	Scurves  map[string]*Surface_curve
}

func (o *STEP) ParseDATA(dat string) (err error) {

	// remove newlines and split commands
	sdat := strings.Replace(dat, "\n", "", -1)
	res := strings.Split(sdat, ";")

	// allocate maps
	o.Points = make(map[string]*Cartesian_point)
	o.BScurves = make(map[string]*B_spline_curve_with_knots)
	o.Scurves = make(map[string]*Surface_curve)

	// for each line
	reading_data := false
	for _, lin := range res {

		// activate data reading
		lin = strings.TrimSpace(strings.ToLower(lin))
		if reading_data {
			if strings.HasPrefix(lin, "endsec") {
				return
			}
		} else {
			if strings.HasPrefix(lin, "data") {
				reading_data = true
			}
			continue
		}

		// left- and right-hand-sides => key = function
		lhs_rhs := strings.Split(lin, "=")
		if len(lhs_rhs) != 2 {
			continue
		}

		// key
		lhs := strings.TrimSpace(lhs_rhs[0])

		// function call
		rhs := strings.ToLower(strings.TrimSpace(lhs_rhs[1]))

		// extract entities
		switch {

		// Cartesian point
		case strings.HasPrefix(rhs, "cartesian_point"):
			sargs := strings.TrimPrefix(rhs, "cartesian_point")
			args := io.SplitWithinParentheses(sargs)
			n := len(args)
			if n != 2 {
				err = chk.Err("cartesian_point has the wrong number of arguments. %d != %d", n, 2)
				return
			}
			p := Cartesian_point{
				Name:        args[0],
				Coordinates: io.SplitFloats(args[1]),
			}
			o.Points[lhs] = &p

		// B-spline curve
		case strings.HasPrefix(rhs, "b_spline_curve_with_knots"):
			sargs := strings.TrimPrefix(rhs, "b_spline_curve_with_knots")
			args := io.SplitWithinParentheses(sargs)
			n := len(args)
			if n != 9 {
				err = chk.Err("b_spline_curve_with_knots has the wrong number of arguments. %d != %d", n, 9)
				return
			}
			b := B_spline_curve_with_knots{
				Name:                args[0],
				Degree:              io.Atoi(args[1]),
				Control_points_list: strings.Fields(args[2]),
				Curve_form:          args[3],
				Closed_curve:        atob(args[4]),
				Self_intersect:      atob(args[5]),
				Knot_multiplicities: io.SplitInts(args[6]),
				Knots:               io.SplitFloats(args[7]),
				Knot_spec:           args[8],
			}
			chk.IntAssert(len(b.Knot_multiplicities), len(b.Knots))
			o.BScurves[lhs] = &b

		// surface curve
		case strings.HasPrefix(rhs, "surface_curve"):
			sargs := strings.TrimPrefix(rhs, "surface_curve")
			args := io.SplitWithinParentheses(sargs)
			n := len(args)
			if n != 4 {
				err = chk.Err("surface_curve has the wrong number of arguments. %d != %d", n, 4)
				return
			}
			b := Surface_curve{
				Name:                  args[0],
				Curve_3d:              args[1],
				Associated_geometry:   io.SplitWithinParentheses(args[2]),
				Master_representation: args[3],
			}
			o.Scurves[lhs] = &b
		}
	}
	return
}
