// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func Test_plot01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot01")

	x := utl.LinSpace(0.0, 1.0, 11)
	y := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		y[i] = x[i] * x[i]
	}
	Plot(x, y, "'ro', ls='-', lw=2, clip_on=0")
	Gll(`$\varepsilon$`, `$\sigma$`, "")
	if chk.Verbose {
		SaveD("/tmp/gosl", "t_plot01.eps")
	}
}
