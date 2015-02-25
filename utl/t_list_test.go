// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"reflect"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_list01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("list01. DblList.Append")

	var m DblList
	m.Append(2, 2.0)
	io.Pforan("m = %v\n", m)
	equal := reflect.DeepEqual(m.Vals, [][]float64{{}, {}, {2}})
	if !equal {
		chk.PrintFail("DblList Append")
	}

	m.Append(0, 0.0)
	m.Append(1, 1.0)
	io.Pforan("m = %v\n", m)
	equal = reflect.DeepEqual(m.Vals, [][]float64{{0}, {1}, {2}})
	if !equal {
		chk.PrintFail("DblList Append")
	}
}
