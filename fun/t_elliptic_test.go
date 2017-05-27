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

func Test_elliptic01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("elliptic01")

	d2r := func(deg float64) float64 { return deg * math.Pi / 180.0 }

	φ, k := d2r(5), 0.0
	F := Elliptic1(φ, k)
	io.Pf("F(5°,0) = %v\n", F)
	chk.Scalar(tst, "F(5°,0)", 1e-15, F, φ)

	φ, k = d2r(30), math.Sin(d2r(68))
	F = Elliptic1(φ, k)
	io.Pf("F(30°,68°) = %v\n", F)
	chk.Scalar(tst, "F(30°,68°)", 1e-9, F, 0.54527182)

	φ, k = d2r(65), math.Sin(d2r(88))
	F = Elliptic1(φ, k)
	io.Pf("F(65°,88°) = %v\n", F)
	chk.Scalar(tst, "F(65°,88°)", 1e-8, F, 1.50537033)

	φ, k = d2r(85), math.Sin(d2r(90))
	F = Elliptic1(φ, k)
	io.Pf("F(85°,90°) = %v\n", F)
	chk.Scalar(tst, "F(85°,90°)", 1e-8, F, 3.13130133)

	p90 := d2r(90)
	φ, k = p90, math.Sin(d2r(90))
	F = Elliptic1(φ, k)
	io.Pf("F(90°,90°) = %v\n", F)
	if F != math.Inf(1) {
		tst.Errorf("F(90°,90°) should be +Inf")
		return
	}

	// load data
	keys, dat, err := io.ReadTable("data/as-17-elliptic-integrals-table17.5-small.cmp")
	//keys, dat, err := io.ReadTable("data/as-17-elliptic-integrals-table17.5-big.cmp")
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	if false {
		io.Pf("Keys = %v\n", keys)
		io.Pf("dat = %v\n", dat["F"])
	}
	for i, p := range dat["phi"] {
		k := dat["k"][i]
		F := Elliptic1(p, k)
		if math.Abs(k-1.0) < 1e-15 && math.Abs(p-p90) < 1e-15 {
			if F != math.Inf(1) {
				tst.Errorf("F(90°,90°) should be +Inf")
				return
			}
		} else {
			chk.Scalar(tst, io.Sf("F(%.8f,%.8f)=%23.15e", p, k, F), 1e-14, F, dat["F"][i])
		}
	}
}
