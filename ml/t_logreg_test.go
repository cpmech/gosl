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
	d := NewRegData(nr*nc, 2)
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
			d.thetaVec[0] = θ0
			d.thetaVec[1] = θ1
			r.Deriv(dCdθ, d)
			chk.DerivScaVec(tst, "dCdθ", 1e-7, dCdθ, d.thetaVec, 1e-6, chk.Verbose, func(th []float64) float64 {
				copy(d.thetaVec, th)
				return r.Cost(d)
			})
		}
	}

	// Gradient descent
	g := NewGradDesc(10)
	g.SetControl(100.0, 0, 0)
	g.Run(d, r, []float64{0, 0, 0})
	C := r.Cost(d)
	io.Pl()
	io.Pf("GradDesc: θ = %.8f\n", d.thetaVec)
	io.Pf("cost(θ) = %.8f\n", C)

	// train linear regression using Newton's method
	if true {
		checkJac := true
		tolJac0 := 1e-5
		tolJac1 := 1e-5
		io.Pl()
		d.thetaVec.Fill(0)
		r.CalcTheta(d, chk.Verbose, checkJac, tolJac0, tolJac1, map[string]float64{"ftol": 0.0005})
		C = r.Cost(d)
		io.Pl()
		io.Pf("Newton: θ = %.8f\n", d.thetaVec)
		io.Pf("cost(θ) = %.8f\n", C)
	}

	// plot
	if chk.Verbose {
		io.Pl()
		plt.Reset(true, &plt.A{WidthPt: 350, Dpi: 150, Prop: 1.7})

		argsContour := &plt.A{Colors: []string{plt.C(1, 0)}, Levels: []float64{0.5}}
		plt.Subplot(2, 1, 1)
		d.PlotX(0, 1, map[int]*plt.A{0: {C: "b", M: "o", Void: true}, 1: {C: "r", M: "x"}})
		d.PlotContModel(r, 0, 1, 11, 11, nil, nil, false, argsContour)
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
	Xraw := io.ReadMatrix("./data/angEx2data1.txt")
	d := NewRegDataTable(Xraw)
	io.Pl()
	chk.Int(tst, "nData", d.Ndata(), 100)
	chk.Int(tst, "nFeatures", d.Nfeatures(), 2)

	// initial theta
	θ := []float64{-24, 0.2, 0.2}

	// model
	r := NewLogReg(d)
	d.thetaVec.Apply(1, θ)
	C := r.Cost(d)
	io.Pl()
	io.Pf("Initial: θ = %.8f\n", d.thetaVec)
	io.Pf("cost(θini) = %.8f\n", C)
	chk.Float64(tst, "cost(θini)", 1e-15, C, 2.183301938265978e-01)

	// Gradient descent
	g := NewGradDesc(10)
	g.SetControl(0.002, 0, 0)
	g.Run(d, r, d.thetaVec)
	C = r.Cost(d)
	io.Pl()
	io.Pf("GradDesc: θ = %.8f\n", d.thetaVec)
	io.Pf("cost(θ) = %.8f\n", C)
	chk.Float64(tst, "cost(θgrad)", 1e-15, C, 2.037591668976244e-01)
	chk.Array(tst, "θ", 1e-14, d.thetaVec, []float64{-2.400009669708430e+01, 1.957478716620902e-01, 1.933175159514175e-01})

	// Newton's method
	checkJac := true
	tolJac0 := 1e-3
	tolJac1 := 1e-4
	useZeroGuess := true
	if useZeroGuess {
		d.thetaVec.Fill(0) // start with challenging initial θ
	} else {
		d.thetaVec.Apply(1, θ)
	}
	io.Pl()
	r.CalcTheta(d, chk.Verbose, checkJac, tolJac0, tolJac1, nil)
	C = r.Cost(d)
	io.Pl()
	io.Pf("Newton: θ = %.8f\n", d.thetaVec)
	io.Pf("cost(θ) = %.8f\n", C)
	chk.Float64(tst, "cost(θnewt)", 1e-15, C, 2.034977015894404e-01)
	chk.Array(tst, "θ", 1e-14, d.thetaVec, []float64{-2.516133256589910e+01, 2.062317052577260e-01, 2.014715922708144e-01})

	// plot
	if chk.Verbose {
		io.Pl()
		plt.Reset(true, &plt.A{WidthPt: 350, Dpi: 150, Prop: 1.7})

		argsContour := &plt.A{Colors: []string{plt.C(1, 0)}, Levels: []float64{0.5}}
		plt.Subplot(2, 1, 1)
		d.PlotX(0, 1, map[int]*plt.A{0: {C: "b", M: "o", Void: true}, 1: {C: "r", M: "x"}})
		d.PlotContModel(r, 0, 1, 41, 41, nil, nil, false, argsContour)

		plt.Subplot(2, 1, 2)
		g.PlotCostIter(nil)
		plt.Gll("$iteration$", "$cost$", nil)
		plt.SetTicksNormal()
		plt.HideTRborders()

		plt.Save("/tmp/gosl/ml", "logreg01")
	}
}
