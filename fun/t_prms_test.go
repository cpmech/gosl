// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_prms01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("prms01")

	var prms Prms
	prms = []*Prm{
		&Prm{N: "a", V: 1.0},
		&Prm{N: "b", V: 2.0},
		&Prm{N: "c", V: 3.0},
	}
	io.Pforan("%v\n", prms)

	var a, b, c, A, B float64
	e := prms.Connect(&a, "a", "test call")
	e += prms.Connect(&b, "b", "test call")
	e += prms.Connect(&c, "c", "test call")
	e += prms.Connect(&A, "a", "test call")
	e += prms.Connect(&B, "b", "test call")
	if e != "" {
		tst.Error("connect failed: %v\n", e)
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

	var prms Prms
	prms = []*Prm{
		&Prm{N: "klx", V: 1.0},
		&Prm{N: "kly", V: 2.0},
		&Prm{N: "klz", V: 3.0},
	}
	io.Pforan("%v\n", prms)

	var klx, kly, klz float64
	err_msg := prms.ConnectSet([]*float64{&klx, &kly, &klz}, []string{"klx", "kly", "klz"}, "Test_prms02")
	if err_msg != "" {
		tst.Error("connect set failed: %v\n", err_msg)
		return
	}

	chk.Scalar(tst, "klx", 1e-15, klx, 1.0)
	chk.Scalar(tst, "kly", 1e-15, kly, 2.0)
	chk.Scalar(tst, "klz", 1e-15, klz, 3.0)

	prms[1].Set(2.2)
	chk.Scalar(tst, "kly", 1e-15, kly, 2.2)
}
