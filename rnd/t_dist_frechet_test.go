// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/io"
)

func Test_dist_frechet_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_frechet_01")

	_, dat := io.ReadTable("data/frechet.dat")

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	L, ok := dat["l"] // location
	if !ok {
		tst.Errorf("cannot get l values\n")
		return
	}
	C, ok := dat["c"] // scale
	if !ok {
		tst.Errorf("cannot get c values\n")
		return
	}
	A, ok := dat["a"] // shape
	if !ok {
		tst.Errorf("cannot get a values\n")
		return
	}
	YpdfCmp, ok := dat["ypdf"]
	if !ok {
		tst.Errorf("cannot get ypdf values\n")
		return
	}
	YcdfCmp, ok := dat["ycdf"]
	if !ok {
		tst.Errorf("cannot get ycdf values\n")
		return
	}

	var dist DistFrechet

	nx := len(X)
	for i := 0; i < nx; i++ {
		dist.Init(&Variable{L: L[i], C: C[i], A: A[i]})
		Ypdf := dist.Pdf(X[i])
		Ycdf := dist.Cdf(X[i])
		err := chk.PrintAnaNum("ypdf", 1e-14, YpdfCmp[i], Ypdf, chk.Verbose)
		if err != nil {
			tst.Errorf("pdf failed: %v\n", err)
			return
		}
		err = chk.PrintAnaNum("ycdf", 1e-15, YcdfCmp[i], Ycdf, chk.Verbose)
		if err != nil {
			tst.Errorf("cdf failed: %v\n", err)
			return
		}
	}
}

func Test_frechet_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_frechet_03")

	μ := 10.0
	σ := 5.0
	δ := σ / μ
	d := 1.0 + δ*δ
	io.Pforan("μ=%v σ=%v δ=%v d=%v\n", μ, σ, δ, d)

	k := 0.2441618
	α := 1.0 / k
	l := μ - math.Gamma(1.0-k)
	io.Pfpink("l=%v α=%v\n", l, α)

	l = 8.782275
	α = 4.095645

	var dist DistFrechet
	dist.Init(&Variable{L: l, A: α})
	io.Pforan("dist = %+#v\n", dist)
	io.Pforan("mean = %v\n", dist.Mean())
	io.Pforan("var  = %v\n", dist.Variance())
	io.Pforan("σ    = %v\n", math.Sqrt(dist.Variance()))
}
