// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_ts01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts01. Decreasing Reference Model")

	ya := 1.0
	yb := -0.5
	λ1 := 1.0

	o, err := New("ref-dec-gen", []*P{
		&P{N: "bet", V: 5.0},
		&P{N: "a", V: -λ1},
		&P{N: "b", V: -1.0},
		&P{N: "c", V: ya},
		&P{N: "A", V: 0.0},
		&P{N: "B", V: λ1},
		&P{N: "xini", V: 0.0},
		&P{N: "yini", V: yb},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmax := 3.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(o, "", "", 0.0, tmax, xcte, 201)
		plt.Subplot(3, 1, 1)
		plt.Plot([]float64{0, tmax}, []float64{ya, ya - λ1*tmax}, &plt.A{C: "k", Ls: "--"})
		plt.Equal()
		plt.Save("/tmp/gosl/fun", "ref-dec-gen")
	}
	//
	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, o, 0.0, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts02. Dec Ref Model (specialised)")

	ya := 1.0
	yb := -50.0
	λ1 := 1.0

	o, err := New("ref-dec-sp1", []*P{
		&P{N: "bet", V: 5.0},
		&P{N: "lam1", V: λ1},
		&P{N: "ya", V: ya},
		&P{N: "yb", V: yb},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 300.0
	//tmax := 140.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(o, "", "", tmin, tmax, xcte, 201)
		plt.Plot([]float64{0, tmax}, []float64{ya, ya - λ1*tmax}, &plt.A{C: "k", Ls: "--"})
		plt.Equal()
		plt.Save("/tmp/gosl/fun", "ref-dec-sp1")
	}

	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, o, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts03. add, cte, srmps")

	cte, err := New("cte", []*P{&P{N: "c", V: 30}})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	srmps, err := New("srmps", []*P{
		&P{N: "ca", V: 0},
		&P{N: "cb", V: 1},
		&P{N: "ta", V: 0},
		&P{N: "tb", V: 1},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	add, err := New("add", []*P{
		&P{N: "a", V: 1},
		&P{N: "b", V: 1},
		&P{N: "fa", Fcn: cte},
		&P{N: "fb", Fcn: srmps},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(cte, "/tmp/gosl/fun", "cte", tmin, tmax, xcte, 41)
		plt.Clf()
		PlotT(srmps, "/tmp/gosl/fun", "srmps", tmin, tmax, xcte, 41)
		plt.Clf()
		PlotT(add, "/tmp/gosl/fun", "add", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-9
	ver := chk.Verbose
	tskip := []float64{tmin, tmax}
	CheckDerivT(tst, cte, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	io.Pf("\n")
	CheckDerivT(tst, srmps, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
	io.Pf("\n")
	CheckDerivT(tst, add, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
}

func Test_ts04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts04. lin")

	lin, err := New("lin", []*P{
		&P{N: "m", V: 0.5},
		&P{N: "ts", V: 0},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(lin, "/tmp/gosl/fun", "lin", tmin, tmax, xcte, 11)
	}

	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, lin, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts05. zero and one")

	io.Pforan("Zero(666,nil) = %v\n", Zero.F(666, nil))
	io.Pforan("One(666,nil)  = %v\n", One.F(666, nil))
	chk.Float64(tst, "zero", 1e-17, Zero.F(666, nil), 0)
	chk.Float64(tst, "one ", 1e-17, One.F(666, nil), 1)
}

func Test_ts06a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts06a. pts")

	fun, err := New("pts", []*P{
		&P{N: "t", V: 0.00}, {N: "y", V: 0.50},
		&P{N: "t", V: 1.00}, {N: "y", V: 0.20},
		&P{N: "t", V: 2.00}, {N: "y", V: 0.20},
		&P{N: "t", V: 3.00}, {N: "y", V: 0.05},
		&P{N: "t", V: 4.00}, {N: "y", V: 0.01},
		&P{N: "t", V: 5.00}, {N: "y", V: 0.00},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := -1.0
	tmax := 6.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "pts", tmin, tmax, xcte, 8)
	}

	tmin = 0.01
	tmax = 4.99
	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts06b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts06b. pts")

	fun, err := New("pts", []*P{
		&P{N: "t0", V: 0.0}, {N: "y0", V: 0.50},
		&P{N: "dy", Extra: "-0.3  0  -0.15  -0.04  -0.01"},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "ptsB", tmin, tmax, xcte, 8)
	}

	tmin = 0.01
	tmax = 0.99
	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts06c(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts06c. pts")

	fun, err := New("pts", []*P{
		// T =                     0 0.05 0.1 0.2 0.3 0.5  0.75 1
		&P{N: "y=dt", Extra: "0.05 0.05 0.1 0.1 0.2 0.25 0.25 0"},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "ptsC", tmin, tmax, xcte, 8)
	}

	tmin = 0.01
	tmax = 0.99
	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 10, nil, sktol, dtol, dtol2, ver)
}

func Test_ts07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts07. exc1")

	fun, err := New("exc1", []*P{
		&P{N: "a", V: 200},
		&P{N: "b", V: 2},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "exc1", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-7
	dtol2 := 1e-6
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts08. exc2")

	fun, err := New("exc2", []*P{
		&P{N: "ta", V: 5},
		&P{N: "a", V: 3},
		&P{N: "b", V: 0.2},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 7.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "exc2", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts09(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts09. cos")

	fun, err := New("cos", []*P{
		&P{N: "a", V: 10},
		&P{N: "b", V: math.Pi},
		&P{N: "c", V: 1.0},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 2.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "cos", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-8
	dtol2 := 1e-7
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts10(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts10. rmp")

	fun, err := New("rmp", []*P{
		&P{N: "ta", V: 1},
		&P{N: "tb", V: 2},
		&P{N: "ca", V: 0.5},
		&P{N: "cb", V: -1.5},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 3.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "rmp", tmin, tmax, xcte, 4)
	}

	sktol := 1e-10
	dtol := 1e-12
	dtol2 := 1e-17
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts11(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts11. ref-inc-rl1")

	fun, err := New("ref-inc-rl1", []*P{
		&P{N: "lam0", V: 0.001},
		&P{N: "lam1", V: 1.2},
		&P{N: "alp", V: 0.01},
		&P{N: "bet", V: 10},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "ref-inc-rl1", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts12(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts12. mul")

	cos, err := New("cos", []*P{
		&P{N: "a", V: 1},
		&P{N: "b/pi", V: 2},
		&P{N: "c", V: 1},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	lin, err := New("lin", []*P{
		&P{N: "m", V: 0.5},
		&P{N: "ts", V: 0},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	mul, err := New("mul", []*P{
		&P{N: "fa", Fcn: cos},
		&P{N: "fb", Fcn: lin},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 1.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(cos, "/tmp/gosl/fun", "cosB", tmin, tmax, xcte, 41)
		plt.Clf()
		PlotT(lin, "/tmp/gosl/fun", "linB", tmin, tmax, xcte, 41)
		plt.Clf()
		PlotT(mul, "/tmp/gosl/fun", "mul", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-9
	dtol2 := 1e-8
	ver := chk.Verbose
	tskip := []float64{tmin, tmax}
	CheckDerivT(tst, cos, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	io.Pf("\n")
	CheckDerivT(tst, lin, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
	io.Pf("\n")
	CheckDerivT(tst, mul, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
}

func Test_ts13(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts13. pulse")

	pulse, err := New("pulse", []*P{
		&P{N: "ca", V: 0.2},
		&P{N: "cb", V: 2.0},
		&P{N: "ta", V: 1.0},
		&P{N: "tb", V: 2.5},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 5.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(pulse, "/tmp/gosl/fun", "pulse", tmin, tmax, xcte, 61)
	}

	sktol := 1e-17
	dtol := 1e-10
	dtol2 := 1e-10
	ver := chk.Verbose
	tskip := []float64{1, 4}
	CheckDerivT(tst, pulse, tmin, tmax, xcte, 11, tskip, sktol, dtol, dtol2, ver)
}

func Test_ts14(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts14. sin")

	fun, err := New("sin", []*P{
		&P{N: "a", V: 10},
		&P{N: "b", V: math.Pi},
		&P{N: "c", V: 1.0},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 2.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "sin", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-8
	dtol2 := 1e-7
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts15(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts15. cut-sin; test cut positive values.")

	fun, err := New("cut-sin", []*P{
		&P{N: "a", V: 10},
		&P{N: "b", V: math.Pi},
		&P{N: "c", V: 1.0},
		&P{N: "cps", V: 0.0},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 2.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "cut-sin-positive", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-8
	dtol2 := 1e-7
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}

func Test_ts16(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ts16. cut-sin; test cut negative values.")

	fun, err := New("cut-sin", []*P{
		&P{N: "a", V: 10},
		&P{N: "b", V: math.Pi},
		&P{N: "c", V: 1.0},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n", err)
		return
	}

	tmin := 0.0
	tmax := 2.0
	xcte := []float64{0, 0, 0}
	if chk.Verbose {
		plt.Reset(false, nil)
		PlotT(fun, "/tmp/gosl/fun", "cut-sin-negative", tmin, tmax, xcte, 41)
	}

	sktol := 1e-10
	dtol := 1e-8
	dtol2 := 1e-7
	ver := chk.Verbose
	CheckDerivT(tst, fun, tmin, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
}
