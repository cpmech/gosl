// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestSetVzero(tst *testing.T) {

	//verbose()
	chk.PrintTitle("SetVzero")

	v := []float64{1, 2, 3, 4}
	setvzero(v)
	chk.Array(tst, "v=zero", 1e-15, v, nil)
}
