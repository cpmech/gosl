// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestAdata01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Adata01. conversions")

	var obj interface{}

	obj = FromString("abc")
	chk.String(tst, ToString(obj), "abc")

	obj = FromStrings([]string{"abc", "edf"})
	chk.Strings(tst, "ToStrings", ToStrings(obj), []string{"abc", "edf"})

	obj = FromInt(123)
	chk.Int(tst, "ToInt", ToInt(obj), 123)

	obj = FromInts([]int{123, 456})
	chk.Ints(tst, "ToInt", ToInts(obj), []int{123, 456})

	obj = FromFloat64(456.123)
	chk.Float64(tst, "ToFloat64", 1e-15, ToFloat64(obj), 456.123)

	obj = FromFloat64s([]float64{456.123, 123.456})
	chk.Array(tst, "ToFloat64s", 1e-15, ToFloat64s(obj), []float64{456.123, 123.456})

	obj = FromBool(false)
	if ToBool(obj) {
		tst.Errorf("ToBool(false) failed\n")
		return
	}

	obj = FromBool(true)
	if !ToBool(obj) {
		tst.Errorf("ToBool(true) failed\n")
		return
	}

	obj = FromBools([]bool{true, false})
	vals := ToBools(obj)
	if vals[0] != true || vals[1] != false {
		tst.Errorf("ToBools failed\n")
		return
	}
}
