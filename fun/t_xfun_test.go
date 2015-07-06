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

	if false {
		//if true {
		tcte := 0.0
		xmin := []float64{-1, -1}
		xmax := []float64{2, 2}
		np := 21
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

	/*
		sktol := 1e-10
		dtol := 1e-10
		dtol2 := 1e-10
		ver := chk.Verbose
		CheckT(tst, o, 0.0, tmax, xcte, 11, nil, sktol, dtol, dtol2, ver)
	*/
}
