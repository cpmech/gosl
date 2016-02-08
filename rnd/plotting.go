// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// PlotPdf plots PDF
func (o VarData) PlotPdf(np int, args string) {
	X := utl.LinSpace(o.Min, o.Max, np)
	Y := make([]float64, np)
	for i := 0; i < np; i++ {
		Y[i] = o.Distr.Pdf(X[i])
	}
	if args == "" {
		args = "'b-'"
	}
	plt.Plot(X, Y, args)
	plt.Gll("$x$", "$f(x)$", "")
}
