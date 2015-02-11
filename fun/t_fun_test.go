// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	//T_FUN_SAVE = true
	T_FUN_SAVE = false
)

func Test_fun01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun01. Decreasing Reference Model")

	ya := 1.0
	yb := -0.5
	λ1 := 1.0

	o := New("ref-dec-gen", []*Prm{
		&Prm{N: "bet", V: 5.0},
		&Prm{N: "a", V: -λ1},
		&Prm{N: "b", V: -1.0},
		&Prm{N: "c", V: ya},
		&Prm{N: "A", V: 0.0},
		&Prm{N: "B", V: λ1},
		&Prm{N: "xini", V: 0.0},
		&Prm{N: "yini", V: yb},
	})

	tmax := 3.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		plt.Reset()
		withG, withH, save, show := true, true, false, true
		PlotT(o, "/tmp/gosl", "ref-dec-gen-01.png", 0.0, tmax, xcte, 201, "", withG, withH, save, show, func() {
			plt.Plot([]float64{0, tmax}, []float64{ya, ya - λ1*tmax}, "'k-'")
			plt.Equal()
		})
	}

	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := true
	CheckT(tst, o, 0.0, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_fun02(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun02. Dec Ref Model (specialised)")

	ya := 1.0
	yb := -50.0
	λ1 := 1.0

	o := New("ref-dec-sp1", []*Prm{
		&Prm{N: "bet", V: 5.0},
		&Prm{N: "lam1", V: λ1},
		&Prm{N: "ya", V: ya},
		&Prm{N: "yb", V: yb},
	})

	tmin := 0.0
	tmax := 300.0
	//tmax := 140.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		plt.Reset()
		withG, withH, save, show := true, true, false, true
		PlotT(o, "/tmp/gosl", "ref-dec-sp1-01.png", tmin, tmax, xcte, 201, "lw=2,color='orange'", withG, withH, save, show, func() {
			plt.Plot([]float64{0, tmax}, []float64{ya, ya - λ1*tmax}, "'k--'")
			plt.Equal()
		})
	}

	if true {
		//if false {
		sktol := 1e-10
		dtol := 1e-10
		dtol2 := 1e-10
		ver := true
		CheckT(tst, o, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	}
}

func Test_fun03(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun03. add, cte, srmps")

	cte := New("cte", []*Prm{&Prm{N: "C", V: 30}})
	srmps := New("srmps", []*Prm{
		&Prm{N: "ca", V: 0},
		&Prm{N: "cb", V: 1},
		&Prm{N: "ta", V: 0},
		&Prm{N: "tb", V: 1},
	})
	add := New("add", []*Prm{
		&Prm{N: "a", V: 1},
		&Prm{N: "b", V: 1},
		&Prm{N: "fa", Fcn: cte},
		&Prm{N: "fb", Fcn: srmps},
	})

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		withG, withH, save, show := true, true, false, true
		plt.Reset()
		PlotT(cte, "/tmp/gosl", "fun-cte-01.png", tmin, tmax, xcte, 41, "", withG, withH, save, show, nil)
		plt.Reset()
		PlotT(srmps, "/tmp/gosl", "fun-srmps-01.png", tmin, tmax, xcte, 41, "", withG, withH, save, show, nil)
		plt.Reset()
		PlotT(add, "/tmp/gosl", "fun-add-01.png", tmin, tmax, xcte, 41, "", withG, withH, save, show, nil)
	}

	if true {
		//if false {
		sktol := 1e-10
		dtol := 1e-10
		dtol2 := 1e-9
		ver := true
		tskip := []float64{tmin, tmax}
		CheckT(tst, cte, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
		utl.Pf("\n")
		CheckT(tst, srmps, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
		utl.Pf("\n")
		CheckT(tst, add, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
	}
}

func Test_fun04(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun04. lin")

	lin := New("lin", []*Prm{
		&Prm{N: "m", V: 0.5},
		&Prm{N: "ts", V: 0},
	})

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		plt.Reset()
		withG, withH, save, show := true, true, false, true
		PlotT(lin, "/tmp/gosl", "fun-lin-01.png", tmin, tmax, xcte, 11, "", withG, withH, save, show, nil)
	}

	if true {
		//if false {
		sktol := 1e-10
		dtol := 1e-10
		dtol2 := 1e-10
		ver := true
		CheckT(tst, lin, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	}
}

func Test_fun05(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun05. zero and one")

	utl.Pforan("Zero(666,nil) = %v\n", Zero.F(666, nil))
	utl.Pforan("One(666,nil)  = %v\n", One.F(666, nil))
	utl.CheckScalar(tst, "zero", 1e-17, Zero.F(666, nil), 0)
	utl.CheckScalar(tst, "one ", 1e-17, One.F(666, nil), 1)
}

func Test_fun06(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun06. pts")

	fun := New("pts", []*Prm{
		&Prm{N: "t0", V: 0.00}, {N: "y0", V: 0.50},
		&Prm{N: "t1", V: 1.00}, {N: "y1", V: 0.20},
		&Prm{N: "t2", V: 2.00}, {N: "y2", V: 0.20},
		&Prm{N: "t3", V: 3.00}, {N: "y3", V: 0.05},
		&Prm{N: "t4", V: 4.00}, {N: "y4", V: 0.01},
		&Prm{N: "t5", V: 5.00}, {N: "y5", V: 0.00},
	})

	tmin := 0.0
	tmax := 5.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		plt.Reset()
		withG, withH, save, show := true, true, false, true
		PlotT(fun, "/tmp/gosl", "fun-pts-01.png", tmin, tmax, xcte, 6, "'o-'", withG, withH, save, show, nil)
	}

	if true {
		tmin = 0.01
		tmax = 4.99
		//if false {
		sktol := 1e-10
		dtol := 1e-10
		dtol2 := 1e-10
		ver := true
		CheckT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	}
}

func Test_fun07(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun07. pts")

	fun := New("exc1", []*Prm{
		&Prm{N: "A", V: 200},
		&Prm{N: "b", V: 2},
	})

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		plt.Reset()
		withG, withH, save, show := true, true, false, true
		PlotT(fun, "/tmp/gosl", "fun-exc1-01.png", tmin, tmax, xcte, 41, "'o-'", withG, withH, save, show, nil)
	}

	if true {
		sktol := 1e-10
		dtol := 1e-7
		dtol2 := 1e-6
		ver := true
		CheckT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	}
}

func Test_fun08(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("fun08. pts")

	fun := New("exc2", []*Prm{
		&Prm{N: "ta", V: 5},
		&Prm{N: "A", V: 3},
		&Prm{N: "b", V: 0.2},
	})

	tmin := 0.0
	tmax := 7.0
	xcte := []float64{0, 0, 0}
	if T_FUN_SAVE {
		plt.Reset()
		withG, withH, save, show := true, true, false, true
		PlotT(fun, "/tmp/gosl", "fun-exc2-01.png", tmin, tmax, xcte, 41, "'o-'", withG, withH, save, show, nil)
	}

	if true {
		sktol := 1e-10
		dtol := 1e-10
		dtol2 := 1e-10
		ver := true
		CheckT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	}
}
