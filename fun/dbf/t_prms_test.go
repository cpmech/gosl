// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_prms01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("prms01")

	var prms Params
	prms = []*P{
		&P{N: "a", V: 1.0},
		&P{N: "b", V: 2.0},
		&P{N: "c", V: 3.0},
	}
	io.Pforan("%v\n", prms)

	var a, b, c, A, B float64
	e := prms.Connect(&a, "a", "test call")
	e += prms.Connect(&b, "b", "test call")
	e += prms.Connect(&c, "c", "test call")
	e += prms.Connect(&A, "a", "test call")
	e += prms.Connect(&B, "b", "test call")
	if e != "" {
		tst.Error("connect failed\n")
		return
	}

	chk.Scalar(tst, "a", 1e-15, a, 1.0)
	chk.Scalar(tst, "b", 1e-15, b, 2.0)
	chk.Scalar(tst, "c", 1e-15, c, 3.0)
	chk.Scalar(tst, "A", 1e-15, A, 1.0)
	chk.Scalar(tst, "B", 1e-15, B, 2.0)

	prm := prms.Find("a")
	if prm == nil {
		tst.Error("cannot find parameter 'a'\n")
		return
	}
	prm.Set(123)
	chk.Scalar(tst, "a", 1e-15, a, 123)
	chk.Scalar(tst, "A", 1e-15, A, 123)
}

func Test_prms02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("prms02")

	var prms Params
	prms = []*P{
		&P{N: "klx", V: 1.0},
		&P{N: "kly", V: 2.0},
		&P{N: "klz", V: 3.0},
	}
	io.Pforan("%v\n", prms)

	var klx, kly, klz float64
	err_msg := prms.ConnectSet([]*float64{&klx, &kly, &klz}, []string{"klx", "kly", "klz"}, "Test_prms02")
	if err_msg != "" {
		tst.Error("connect set failed\n")
		return
	}

	chk.Scalar(tst, "klx", 1e-15, klx, 1.0)
	chk.Scalar(tst, "kly", 1e-15, kly, 2.0)
	chk.Scalar(tst, "klz", 1e-15, klz, 3.0)

	prms[1].Set(2.2)
	chk.Scalar(tst, "kly", 1e-15, kly, 2.2)
}

func Test_prms03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("prms03")

	var prms Params
	prms = []*P{
		&P{N: "klx", V: 1.0},
		&P{N: "kly", V: 2.0},
		&P{N: "klz", V: 3.0},
	}
	io.Pforan("%v\n", prms)

	values, found := prms.GetValues([]string{"klx", "kly", "klz"})
	chk.Vector(tst, "values", 1e-15, values, []float64{1, 2, 3})
	chk.Bools(tst, "found", found, []bool{true, true, true})

	values, found = prms.GetValues([]string{"klx", "klY", "klz"})
	chk.Vector(tst, "values", 1e-15, values, []float64{1, 0, 3})
	chk.Bools(tst, "found", found, []bool{true, false, true})
}

func Test_prms04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("prms04")

	var prms Params
	prms = []*P{
		&P{N: "a", V: 123, Min: 0, Max: math.MaxFloat64},
		&P{N: "b", V: 456, Min: 0, Max: math.MaxFloat64},
	}

	prms.CheckLimits()

	values := prms.CheckAndGetValues([]string{"a", "b"})
	chk.Vector(tst, "values", 1e-15, values, []float64{123, 456})

	var a, b float64
	prms.CheckAndSetVars([]string{"a", "b"}, []*float64{&a, &b})
	chk.Scalar(tst, "a", 1e-15, a, 123)
	chk.Scalar(tst, "b", 1e-15, b, 456)
}
