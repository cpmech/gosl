// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tri

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
)

func Test_tri01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("tri01")

	V := [][]float64{
		{0.0, 0.0},
		{1.0, 0.0},
		{1.0, 1.0},
		{0.0, 1.0},
		{0.5, 0.5},
	}

	C := [][]int{
		{0, 1, 4},
		{1, 2, 4},
		{2, 3, 4},
		{3, 0, 4},
	}

	if chk.Verbose {
		plt.SetForPng(1, 300, 150, nil)
		Draw(V, C, nil)
		plt.Equal()
		plt.AxisRange(-0.1, 1.1, -0.1, 1.1)
		plt.Gll("x", "y", nil)
		plt.SaveD("/tmp/gosl", "t_tri01.png")
	}
}
