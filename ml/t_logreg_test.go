// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

func TestLogReg00(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg00. Logistic Regression")

	// data
	nr, nc := 5, 5
	d := NewDataMatrix(nr*nc, 2, true)
	k := 0
	for j := 0; j < nc; j++ {
		for i := 0; i < nr; i++ {
			x0 := -1 + 2*float64(i)/float64(nc-1)
			x1 := -1 + 2*float64(j)/float64(nr-1)
			d.SetX(k, 0, x0)
			d.SetX(k, 1, x1)
			if x0+x1+1e-15 >= 0 {
				d.SetY(k, 1)
			}
			k++
		}
	}

	// logistic regression
	r := NewLogReg(d)

	// check dCdθ
	io.Pl()
	dCdθ := la.NewVector(d.Nparams())
	for _, θ0 := range []float64{0.5, 1.0} {
		for _, θ1 := range []float64{0.5, 1.0} {
			r.θ[0] = θ0
			r.θ[1] = θ1
			r.Deriv(dCdθ, d)
			chk.DerivScaVec(tst, "dCdθ", 1e-7, dCdθ, r.θ, 1e-6, chk.Verbose, func(th []float64) float64 {
				copy(r.θ, th)
				return r.Cost(d)
			})
		}
	}

	// Gradient descent
	g := NewGradDesc(10)
	g.SetControl(100.0, 0, 0)
	g.Run(d, r, []float64{0, 0, 0}, 0)
	C := r.Cost(d)
	io.Pl()
	io.Pf("GradDesc: θ = %.8f\n", r.θ)
	io.Pf("cost(θ) = %.8f\n", C)

	// train linear regression using Newton's method
	if true {
		checkJac := true
		tolJac0 := 1e-5
		tolJac1 := 1e-5
		io.Pl()
		r.θ.Fill(0)
		r.CalcTheta(d, chk.Verbose, checkJac, tolJac0, tolJac1, map[string]float64{"ftol": 0.0005})
		C = r.Cost(d)
		io.Pl()
		io.Pf("Newton: θ = %.8f\n", r.θ)
		io.Pf("cost(θ) = %.8f\n", C)
	}

	// plot
	if chk.Verbose {
		io.Pl()
		plt.Reset(true, &plt.A{WidthPt: 350, Dpi: 150, Prop: 1.7})

		argsContour := &plt.A{Colors: []string{plt.C(1, 0)}, Levels: []float64{0.5}}
		plt.Subplot(2, 1, 1)
		d.PlotX(0, 1, map[int]*plt.A{0: {C: "b", M: "o", Void: true}, 1: {C: "r", M: "x"}})
		d.PlotContModel(r, 0, 1, 11, 11, nil, nil, nil, false, argsContour)
		plt.Equal()
		plt.HideTRborders()

		plt.Subplot(2, 1, 2)
		g.PlotCostIter(nil)
		plt.Gll("$iteration$", "$cost$", nil)
		plt.HideTRborders()

		plt.Save("/tmp/gosl/ml", "logreg00")
	}
}

func TestLogReg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg01. Logistic Reg: ANG test 1")

	// data
	XYraw := io.ReadMatrix("./data/angEx2data1.txt")
	d := NewDataMatrixTable(XYraw)
	io.Pl()
	chk.Int(tst, "nData", d.Nsamples(), 100)
	chk.Int(tst, "nFeatures", d.Nfeatures(), 2)

	// initial theta
	θ := []float64{-24, 0.2, 0.2}

	// model
	r := NewLogReg(d)
	r.θ.Apply(1, θ)
	C := r.Cost(d)
	io.Pl()
	io.Pf("Initial: θ = %.8f\n", r.θ)
	io.Pf("cost(θini) = %.8f\n", C)
	chk.Float64(tst, "cost(θini)", 1e-15, C, 2.183301938265978e-01)

	// Gradient descent
	g := NewGradDesc(10)
	g.SetControl(0.002, 0, 0)
	g.Run(d, r, r.θ, r.b)
	C = r.Cost(d)
	io.Pl()
	io.Pf("GradDesc: θ = %.8f\n", r.θ)
	io.Pf("cost(θ) = %.8f\n", C)
	chk.Float64(tst, "cost(θgrad)", 1e-15, C, 2.037591668976244e-01)
	chk.Array(tst, "θ", 1e-14, r.θ, []float64{-2.400009669708430e+01, 1.957478716620902e-01, 1.933175159514175e-01})

	// Newton's method
	checkJac := true
	tolJac0 := 1e-3
	tolJac1 := 1e-4
	useZeroGuess := true
	if useZeroGuess {
		r.θ.Fill(0) // start with challenging initial θ
	} else {
		r.θ.Apply(1, θ)
	}
	io.Pl()
	r.CalcTheta(d, chk.Verbose, checkJac, tolJac0, tolJac1, nil)
	C = r.Cost(d)
	io.Pl()
	io.Pf("Newton: θ = %.8f\n", r.θ)
	io.Pf("cost(θ) = %.8f\n", C)
	chk.Float64(tst, "cost(θnewt)", 1e-15, C, 2.034977015894404e-01)
	chk.Array(tst, "θ", 1e-14, r.θ, []float64{-2.516133256589910e+01, 2.062317052577260e-01, 2.014715922708144e-01})

	// plot
	if chk.Verbose {
		io.Pl()
		plt.Reset(true, &plt.A{WidthPt: 350, Dpi: 150, Prop: 1.7})

		argsContour := &plt.A{Colors: []string{plt.C(1, 0)}, Levels: []float64{0.5}}
		plt.Subplot(2, 1, 1)
		d.PlotX(0, 1, map[int]*plt.A{0: {C: "b", M: "o", Void: true}, 1: {C: "r", M: "x"}})
		d.PlotContModel(r, 0, 1, 41, 41, nil, nil, nil, false, argsContour)

		plt.Subplot(2, 1, 2)
		g.PlotCostIter(nil)
		plt.Gll("$iteration$", "$cost$", nil)
		plt.SetTicksNormal()
		plt.HideTRborders()

		plt.Save("/tmp/gosl/ml", "logreg01")
	}
}

