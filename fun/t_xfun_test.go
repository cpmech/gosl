// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func Test_xfun01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("xfun01. 2D halo => circle")

	o, err := New("halo", []*Prm{
		&Prm{N: "r", V: 0.5},
		&Prm{N: "xc", V: 0.5},
		&Prm{N: "yc", V: 0.5},
	})
	if err != nil {
		tst.Errorf("test failed: %v\n")
		return
	}

	tcte := 0.0
	xmin := []float64{-1, -1}
	xmax := []float64{2, 2}
	np := 21
	if false {
		//if true {
		withGrad := true
		hlZero := true
		axEqual := true
		save := true
		show := false
		plt.Reset()
		PlotX(o, "/tmp/gosl", "xfun01.png", tcte, xmin, xmax, np, "", withGrad, hlZero, axEqual, save, show, func() {
			plt.Equal()
		})
	}

	np = 4
	sktol := 1e-10
	dtol := 1e-10
	ver := chk.Verbose
	CheckDerivX(tst, o, tcte, xmin, xmax, np, nil, sktol, dtol, ver)
}
