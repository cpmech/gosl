// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func Test_plot01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose() = false
	chk.PrintTitle("plot01")

	x := utl.LinSpace(0.0, 1.0, 11)
	y := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		y[i] = x[i] * x[i]
	}
	Plot(x, y, "'ro', ls='-', lw=2, clip_on=0")
	Gll(`$\varepsilon$`, `$\sigma$`, "")
	//Save("/tmp/gosl", "t_plot01.eps")
}
