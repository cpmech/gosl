// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"gosl/chk"
	"gosl/io"
)

func Test_stat01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("stat01")

	x := []float64{100, 100, 102, 98, 77, 99, 70, 105, 98}

	xmin, xave, xmax, xdevS := StatBasic(x, true)
	_, _, _, xdevA := StatBasic(x, false)
	xdevS1 := StatDev(x, true)
	xdevA1 := StatDev(x, false)
	xave1 := StatAve(x)
	xave2a, xdev2a := StatAveDev(x, true)
	xave2b, xdev2b := StatAveDev(x, false)

	sum, mean, adev, sdev, vari, skew, kurt := StatMoments(x)
	io.Pforan("x    = %v\n", x)
	io.Pforan("sum  = %v\n", sum)
	io.Pforan("mean = %v  (%v)  (%v)\n", mean, xave, xave1)
	io.Pforan("adev = %v  (%v)  (%v)\n", adev, xdevA, xdevA1)
	io.Pforan("sdev = %v  (%v)  (%v)\n", sdev, xdevS, xdevS1)
	io.Pforan("vari = %v\n", vari)
	io.Pforan("skew = %v\n", skew)
	io.Pforan("kurt = %v\n", kurt)
	chk.Float64(tst, "sum  ", 1e-17, sum, 849)
	chk.Float64(tst, "mean ", 1e-17, mean, 849.0/9.0)
	chk.Float64(tst, "sdev ", 1e-17, sdev, 12.134661099511597)
	chk.Float64(tst, "vari ", 1e-17, vari, 147.25)
	chk.Float64(tst, "xmin ", 1e-17, xmin, 70)
	chk.Float64(tst, "xave ", 1e-17, xave, 849.0/9.0)
	chk.Float64(tst, "xave1", 1e-17, xave1, 849.0/9.0)
	chk.Float64(tst, "xmax ", 1e-17, xmax, 105)
	chk.Float64(tst, "xdevA", 1e-17, xdevA, adev)
	chk.Float64(tst, "xdevS", 1e-17, xdevS, 12.134661099511597)
	chk.Float64(tst, "xdevA1", 1e-17, xdevA1, adev)
	chk.Float64(tst, "xdevS1", 1e-17, xdevS1, 12.134661099511597)
	chk.Float64(tst, "xave2a", 1e-17, xave2a, xave)
	chk.Float64(tst, "xdev2a", 1e-17, xdev2a, xdevS)
	chk.Float64(tst, "xave2b", 1e-17, xave2b, xave)
	chk.Float64(tst, "xdev2b", 1e-17, xdev2b, xdevA)
	// TODO: add checks for adev, skew and kurt
}

func Test_stat02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("stat02")

	x := [][]float64{
		{100, 100, 102, 98, 77, 99, 70, 105, 98},
		{80, 101, 12, 58, 47, 80, 20, 111, 89},
		{50, 130, 72, 38, 71, 15, 10, 12, 55},
	}

	y, z := StatTable(x, true, true)
	io.Pforan("\nmin\n")
	chk.Float64(tst, "y00=min(x[0,:])", 1e-17, y[0][0], 70)
	chk.Float64(tst, "y01=min(x[1,:])", 1e-17, y[0][1], 12)
	chk.Float64(tst, "y02=min(x[2,:])", 1e-17, y[0][2], 10)
	io.Pforan("\nave\n")
	chk.Float64(tst, "y10=ave(x[0,:])", 1e-17, y[1][0], 849.0/9.0)
	chk.Float64(tst, "y11=ave(x[1,:])", 1e-17, y[1][1], 598.0/9.0)
	chk.Float64(tst, "y12=ave(x[2,:])", 1e-17, y[1][2], 453.0/9.0)
	io.Pforan("\nmax\n")
	chk.Float64(tst, "y20=max(x[0,:])", 1e-17, y[2][0], 105)
	chk.Float64(tst, "y21=max(x[1,:])", 1e-17, y[2][1], 111)
	chk.Float64(tst, "y22=max(x[2,:])", 1e-17, y[2][2], 130)
	io.Pforan("\ndev\n")
	chk.Float64(tst, "y30=dev(x[0,:])", 1e-17, y[3][0], 12.134661099511597)
	chk.Float64(tst, "y31=dev(x[1,:])", 1e-17, y[3][1], 34.688294535444918)
	chk.Float64(tst, "y32=dev(x[2,:])", 1e-17, y[3][2], 38.343839140075687)
	io.Pfyel("\nmin\n")
	chk.Float64(tst, "z00=min(y[0,:])=min(min)", 1e-17, z[0][0], 10)
	chk.Float64(tst, "z01=min(y[1,:])=min(ave)", 1e-17, z[0][1], 453.0/9.0)
	chk.Float64(tst, "z02=min(y[2,:])=min(max)", 1e-17, z[0][2], 105)
	chk.Float64(tst, "z03=min(y[3,:])=min(dev)", 1e-17, z[0][3], 12.134661099511597)
	io.Pfyel("\nave\n")
	chk.Float64(tst, "z10=ave(y[0,:])=ave(min)", 1e-17, z[1][0], 92.0/3.0)
	chk.Float64(tst, "z11=ave(y[1,:])=ave(ave)", 1e-17, z[1][1], ((849.0+598.0+453.0)/9.0)/3.0)
	chk.Float64(tst, "z12=ave(y[2,:])=ave(max)", 1e-17, z[1][2], 346.0/3.0)
	chk.Float64(tst, "z13=ave(y[3,:])=ave(dev)", 1e-17, z[1][3], (12.134661099511597+34.688294535444918+38.343839140075687)/3.0)
	io.Pfyel("\nmax\n")
	chk.Float64(tst, "z20=max(y[0,:])=max(min)", 1e-17, z[2][0], 70)
	chk.Float64(tst, "z21=max(y[1,:])=max(ave)", 1e-17, z[2][1], 849.0/9.0)
	chk.Float64(tst, "z22=max(y[2,:])=max(max)", 1e-17, z[2][2], 130)
	chk.Float64(tst, "z23=max(y[3,:])=max(dev)", 1e-17, z[2][3], 38.343839140075687)
	io.Pfyel("\ndev\n")
	chk.Float64(tst, "z30=dev(y[0,:])=dev(min)", 1e-17, z[3][0], 34.078341117685483)
	chk.Float64(tst, "z31=dev(y[1,:])=dev(ave)", 1e-17, z[3][1], 22.261169573539771)
	chk.Float64(tst, "z32=dev(y[2,:])=dev(max)", 1e-17, z[3][2], 13.051181300301263)
	chk.Float64(tst, "z33=dev(y[3,:])=dev(dev)", 1e-17, z[3][3], 14.194778389023206)
}
