// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
)

func d2r(deg float64) float64 { return deg * math.Pi / 180.0 }

func Test_elliptic01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("elliptic01")

	φ, k := d2r(5), 0.0
	F := Elliptic1(φ, k)
	io.Pf("F(5°,0) = %v\n", F)
	chk.Float64(tst, "F(5°,0)", 1e-15, F, φ)

	φ, k = d2r(30), math.Sin(d2r(68))
	F = Elliptic1(φ, k)
	io.Pf("F(30°,68°) = %v\n", F)
	chk.Float64(tst, "F(30°,68°)", 1e-9, F, 0.54527182)

	φ, k = d2r(65), math.Sin(d2r(88))
	F = Elliptic1(φ, k)
	io.Pf("F(65°,88°) = %v\n", F)
	chk.Float64(tst, "F(65°,88°)", 1e-8, F, 1.50537033)

	φ, k = d2r(85), math.Sin(d2r(90))
	F = Elliptic1(φ, k)
	io.Pf("F(85°,90°) = %v\n", F)
	chk.Float64(tst, "F(85°,90°)", 1e-8, F, 3.13130133)

	p90 := d2r(90)
	φ, k = p90, math.Sin(d2r(90))
	F = Elliptic1(φ, k)
	io.Pf("F(90°,90°) = %v\n", F)
	if F != math.Inf(1) {
		tst.Errorf("F(90°,90°) should be +Inf")
		return
	}

	// load data
	_, dat := io.ReadTable("data/as-17-elliptic-integrals-table17.5-small.cmp")
	for i, p := range dat["phi"] {
		k := dat["k"][i]
		F := Elliptic1(p, k)
		if math.Abs(k-1.0) < 1e-15 && math.Abs(p-p90) < 1e-15 {
			if F != math.Inf(1) {
				tst.Errorf("F(90°,1) should be +Inf")
				return
			}
		} else {
			chk.Float64(tst, io.Sf("F(%.8f,%.8f)=%23.15e", p, k, F), 1e-14, F, dat["F"][i])
		}
	}
}

func Test_elliptic02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("elliptic02")

	φ, k := d2r(5), 0.0
	E := Elliptic2(φ, k)
	io.Pf("E(5°,0) = %v\n", E)
	chk.Float64(tst, "E(5°,0)", 1e-15, E, φ)

	φ, k = d2r(30), math.Sin(d2r(68))
	E = Elliptic2(φ, k)
	io.Pf("E(30°,68°) = %v\n", E)
	chk.Float64(tst, "E(30°,68°)", 1e-9, E, 0.50343686)

	φ, k = d2r(65), math.Sin(d2r(88))
	E = Elliptic2(φ, k)
	io.Pf("E(65°,88°) = %v\n", E)
	chk.Float64(tst, "E(65°,88°)", 1e-8, E, 0.90667305)

	φ, k = d2r(85), math.Sin(d2r(90))
	E = Elliptic2(φ, k)
	io.Pf("E(85°,90°) = %v\n", E)
	chk.Float64(tst, "E(85°,90°)", 1e-8, E, 0.99619470)

	p90 := d2r(90)
	φ, k = p90, math.Sin(d2r(90))
	E = Elliptic2(φ, k)
	io.Pf("E(90°,90°) = %v\n", E)
	chk.Float64(tst, "E(85°,90°)", 1e-8, E, 1)

	// load data
	_, dat := io.ReadTable("data/as-17-elliptic-integrals-table17.6-small.cmp")
	//_, dat := io.ReadTable("data/as-17-elliptic-integrals-table17.6-big.cmp")
	for i, p := range dat["phi"] {
		k := dat["k"][i]
		E := Elliptic2(p, k)
		chk.Float64(tst, io.Sf("E(%.8f,%.8f)=%23.15e", p, k, E), 1e-14, E, dat["E"][i])
	}
}

func Test_elliptic03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("elliptic03")

	n, φ, k := 0.1, d2r(0), 0.0
	P := Elliptic3(n, φ, k)
	io.Pf("Π(0,15°,0) = %v\n", P)
	chk.Float64(tst, "Π(0,15°,0)", 1e-15, P, 0)

	n, φ, k = 0.1, d2r(15), 0.0
	P = Elliptic3(n, φ, k)
	io.Pf("Π(0,90°,0) = %v\n", P)
	chk.Float64(tst, "Π(0,90°,0)", 1e-5, P, 0.26239)

	// load data
	_, dat := io.ReadTable("data/as-17-elliptic-integrals-table17.9-small.cmp")
	//_, dat := io.ReadTable("data/as-17-elliptic-integrals-table17.9-big.cmp")
	for i, n := range dat["n"] {
		p := dat["phi"][i]
		k := dat["k"][i]
		P := Elliptic3(n, p, k)
		s := math.Sin(p)
		if math.Abs(k*s-1.0) < 1e-15 {
			if P != math.Inf(1) {
				tst.Errorf("Π(n,90°,1) should be +Inf")
				return
			}
		} else if math.Abs(n*s-1.0) < 1e-15 {
			if P != math.Inf(1) {
				tst.Errorf("Π(1,90°,k) should be +Inf")
				return
			}
		} else {
			chk.Float64(tst, io.Sf("Π(%.2f,%.8f,%.8f)=%23.15e", n, p, k, P), 1e-14, P, dat["PI"][i])
		}
	}
}
