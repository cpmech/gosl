// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chk

import "testing"

func Test_err01(tst *testing.T) {

	//Verbose = true
	defer func() { Verbose = false }()
	PrintTitle("err01")

	err := Err("TEST_ERR = %d", 666)
	if err.Error() != "TEST_ERR = 666" {
		panic("Err failed")
	}

	PrintOk("hello from PrintOk => ")
	PrintFail("hello from PrintFail => ")
}
