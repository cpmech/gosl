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

func TestLogReg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("LogReg01. Logistic Regression")

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

	// train linear regression
	io.Pl()
	g := NewGradDesc(20)
	g.SetControl(5.0, 0, 0)
	g.Run(d, r, []float64{5, 5, 5})
	io.Pf("grad.desc: θ = %v\n", d.thetaVec)

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})

		args := &plt.A{Colors: []string{"k"}, Levels: []float64{0.5}}
		plt.Subplot(2, 1, 1)
		d.PlotX(0, 1, map[int]*plt.A{0: {C: "b", M: "o", Void: true}, 1: {C: "r", M: "x"}})
		d.PlotContModel(r, 0, 1, 11, 11, nil, nil, false, args)
		plt.Equal()
		plt.HideTRborders()

		plt.Subplot(2, 1, 2)
		g.PlotCostIter(nil)
		plt.Gll("$iteration$", "$cost$", nil)
		plt.HideTRborders()

		plt.Save("/tmp/gosl/ml", "logreg01")
	}
}