func TestLogReg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg02. Logistic Reg: ANG test 2")

	// raw data
	XYraw := io.ReadMatrix("./data/angEx2data2.txt")

	// mapped data
	nOriFeatures := len(XYraw[0]) - 1 // -1 because y-column
	mapper := NewPolyDataMapper(nOriFeatures, 0, 1, 6)
	d := mapper.GetMapped(XYraw, true)
	chk.Int(tst, "nData", d.Nsamples(), 118)
	chk.Int(tst, "nFeatures", d.Nfeatures(), 27)

	// initial theta
	θ := la.NewVector(d.Nparams())
	θ.Fill(1)

	// model
	λ := 1.0
	r := NewLogReg(d)
	r.SetRegularization(λ)
	r.θ.Apply(1, θ)
	C := r.Cost(d)
	io.Pl()
	io.Pf("Initial: θ[:5] = %.8f\n", r.θ[:5])
	io.Pf("cost(θini) = %.8f\n", C)
	chk.Float64(tst, "cost(θini)", 1e-15, C, 2.134848314666066)

	// check dCdθ
	io.Pl()
	dCdθ := la.NewVector(d.Nparams())
	verb := false
	for _, θ0 := range []float64{0.5, 1.0} {
		for _, θ1 := range []float64{0.5, 1.0} {
			r.θ[0] = θ0
			r.θ[1] = θ1
			r.Deriv(dCdθ, d)
			chk.DerivScaVec(tst, "dCdθ", 1e-7, dCdθ, r.θ, 1e-6, verb, func(th []float64) float64 {
				copy(r.θ, th)
				return r.Cost(d)
			})
		}
	}

	// Gradient descent
	g := NewGradDesc(10)
	g.SetControl(5.0, 0, 0)
	r.θ.Apply(1, θ)
	g.Run(d, r, r.θ, r.b)
	C = r.Cost(d)
	io.Pl()
	io.Pf("GradDesc: θ[:5] = %.8f\n", r.θ[:5])
	io.Pf("cost(θ) = %.8f\n", C)
	chk.Float64(tst, "cost(θgrad)", 1e-15, C, 5.920108560779025e-01)
	chk.Array(tst, "θ[:5]", 1e-14, r.θ[:5], []float64{6.527575848138054e-01, -1.730594903181217e-01, 3.615618466861891e-01, -1.194645899263627e+00, -4.186288373383852e-01})

	// Newton's method
	checkJac := true
	tolJac0 := 1e-3
	tolJac1 := 1e-4
	useZeroGuess := true
	if useZeroGuess {
		r.θ.Fill(0)
	} else {
		r.θ.Apply(1, θ)
	}
	io.Pl()
	r.CalcTheta(d, chk.Verbose, checkJac, tolJac0, tolJac1, nil)
	C = r.Cost(d)
	io.Pl()
	io.Pf("Newton: θ[:5] = %.8f\n", r.θ[:5])
	io.Pf("cost(θ) = %.15e\n", C)
	chk.Float64(tst, "cost(θnewt)", 1e-15, C, 5.290027411158117e-01)
	chk.Array(tst, "θ[:5]", 1e-14, r.θ[:5], []float64{1.272656700281225, 6.252526148274546e-01, 1.180976145721166, -2.019842398401904, -9.173659359499787e-01})

	// plot
	if chk.Verbose {
		io.Pl()
		plt.Reset(true, &plt.A{WidthPt: 350, Dpi: 150, Prop: 1.7})

		argsContour := &plt.A{Colors: []string{plt.C(1, 0)}, Levels: []float64{0.5}}
		plt.Subplot(2, 1, 1)
		d.PlotX(0, 1, map[int]*plt.A{0: {C: "b", M: "o", Void: true}, 1: {C: "r", M: "x"}})
		d.PlotContModel(r, 0, 1, 41, 41, mapper, nil, nil, false, argsContour)

		plt.Subplot(2, 1, 2)
		g.PlotCostIter(nil)
		plt.Gll("$iteration$", "$cost$", nil)
		plt.SetTicksNormal()
		plt.HideTRborders()

		plt.Save("/tmp/gosl/ml", "logreg02")
	}
}
