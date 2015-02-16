// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"testing"

	"github.com/cpmech/gosl/utl"
)

func Test_draw01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("draw01")

	P := [][]float64{
		{-2.5, 0.0},
		{-5.5, 4.0},
		{0.0, 10.0},
		{5.5, 4.0},
		{2.5, 0.0},
	}

	var sd ShapeData
	sd.Init()
	sd.Closed = true
	DrawPolyline(P, &sd, "")
	AutoScale(P)
	Equal()
	DrawLegend([]LineData{
		LineData{"first", "red", 1, -1, "o", "-"},
		LineData{"second", "green", 2, 0, "s", "-"},
		LineData{"third", "blue", 3, 10, "+", "-"},
	}, 10, "best", false, "")
	//Show()
}

func Test_draw02(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("draw02")

	d := LineData{"gofem", "red", 1.2, 10, "o", "--"}
	l := d.GetArgs("clip_on=0")
	utl.Pforan("l = %q\n", l)
	utl.CheckString(tst, l, "clip_on=0,label='gofem',color='red',lw=1.2,ms=10,marker='o',ls='--'")
}
