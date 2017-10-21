// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestLogReg01(tst *testing.T) {

	chk.Verbose = true
	chk.PrintTitle("LogReg01. Logistic Regression function")

	cost := func(x, y, θ float64) float64 {
		hθ := 1 / (1 + math.Exp(-θ*x))
		return -y*math.Log(hθ) - (1-y)*math.Log(1-hθ)
	}

	θ := 10.5
	xx := utl.LinSpace(-1, 1, 101)
	y1 := utl.GetMapped(xx, func(x float64) float64 { return cost(x, 0, θ) })
	y2 := utl.GetMapped(xx, func(x float64) float64 { return cost(x, 1, θ) })
	plt.Reset(true, nil)
	plt.Plot(xx, y1, &plt.A{C: plt.C(0, 0), NoClip: true})
	plt.Plot(xx, y2, &plt.A{C: plt.C(1, 0), NoClip: true})
	plt.Gll("$x$", "$cost$", nil)
	plt.HideTRborders()
	plt.Save("/tmp/gosl/ml", "logreg01")
}
