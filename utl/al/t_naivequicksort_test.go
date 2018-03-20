// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestNaiveQuickSort01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NaiveQuickSort01.")

	A := []int{0, 1, 4, 5, -1, 3, 8, 2}
	NaiveQuickSort(A, func(a, b int) int {
		if ToInt(a) > ToInt(b) {
			return +1
		}
		if ToInt(a) < ToInt(b) {
			return -1
		}
		return 0
	})
	chk.Ints(tst, "A.sorted", A, []int{-1, 0, 1, 2, 3, 4, 5, 8})
}
