// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func PlotNurbs(dirout, fnkey string, b *Nurbs) {
	plt.SetForEps(0.75, 400)
	b.DrawCtrl2D(true)
	b.DrawElems2D(21, true, "", "")
	plt.Equal()
	plt.SaveD(dirout, fnkey+".eps")
}

func PlotNurbsBasis(dirout, fnkey string, b *Nurbs, la, lb int) {
	npts := 41
	plt.SetForEps(1.2, 500)

	plt.Subplot(3, 2, 1)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	t0 := time.Now()
	b.PlotBasis(la, "", 11, 0) // 0 => CalcBasis
	io.Pfcyan("time elapsed (calcbasis) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(3, 2, 2)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 0) // 0 => CalcBasis
	plt.Equal()

	plt.Subplot(3, 2, 3)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(la, "", 11, 1) // 1 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(3, 2, 4)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 1) // 1 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(3, 2, 5)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	t0 = time.Now()
	b.PlotBasis(la, "", 11, 2) // 2 => RecursiveBasis
	io.Pfcyan("time elapsed (recursive) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(3, 2, 6)
	b.DrawCtrl2D(false)
	b.DrawElems2D(npts, false, "", "")
	b.PlotBasis(lb, "", 11, 2) // 2 => RecursiveBasis
	plt.Equal()

	plt.SaveD(dirout, fnkey+".eps")
}

func PlotNurbsDerivs(dirout, fnkey string, b *Nurbs, la, lb int) {
	npts := 41
	plt.SetForEps(1.5, 500)

	plt.Subplot(4, 2, 1)
	t0 := time.Now()
	b.PlotDeriv(la, 0, "", npts, 0) // 0 => CalcBasisAndDerivs
	io.Pfcyan("time elapsed (calcbasis) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(4, 2, 2)
	t0 = time.Now()
	b.PlotDeriv(la, 0, "", npts, 1) // 1 => NumericalDeriv
	io.Pfcyan("time elapsed (numerical) = %v\n", time.Now().Sub(t0))
	plt.Equal()

	plt.Subplot(4, 2, 3)
	b.PlotDeriv(la, 1, "", npts, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 4)
	b.PlotDeriv(la, 1, "", npts, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.Subplot(4, 2, 5)
	b.PlotDeriv(lb, 0, "", npts, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 6)
	b.PlotDeriv(lb, 0, "", npts, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.Subplot(4, 2, 7)
	b.PlotDeriv(lb, 1, "", npts, 0) // 0 => CalcBasisAndDerivs
	plt.Equal()

	plt.Subplot(4, 2, 8)
	b.PlotDeriv(lb, 1, "", npts, 1) // 0 => NumericalDeriv
	plt.Equal()

	plt.SaveD(dirout, fnkey+".eps")
}

func PlotNurbsRefined(dirout, fnkey string, b, c *Nurbs) {
	plt.SetForEps(1.5, 400)
	plt.Subplot(3, 1, 1)
	b.DrawCtrl2D(true)
	b.DrawElems2D(21, true, "", "")
	plt.Equal()

	plt.Subplot(3, 1, 2)
	c.DrawCtrl2D(true)
	c.DrawElems2D(21, true, "", "")
	plt.Equal()

	plt.Subplot(3, 1, 3)
	b.DrawElems2D(21, true, ", lw=3", "")
	c.DrawElems2D(21, true, ", color='red', marker='+', markevery=10", "color='magenta', size=8, va='bottom'")
	plt.Equal()

	plt.SaveD(dirout, fnkey+".eps")
}
