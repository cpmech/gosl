// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"bytes"
	"testing"
)

func Test_genarrays01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("genarrays01")

	u := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 11, 15, 16, 38, 39}
	v := []float64{1, 2, 3, -1, -2, 0, 8, -3}

	var b0 bytes.Buffer
	Ff(&b0, "from gosl import *\n")

	var b1, b2 bytes.Buffer
	GenArray(&b1, "u", u)
	Gen2Arrays(&b2, "u", "v", u, v)
	Pfcyan("just u:\n%v\n\n", b1.String())
	Pfcyan("u and v:\n%v\n", b2.String())

	var b3 bytes.Buffer
	GenFuncT(&b3, "t", "f", 0.0, 1.0, 0.1, func(t float64) float64 { return t * t })
	Pfcyan("t and F(t):\n%v\n", b3.String())
	Ff(&b3, "plot(t,f)\nshow()\n")

	WriteFileD("/tmp/gosl", "plot_genscript_test.py", &b0, &b1, &b2, &b3)
}

func Test_genarrays02(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("genarrays02")

	a := []float64{1, 2, 3}
	b := []float64{1, 2, 3, 4}
	c := []float64{1, 2, 3, 4, 5}
	d := []float64{1, 2, 3, 4, 5, 6}
	e := []float64{1, 2, 3, 4, 5, 6, 7}
	f := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	g := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	h := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	var buf bytes.Buffer
	Ff(&buf, "from gosl import *\n")
	Gen8Arrays(&buf, "a", "b", "c", "d", "e", "f", "g", "h", a, b, c, d, e, f, g, h)
	WriteFileD("/tmp/gosl", "genarrays02.py", &buf)
}

func Test_genmat01(tst *testing.T) {

	prevTs := Tsilent
	defer func() {
		Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//Tsilent = false
	TTitle("genmat01")

	var buf bytes.Buffer
	m := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}
	GenMat(&buf, "m", m)
	WriteFileD("/tmp/gosl", "genmat01.py", &buf)
}
