// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Reference:
// [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
//     and Mathematical Tables. U.S. Department of Commerce, NIST

func TestBessel01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Bessel01")

	// load data
	_, dat, err := io.ReadTable("data/as-9-bessel-integer-sml.cmp")
	//_, dat, err := io.ReadTable("data/as-9-bessel-integer-big.cmp")
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check J0
	for i, x := range dat["x"] {
		j0 := math.J0(x)
		chk.Scalar(tst, io.Sf("J0(%11.8f)=%23.15e", x, j0), 1e-15, j0, dat["J0"][i])
	}

	// check J1
	io.Pl()
	for i, x := range dat["x"] {
		j1 := math.J1(x)
		chk.Scalar(tst, io.Sf("J1(%11.8f)=%23.15e", x, j1), 1e-15, j1, dat["J1"][i])
	}

	// check J2
	io.Pl()
	for i, x := range dat["x"] {
		j2 := math.Jn(2, x)
		chk.Scalar(tst, io.Sf("J2(%11.8f)=%23.15e", x, j2), 1e-15, j2, dat["J2"][i])
	}

	// check Y0
	io.Pl()
	for i, x := range dat["x"] {
		y0 := math.Y0(x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(y0, -1) {
				tst.Errorf("Y0(0) should be -Inf. %v is incorrect\n", y0)
				return
			}
			if !math.IsInf(dat["Y0"][i], -1) {
				tst.Errorf("Reference data is incorrect: Y0(0)=%v should be -Inf\n", dat["Y0"][i])
				return
			}
		} else {
			chk.Scalar(tst, io.Sf("Y0(%11.8f)=%23.15e", x, y0), 1e-15, y0, dat["Y0"][i])
		}
	}

	// check Y1
	io.Pl()
	for i, x := range dat["x"] {
		y0 := math.Y1(x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(y0, -1) {
				tst.Errorf("Y1(0) should be -Inf. %v is incorrect\n", y0)
				return
			}
			if !math.IsInf(dat["Y1"][i], -1) {
				tst.Errorf("Reference data is incorrect: Y1(0)=%v should be -Inf\n", dat["Y1"][i])
				return
			}
		} else {
			chk.Scalar(tst, io.Sf("Y1(%11.8f)=%23.15e", x, y0), 1e-15, y0, dat["Y1"][i])
		}
	}

	// check Y2
	io.Pl()
	for i, x := range dat["x"] {
		y0 := math.Yn(2, x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(y0, -1) {
				tst.Errorf("Y2(0) should be -Inf. %v is incorrect\n", y0)
				return
			}
			if !math.IsInf(dat["Y2"][i], -1) {
				tst.Errorf("Reference data is incorrect: Y2(0)=%v should be -Inf\n", dat["Y2"][i])
				return
			}
		} else {
			chk.Scalar(tst, io.Sf("Y2(%11.8f)=%23.15e", x, y0), 1e-15, y0, dat["Y2"][i])
		}
	}
}
