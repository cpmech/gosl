// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestBasic01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Basic01")

	if !Atob("1") {
		tst.Errorf("Atob(\"1\") should have returned true\n")
		return
	}
	if !Atob("true") {
		tst.Errorf("Atob(\"true\") should have returned true\n")
	}
	if Atob("0") {
		tst.Errorf("Atob(\"0\") should have returned false\n")
	}
	if Atob("false") {
		tst.Errorf("Atob(\"false\") should have returned false\n")
	}

	if Itob(0) {
		tst.Errorf("Itob(0) should have returned false\n")
	}
	if !Itob(-1) {
		tst.Errorf("Itob(-1) should have returned true\n")
	}
	if !Itob(+1) {
		tst.Errorf("Itob(+1) should have returned true\n")
	}

	chk.Int(tst, "true => 1", Btoi(true), 1)
	chk.Int(tst, "false => 0", Btoi(false), 0)

	chk.Int(tst, "\"123\" => 123", Atoi("123"), 123)

	chk.String(tst, Btoa(true), "true")
	chk.String(tst, Btoa(false), "false")

	chk.Float64(tst, "\"123.456\" => 123.456", 1e-15, Atof("123.456"), 123.456)
}

func TestBasic02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Parsing02. Atob panic")

	defer chk.RecoverTstPanicIsOK(tst)
	Atob("dorival")
}

func TestBasic03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Parsing03. Atoi panic")

	defer chk.RecoverTstPanicIsOK(tst)
	Atoi("dorival")
}

func TestBasic04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Parsing04. Atof panic")

	defer chk.RecoverTstPanicIsOK(tst)
	Atof("dorival")
}

func TestBasic05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Basic05. IntSf, DblSf, StrSf")

	res := IntSf("%4d", []int{1, 2, 3}) // note that an inner space is always added
	Pforan("res = %q\n", res)
	chk.String(tst, res, "   1    2    3")

	res = DblSf("%8.3f", []float64{1, 2, 3}) // note that an inner space is always added
	Pforan("res = %q\n", res)
	chk.String(tst, res, "   1.000    2.000    3.000")

	res = StrSf("%s", []string{"a", "b", "c"}) // note that an inner space is always added
	Pforan("res = %q\n", res)
	chk.String(tst, res, "a b c")
}
