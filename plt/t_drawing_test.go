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
		Fmt{C: "red", M: "o", Ls: "-", Lw: 1, Ms: -1, L: "first", Me: -1},
		Fmt{C: "green", M: "s", Ls: "-", Lw: 2, Ms: 0, L: "second", Me: -1},
		Fmt{C: "blue", M: "+", Ls: "-", Lw: 3, Ms: 10, L: "third", Me: -1},
	}, 10, "best", false, "")
	if chk.Verbose {
		SaveD("/tmp/gosl", "draw01.eps")
	}
}

func Test_draw02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("draw02")

	d := Fmt{"red", "o", "--", 1.2, -1, "gofem", 2, 10, "blue", 0.3, true, true}
	l := d.GetArgs("clip_on=0")
	io.Pforan("l = %q\n", l)
	chk.String(tst, l, "clip_on=0,color='red',marker='o',ls='--',lw=1.2,label='gofem',markevery=2,zorder=10,markeredgecolor='blue',mew=0.3,markerfacecolor='none',clip_on=1")
}
