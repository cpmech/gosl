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

	vars := Variables{
		&VarData{D: D_Log, M: 120, S: 12},
		&VarData{D: D_Log, M: 120, S: 12},
		&VarData{D: D_Log, M: 120, S: 12},
		&VarData{D: D_Log, M: 120, S: 12},
		&VarData{D: D_Log, M: 50, S: 15},
		&VarData{D: D_Log, M: 40, S: 12},
	}
	err := vars.Init()
	if err != nil {
		tst.Errorf("init failed:\n%v", err)
		return
	}
}
