// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import (
	"testing"
)

func TestRecover01(tst *testing.T) {

	//Verbose = true
	PrintTitle("Recover01. Panic is OK")

	defer RecoverTstPanicIsOK(tst)
	panic("panic now")
}
