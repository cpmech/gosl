// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// PlotPdf plots PDF
func (o Variable) PlotPdf(np int, args *plt.A) {
	X := utl.LinSpace(o.Min, o.Max, np)
	Y := make([]float64, np)
	for i := 0; i < np; i++ {
		Y[i] = o.Distr.Pdf(X[i])
	}
	if args == nil {
		args = &plt.A{C: "b"}
	}
	plt.Plot(X, Y, args)
	plt.Gll("$x$", "$f(x)$", nil)
}
