// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_vars01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("vars01")

	var vars Variables
	vars.Data = []*VarData{
		&VarData{D: "log", M: 120, S: 12, Std: true},
	}
	err := vars.Init()
	if err != nil {
		tst.Errorf("init failed:\n%v", err)
		return
	}
}
