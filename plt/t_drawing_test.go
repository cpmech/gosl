// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_draw01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("draw01")

	P := [][]float64{
		{-2.5, 0.0},
		{-5.5, 4.0},
		{0.0, 10.0},
		{5.5, 4.0},
		{2.5, 0.0},
	}

	Reset(true, nil)
	Polyline(P, &A{Fc: "#c1d7cf", Ec: "#4db38e", Lw: 4.5, Closed: true, NoClip: true})
	Circle(0, 4, 2.0, &A{Fc: "#b2cfa5", Ec: "#5dba35", Z: 1})
	Arrow(-4, 2, 4, 7, &A{Fc: "cyan", Ec: "blue", Z: 2, Scale: 50, Style: "fancy"})
	Arc(0, 4, 3, 0, 90, nil)
	AutoScale(P)
	Equal()

	if chk.Verbose {
		err := Save("/tmp/gosl", "t_draw01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
