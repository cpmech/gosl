// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestPostProc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Postp01. GetITout")

	Tout := []float64{0, 0.1, 0.2, 0.200001, 0.201, 0.3001, 0.8, 0.99, 0.999, 1}
	Tsel := []float64{0, 0.2, 0.3, 0.6, 0.8, 0.9, 0.99, -1}

	tol := 0.001
	I, T := GetITout(Tout, Tsel, tol)
	io.Pfcyan("Tout = %v\n", Tout)
	io.Pfcyan("Tsel = %v\n", Tsel)
	io.Pforan("I = %v\n", I)
	io.Pforan("T = %v\n", T)

	chk.Ints(tst, "I", I, []int{0, 2, 5, 6, 7, 9})
	chk.Array(tst, "T", 1e-16, T, []float64{0, 0.2, 0.3001, 0.8, 0.99, 1})
}

func TestPostProc02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Postp02. GetStrides")

	nTotal := 2000
	nRequired := 5
	I := GetStrides(nTotal, nRequired)
	io.Pf("I = %v\n", I)
	chk.Ints(tst, "I", I, []int{0, 400, 800, 1200, 1600, 2000})

	nTotal = 2001
	nRequired = 5
	I = GetStrides(nTotal, nRequired)
	io.Pf("I = %v\n", I)
	chk.Ints(tst, "I", I, []int{0, 400, 800, 1200, 1600, 2001})
}

func TestOutputter01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Outputter01")

	out := func(u []float64, t float64) {
		io.Pf("▶▶▶▶ output @ t = %v\n", t)
		u[0], u[1] = 100+t, 200+t
	}
	nu := 2
	dt := 0.1
	dtOut := 0.2
	tmax := 1.0
	o := NewOutputter(dt, dtOut, tmax, nu, out)

	chk.Float64(tst, "Dt", 1e-17, o.Dt, 0.1)
	chk.Float64(tst, "DtOut", 1e-17, o.DtOut, 0.2)
	chk.Float64(tst, "Tmax", 1e-17, o.Tmax, 1.0)
	chk.Int(tst, "Nsteps", o.Nsteps, 10)
	chk.Int(tst, "Every", o.Every, 2)
	chk.Int(tst, "Tidx", o.Tidx, 1)
	chk.Int(tst, "Nmax", o.Nmax, 6)
	chk.Int(tst, "Idx", o.Idx, 1)
	chk.Int(tst, "len(T)", len(o.T), o.Nmax)
	chk.Int(tst, "len(U)", len(o.U), o.Nmax)
	chk.Int(tst, "len(U[0])", len(o.U[0]), nu)

	t := 0.0
	for idxT := 0; idxT < o.Nsteps; idxT++ {
		t += o.Dt
		io.Pforan("%2d : t = %v\n", idxT, t)
		o.MaybeNow(idxT, t)
	}
	io.Pfblue2("final t = %v\n", t)
	io.Pfblue2("Idx     = %v\n", o.Idx)
	chk.Int(tst, "Idx=Nmax", o.Idx, o.Nmax)
	chk.Float64(tst, "t", 1e-15, t, 1.0)
	chk.Array(tst, "T", 1e-15, o.T, []float64{0, 0.2, 0.4, 0.6, 0.8, 1.0})
	chk.Deep2(tst, "U", 1e-17, o.U, [][]float64{
		{100.0, 200.0},
		{100.2, 200.2},
		{100.4, 200.4},
		{100.6, 200.6},
		{100.8, 200.8},
		{101.0, 201.0},
	})
}

func TestOutputter02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Outputter02")

	out := func(u []float64, t float64) {
		io.Pf("▶▶▶▶ output @ t = %v\n", t)
	}
	nu := 1
	dt := 0.12
	dtOut := 0.2
	tmax := 1.0
	o := NewOutputter(dt, dtOut, tmax, nu, out)

	chk.Float64(tst, "Tmax", 1e-17, o.Tmax, 1.08)
	chk.Int(tst, "Nsteps", o.Nsteps, 9)
	chk.Int(tst, "Every", o.Every, 1)
	chk.Int(tst, "Tidx", o.Tidx, 0)
	chk.Int(tst, "Nmax", o.Nmax, 10)
	chk.Int(tst, "Idx", o.Idx, 1)

	t := 0.0
	for idxT := 0; idxT < o.Nsteps; idxT++ {
		t += o.Dt
		io.Pforan("%2d : t = %v\n", idxT, t)
		o.MaybeNow(idxT, t)
	}
	io.Pfblue2("final t = %v\n", t)
	io.Pfblue2("Idx     = %v\n", o.Idx)
	chk.Int(tst, "Idx=Nmax", o.Idx, o.Nmax)
	chk.Float64(tst, "t", 1e-15, t, o.Tmax)
	chk.Array(tst, "T", 1e-15, o.T, []float64{0, 0.12, 0.24, 0.36, 0.48, 0.6, 0.72, 0.84, 0.96, 1.08})
}

