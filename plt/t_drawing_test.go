// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
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

	var sd FmtH
	sd.Init()
	sd.Closed = true
	DrawPolyline(P, &sd, "")
	AutoScale(P)
	Equal()
	DrawLegend([]FmtL{
		FmtL{"first", "red", 1, -1, "o", "-"},
		FmtL{"second", "green", 2, 0, "s", "-"},
		FmtL{"third", "blue", 3, 10, "+", "-"},
	}, 10, "best", false, "")
	//Show()
}

func Test_draw02(tst *testing.T) {

	chk.PrintTitle("draw02")

	d := FmtL{"gofem", "red", 1.2, 10, "o", "--"}
	l := d.GetArgs("clip_on=0")
	io.Pforan("l = %q\n", l)
	chk.String(tst, l, "clip_on=0,label='gofem',color='red',lw=1.2,ms=10,marker='o',ls='--'")
}
