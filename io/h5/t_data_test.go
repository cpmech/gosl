// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package h5

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestData01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Data01. Reading data exported from Matlab")

	dat := Open("./data", "hdf5_sample_from_matlab_online.h5", false)
	defer dat.Close()
	M := dat.GetDeep2("/M")
	chk.Deep2(tst, "M", 1e-15, M, [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
	v := dat.GetArray("/v")
	chk.Array(tst, "v", 1e-15, v, []float64{1, 2, 3})
	w := dat.GetArray("/w")
	chk.Array(tst, "w", 1e-15, w, []float64{14, 32, 50})
}