func TestOutputter03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Outputter03")

	out := func(u []float64, t float64) {
		io.Pf("▶▶▶▶ output @ t = %v\n", t)
	}
	nu := 1
	dt := 0.12
	dtOut := 0.4
	tmax := 1.0
	o := NewOutputter(dt, dtOut, tmax, nu, out)

	chk.Float64(tst, "Tmax", 1e-17, o.Tmax, 1.08)
	chk.Int(tst, "Nsteps", o.Nsteps, 9)
	chk.Int(tst, "Every", o.Every, 3)
	chk.Int(tst, "Tidx", o.Tidx, 2)
	chk.Int(tst, "Nmax", o.Nmax, 4)
	chk.Int(tst, "Idx", o.Idx, 1)

	t := 0.0
	for idxT := 0; idxT < o.Nsteps; idxT++ {
		t += o.Dt
		io.Pforan("%2d : t = %v\n", idxT, t)
		o.MaybeNow(idxT, t)
	}
	io.Pfblue2("final t = %v\n", t)
	io.Pfblue2("Idx     = %v\n", o.Idx)
	chk.Int(tst, "Idx=Nmax", o.Idx, o.Nmax)
	chk.Float64(tst, "t", 1e-15, t, o.Tmax)
	chk.Array(tst, "T", 1e-15, o.T, []float64{0, 0.36, 0.72, 1.08})
}

func TestOutputter04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Outputter04")

	out := func(u []float64, t float64) {
		io.Pf("▶▶▶▶ output @ t = %v\n", t)
	}
	nu := 0
	dt := 0.1738
	dtOut := 0.8485
	tmax := 1.537
	o := NewOutputter(dt, dtOut, tmax, nu, out)

	chk.Float64(tst, "Tmax", 1e-17, o.Tmax, 1.5642)
	chk.Int(tst, "Nsteps", o.Nsteps, 9)
	chk.Int(tst, "Every", o.Every, 4)
	chk.Int(tst, "Tidx", o.Tidx, 3)
	chk.Int(tst, "Nmax", o.Nmax, 4)
	chk.Int(tst, "Idx", o.Idx, 1)

	t := 0.0
	for idxT := 0; idxT < o.Nsteps; idxT++ {
		t += o.Dt
		io.Pforan("%2d : t = %v\n", idxT, t)
		o.MaybeNow(idxT, t)
	}
	io.Pfblue2("final t = %v\n", t)
	io.Pfblue2("Idx     = %v\n", o.Idx)
	chk.Int(tst, "Idx=Nmax", o.Idx, o.Nmax)
	chk.Float64(tst, "t", 1e-15, t, o.Tmax)
	chk.Array(tst, "T", 1e-15, o.T, []float64{0, 4 * 0.1738, 8 * 0.1738, o.Tmax})
}

func TestOutputter05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Outputter05")

	out := func(u []float64, t float64) {
		io.Pf("▶▶▶▶ output @ t = %v\n", t)
	}
	nu := 0
	dt := 0.55
	dtOut := 0.1
	tmax := 1.2
	o := NewOutputter(dt, dtOut, tmax, nu, out)

	chk.Float64(tst, "Dt", 1e-17, o.Dt, 0.55)
	chk.Float64(tst, "DtOut", 1e-17, o.DtOut, 0.55)
	chk.Float64(tst, "Tmax", 1e-15, o.Tmax, 1.65)
	chk.Int(tst, "Nsteps", o.Nsteps, 3)
	chk.Int(tst, "Every", o.Every, 1)
	chk.Int(tst, "Tidx", o.Tidx, 0)
	chk.Int(tst, "Nmax", o.Nmax, 4)
	chk.Int(tst, "Idx", o.Idx, 1)

	t := 0.0
	for idxT := 0; idxT < o.Nsteps; idxT++ {
		t += o.Dt
		io.Pforan("%2d : t = %v\n", idxT, t)
		o.MaybeNow(idxT, t)
	}
	io.Pfblue2("final t = %v\n", t)
	io.Pfblue2("Idx     = %v\n", o.Idx)
	chk.Int(tst, "Idx=Nmax", o.Idx, o.Nmax)
	chk.Float64(tst, "t", 1e-15, t, o.Tmax)
	chk.Array(tst, "T", 1e-15, o.T, []float64{0, 0.55, 1.1, 1.65})
}
