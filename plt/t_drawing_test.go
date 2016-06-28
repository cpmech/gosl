// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
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

	var sd Sty
	sd.Init()
	sd.Closed = true
	DrawPolyline(P, &sd, "")
	AutoScale(P)
	Equal()
	DrawLegend([]Fmt{
		Fmt{"red", "o", "-", 1, -1, "first", -1},
		Fmt{"green", "s", "-", 2, 0, "second", -1},
		Fmt{"blue", "+", "-", 3, 10, "third", -1},
	}, 10, "best", false, "")
	if chk.Verbose {
		SaveD("/tmp/gosl", "draw01.eps")
	}
}

func Test_draw02(tst *testing.T) {

	chk.PrintTitle("draw02")

	d := Fmt{"red", "o", "--", 1.2, -1, "gofem", 2}
	l := d.GetArgs("clip_on=0")
	io.Pforan("l = %q\n", l)
	chk.String(tst, l, "clip_on=0,color='red',marker='o',ls='--',lw=1.2,label='gofem',markevery=2")
}
