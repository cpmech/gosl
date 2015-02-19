// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"bytes"
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_functions01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("functions01")

	x := utl.LinSpace(-2, 2, 21)
	ym := make([]float64, len(x))
	yh := make([]float64, len(x))
	ys := make([]float64, len(x))
	yAbs2Ramp := make([]float64, len(x))
	yHea2Ramp := make([]float64, len(x))
	ySig2Heav := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		ym[i] = Ramp(x[i])
		yh[i] = Heav(x[i])
		ys[i] = Sign(x[i])
		yAbs2Ramp[i] = (x[i] + math.Abs(x[i])) / 2.0
		yHea2Ramp[i] = x[i] * yh[i]
		ySig2Heav[i] = (1.0 + ys[i]) / 2.0
	}
	chk.Vector(tst, "abs => ramp", 1e-17, ym, yAbs2Ramp)
	chk.Vector(tst, "hea => ramp", 1e-17, ym, yHea2Ramp)
	chk.Vector(tst, "sig => heav", 1e-17, yh, ySig2Heav)

	var b bytes.Buffer
	io.Ff(&b, "from gosl_fig import *\n")
	utl.Gen4Arrays(&b, "x", "ym", "yh", "ys", x, ym, yh, ys)
	io.Ff(&b, "subplot(3,1,1)\n")
	io.Ff(&b, "plot(x,ym,label='Ramp/Macaulay',clip_on=0,lw=2,marker='o')\n")
	io.Ff(&b, "axis([axis()[0],axis()[1],-0.1,axis()[3]])\n")
	io.Ff(&b, "Cross()\n")
	io.Ff(&b, "Gll('x','y',leg_loc='upper left')\n")
	io.Ff(&b, "subplot(3,1,2)\n")
	io.Ff(&b, "plot(x,yh,label='Heaviside',clip_on=0,lw=2,marker='o')\n")
	io.Ff(&b, "axis([axis()[0],axis()[1],-0.1,1.1])\n")
	io.Ff(&b, "Cross()\n")
	io.Ff(&b, "Gll('x','y',leg_loc='upper left')\n")
	io.Ff(&b, "subplot(3,1,3)\n")
	io.Ff(&b, "plot(x,ys,label='Sign',clip_on=0,lw=2,marker='o')\n")
	io.Ff(&b, "axis([axis()[0],axis()[1],-1.1,1.1])\n")
	io.Ff(&b, "Cross()\n")
	io.Ff(&b, "Gll('x','y',leg_loc='upper left')\n")
	io.Ff(&b, "show()\n")
	io.WriteFileD("/tmp/gosl/", "functions01.py", &b)
}

// numderiv employs a 1st order forward difference to approximate the derivative of f(x) w.r.t x @ x
func numderiv(f func(x float64) float64, x float64) float64 {
	eps, cte1 := 1e-16, 1e-5
	delta := math.Sqrt(eps * max(cte1, math.Abs(x)))
	return (f(x+delta) - f(x)) / delta
}

func Test_functions02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("functions02")

	β := 6.0
	f := func(x float64) float64 { return Sramp(x, β) }
	ff := func(x float64) float64 { return SrampD1(x, β) }

	np := 401
	//x  := utl.LinSpace(-5e5, 5e5, np)
	//x  := utl.LinSpace(-5e2, 5e2, np)
	x := utl.LinSpace(-5e1, 5e1, np)
	y := make([]float64, np)
	g := make([]float64, np)
	h := make([]float64, np)
	tolg, tolh := 1e-6, 1e-5
	with_err := false
	for i := 0; i < np; i++ {
		y[i] = Sramp(x[i], β)
		g[i] = SrampD1(x[i], β)
		h[i] = SrampD2(x[i], β)
		gnum := numderiv(f, x[i])
		hnum := numderiv(ff, x[i])
		errg := math.Abs(g[i] - gnum)
		errh := math.Abs(h[i] - hnum)
		clrg, clrh := "[1;32m", "[1;32m"
		if errg > tolg {
			clrg, with_err = "[1;31m", true
		}
		if errh > tolh {
			clrh, with_err = "[1;31m", true
		}
		io.Pf("errg = %s%23.15e   errh = %s%23.15e[0m\n", clrg, errg, clrh, errh)
	}

	var b bytes.Buffer
	io.Ff(&b, "from gosl_fig import *\n")
	utl.Gen4Arrays(&b, "x", "y", "g", "h", x, y, g, h)
	io.Ff(&b, "subplot(3,1,1)\n")
	io.Ff(&b, "plot(x,y, 'b-', lw=2)\n")
	io.Ff(&b, "axis('equal')\n")
	io.Ff(&b, "Cross()\n")
	io.Ff(&b, "Gll('x','y',leg=0)\n")
	io.Ff(&b, "subplot(3,1,2)\n")
	io.Ff(&b, "plot(x,g, 'b-', lw=2)\n")
	io.Ff(&b, "Gll('x','g',leg=0)\n")
	io.Ff(&b, "subplot(3,1,3)\n")
	io.Ff(&b, "plot(x,h, 'b-', lw=2)\n")
	io.Ff(&b, "Gll('x','h',leg=0)\n")
	io.Ff(&b, "show()\n")
	io.WriteFileD("/tmp/gosl/", "functions02.py", &b)
	io.PfBlue("file <results/functions02.py> saved\n")

	if with_err {
		chk.Panic("errors found")
	}
}
