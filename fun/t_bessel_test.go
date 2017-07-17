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
	chk.PrintTitle("Bessel01. Standard Bessel functions")

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
		chk.Float64(tst, io.Sf("J0(%11.8f)=%23.15e", x, j0), 1e-15, j0, dat["J0"][i])
	}

	// check J1
	io.Pl()
	for i, x := range dat["x"] {
		j1 := math.J1(x)
		chk.Float64(tst, io.Sf("J1(%11.8f)=%23.15e", x, j1), 1e-15, j1, dat["J1"][i])
	}

	// check J2
	io.Pl()
	for i, x := range dat["x"] {
		j2 := math.Jn(2, x)
		chk.Float64(tst, io.Sf("J2(%11.8f)=%23.15e", x, j2), 1e-15, j2, dat["J2"][i])
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
			chk.Float64(tst, io.Sf("Y0(%11.8f)=%23.15e", x, y0), 1e-15, y0, dat["Y0"][i])
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
			chk.Float64(tst, io.Sf("Y1(%11.8f)=%23.15e", x, y0), 1e-15, y0, dat["Y1"][i])
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
			chk.Float64(tst, io.Sf("Y2(%11.8f)=%23.15e", x, y0), 1e-15, y0, dat["Y2"][i])
		}
	}
}

func TestBessel02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Bessel02. Modified Bessel functions")

	// load data
	_, dneg, err := io.ReadTable("data/as-9-modbessel-integer-neg.cmp")
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	_, dat, err := io.ReadTable("data/as-9-modbessel-integer-sml.cmp")
	//_, dat, err := io.ReadTable("data/as-9-modbessel-integer-big.cmp")
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// check I0 (negative values)
	for i, x := range dneg["x"] {
		I0 := ModBesselI0(x)
		chk.Float64(tst, io.Sf("I0(%11.8f)=%23.15e", x, I0), 1e-8, I0, dneg["I0"][i])
	}

	// check I1 (negative values)
	io.Pl()
	for i, x := range dneg["x"] {
		I1 := ModBesselI1(x)
		chk.Float64(tst, io.Sf("I1(%11.8f)=%23.15e", x, I1), 1e-7, I1, dneg["I1"][i])
	}

	// check I2 (negative values)
	io.Pl()
	for i, x := range dneg["x"] {
		I2 := ModBesselIn(2, x)
		chk.Float64(tst, io.Sf("I2(%11.8f)=%23.15e", x, I2), 1e-7, I2, dneg["I2"][i])
	}

	// check I3 (negative values)
	io.Pl()
	for i, x := range dneg["x"] {
		I3 := ModBesselIn(3, x)
		chk.Float64(tst, io.Sf("I3(%11.8f)=%23.15e", x, I3), 1e-7, I3, dneg["I3"][i])
	}

	// check I0
	io.Pl()
	for i, x := range dat["x"] {
		I0 := ModBesselI0(x)
		chk.Float64(tst, io.Sf("I0(%11.8f)=%23.15e", x, I0), 1e-8, I0, dat["I0"][i])
	}

	// check I1
	io.Pl()
	for i, x := range dat["x"] {
		I1 := ModBesselI1(x)
		chk.Float64(tst, io.Sf("I1(%11.8f)=%23.15e", x, I1), 1e-7, I1, dat["I1"][i])
	}

	// check I2
	io.Pl()
	for i, x := range dat["x"] {
		I2 := ModBesselIn(2, x)
		chk.Float64(tst, io.Sf("I2(%11.8f)=%23.15e", x, I2), 1e-7, I2, dat["I2"][i])
	}

	// check I3
	io.Pl()
	for i, x := range dat["x"] {
		I3 := ModBesselIn(3, x)
		chk.Float64(tst, io.Sf("I3(%11.8f)=%23.15e", x, I3), 1e-7, I3, dat["I3"][i])
	}

	// check K0
	io.Pl()
	for i, x := range dat["x"] {
		K0 := ModBesselK0(x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(K0, +1) {
				tst.Errorf("K0(0) should be +Inf. %v is incorrect\n", K0)
				return
			}
		} else {
			chk.Float64(tst, io.Sf("K0(%11.8f)=%23.15e", x, K0), 1e-15, K0, dat["K0"][i])
		}
	}

	// check K1
	io.Pl()
	for i, x := range dat["x"] {
		K1 := ModBesselK1(x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(K1, +1) {
				tst.Errorf("K1(0) should be +Inf. %v is incorrect\n", K1)
				return
			}
		} else {
			chk.Float64(tst, io.Sf("K1(%11.8f)=%23.15e", x, K1), 1e-15, K1, dat["K1"][i])
		}
	}

	// check K2
	io.Pl()
	for i, x := range dat["x"] {
		K2 := ModBesselKn(2, x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(K2, +1) {
				tst.Errorf("K2(0) should be +Inf. %v is incorrect\n", K2)
				return
			}
		} else {
			chk.Float64(tst, io.Sf("K2(%11.8f)=%23.15e", x, K2), 1e-15, K2, dat["K2"][i])
		}
	}

	// check K3
	io.Pl()
	for i, x := range dat["x"] {
		K3 := ModBesselKn(3, x)
		if math.Abs(x) < 1e-14 {
			if !math.IsInf(K3, +1) {
				tst.Errorf("K3(0) should be +Inf. %v is incorrect\n", K3)
				return
			}
		} else {
			chk.Float64(tst, io.Sf("K3(%11.8f)=%23.15e", x, K3), 1e-14, K3, dat["K3"][i])
		}
	}
}
