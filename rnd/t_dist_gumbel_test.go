// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func Test_dist_gumbel_01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_gumbel_01")

	_, dat := io.ReadTable("data/gumbel.dat")

	X, ok := dat["x"]
	if !ok {
		tst.Errorf("cannot get x values\n")
		return
	}
	U, ok := dat["u"]
	if !ok {
		tst.Errorf("cannot get u values\n")
		return
	}
	B, ok := dat["b"]
	if !ok {
		tst.Errorf("cannot get b values\n")
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

	var dist DistGumbel

	nx := len(X)
	for i := 0; i < nx; i++ {
		dist.U = U[i]
		dist.B = B[i]
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

func Test_gumbel_03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("dist_gumbel_03")

	var dist DistGumbel
	dist.Init(&Variable{M: 61.3, S: 7.52}) // from Haldar & Mahadevan page 90
	io.Pforan("dist = %+#v\n", dist)
	chk.Float64(tst, "u", 0.00011, dist.U, 57.9157)
	chk.Float64(tst, "β", 1e-4, dist.B, 1.0/0.17055)
}
