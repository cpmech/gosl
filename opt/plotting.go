// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// TwoVarsFunc_t defines a function to plot contours (len(x)==2)
type TwoVarsFunc_t func(x []float64) float64

// PlotTwoVarsContour plots contour for two variables problem. len(x) == 2
//  Input
//   dirout  -- directory to save files
//   fnkey   -- file name key for eps figure
//   x       -- solution. can be <nil>
//   np      -- number of points for contour
//   extra   -- called just before saving figure
//   axequal -- axis.equal
//   vmin    -- min 0 values
//   vmax    -- max 1 values
//   f       -- function to plot filled contour. can be <nil>
//   gs      -- functions to plot contour @ level 0. can be <nil>
func PlotTwoVarsContour(dirout, fnkey string, x []float64, np int, extra func(), axequal bool,
	vmin, vmax []float64, f TwoVarsFunc_t, gs ...TwoVarsFunc_t) {
	if fnkey == "" {
		return
	}
	chk.IntAssert(len(vmin), 2)
	chk.IntAssert(len(vmax), 2)
	V0, V1 := utl.MeshGrid2D(vmin[0], vmax[0], vmin[1], vmax[1], np, np)
	var Zf [][]float64
	var Zg [][][]float64
	if f != nil {
		Zf = la.MatAlloc(np, np)
	}
	if len(gs) > 0 {
		Zg = utl.Deep3alloc(len(gs), np, np)
	}
	xtmp := make([]float64, 2)
	for i := 0; i < np; i++ {
		for j := 0; j < np; j++ {
			xtmp[0], xtmp[1] = V0[i][j], V1[i][j]
			if f != nil {
				Zf[i][j] = f(xtmp)
			}
			for k, g := range gs {
				Zg[k][i][j] = g(xtmp)
			}
		}
	}
	plt.Reset()
	plt.SetForEps(0.8, 350)
	if f != nil {
		cmapidx := 0
		plt.Contour(V0, V1, Zf, io.Sf("fsz=7, cmapidx=%d", cmapidx))
	}
	for k, _ := range gs {
		plt.ContourSimple(V0, V1, Zg[k], "zorder=5, levels=[0], colors=['yellow'], linewidths=[2], clip_on=0")
	}
	if x != nil {
		plt.PlotOne(x[0], x[1], "'r*', label='optimum', zorder=10")
	}
	if extra != nil {
		extra()
	}
	if dirout == "" {
		dirout = "."
	}
	plt.Cross("clr='grey'")
	plt.SetXnticks(11)
	plt.SetYnticks(11)
	if axequal {
		plt.Equal()
	}
	plt.AxisRange(vmin[0], vmax[0], vmin[1], vmax[1])
	args := "leg_out='1', leg_ncol=4, leg_hlen=1.5"
	plt.Gll("$x_0$", "$x_1$", args)
	plt.SaveD(dirout, fnkey+".eps")
}
