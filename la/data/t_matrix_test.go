// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestA123(tst *testing.T) {

	verbose()
	chk.PrintTitle("A123.")

	a := new(A123)
	a.Generate()
	chk.Vector(tst, "A123", 1e-17, a.A, []float64{1, 4, 7, 2, 5, 8, 3, 6, 9})
}
