// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestMetrics01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Metrics01")

	// new object
	rin, rou := 2.0, 6.0 // radii
	trf := FactoryTfinite.Surf2dQuarterRing(rin, rou)

	rvals := utl.NonlinSpace(-1, 1, 5, 4, false)
	svals := utl.NonlinSpace(-1, 1, 7, 3, true)
	m := trf.GetMetrics2d(rvals, svals)

	lr := len(rvals) - 1
	ls := len(svals) - 1
	ms := ls / 2
	s2 := 1.0 / math.Sqrt2
	chk.Array(tst, "B0:Nf(r=0.0)", 1e-15, m.Nf[0][0], []float64{0, -1})
	chk.Array(tst, "B0:Nf(r=1.0)", 1e-15, m.Nf[0][lr], []float64{0, -1})
	io.Pl()
	chk.Array(tst, "B1:Nf(s=0.0)", 1e-15, m.Nf[1][0], []float64{1, 0})
	chk.Array(tst, "B1:Nf(s=0.5)", 1e-15, m.Nf[1][ms], []float64{s2, s2})
	chk.Array(tst, "B1:Nf(s=1.0)", 1e-15, m.Nf[1][ls], []float64{0, 1})
	io.Pl()
	chk.Array(tst, "B2:Nf(r=0.0)", 1e-15, m.Nf[2][0], []float64{-1, 0})
	chk.Array(tst, "B2:Nf(r=1.0)", 1e-15, m.Nf[2][lr], []float64{-1, 0})
	io.Pl()
	chk.Array(tst, "B3:Nf(s=0.0)", 1e-15, m.Nf[3][0], []float64{-1, 0})
	chk.Array(tst, "B3:Nf(s=0.5)", 1e-15, m.Nf[3][ms], []float64{-s2, -s2})
	chk.Array(tst, "B3:Nf(s=1.0)", 1e-15, m.Nf[3][ls], []float64{0, -1})

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{31, 31}, true, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 1, NoClip: true})
		m.Draw(0.3, nil, nil)
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "metrics01")
	}
}
