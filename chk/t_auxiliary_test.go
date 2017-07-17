// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import "testing"

func TestMinMax01(tst *testing.T) {

	//Verbose = true
	PrintTitle("MinMax01")

	if min(1, 2) != 1 {
		tst.Errorf("min() failed\n")
		return
	}

	if min(2, 1) != 1 {
		tst.Errorf("min() failed\n")
		return
	}

	if max(1, 2) != 2 {
		tst.Errorf("max() failed\n")
		return
	}

	if max(2, 1) != 2 {
		tst.Errorf("max() failed\n")
		return
	}
}
