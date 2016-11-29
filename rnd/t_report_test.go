// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func Test_report01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Report01. random variables")

	if chk.Verbose {

		dirout := "/tmp"
		fnkey := "gosl-report01"
		genPDF := true

		sets := SetsOfVars{
			&SetOfVars{
				Name: "problem 1",
				Vars: []*VarData{
					&VarData{D: D_Normal, M: 1, S: 0.1},
					&VarData{D: D_Normal, M: 2, S: 0.2},
					&VarData{D: D_Normal, M: 3, S: 0.3},
				},
			},
			&SetOfVars{
				Name: "problem:2",
				Vars: []*VarData{
					&VarData{D: D_Gumbel, M: 1, S: 0.1},
					&VarData{D: D_Gumbel, M: 2, S: 0.2},
					&VarData{D: D_Gumbel, M: 3, S: 0.3},
				},
			},
			&SetOfVars{
				Name: "problem_3",
				Vars: []*VarData{
					&VarData{D: D_Lognormal, M: 1, S: 0.1},
					&VarData{D: D_Lognormal, M: 2, S: 0.2},
					&VarData{D: D_Lognormal, M: 3, S: 0.3},
				},
			},
			&SetOfVars{
				Name: "problem-4",
				Vars: []*VarData{
					&VarData{D: D_Uniform, Min: 1, Max: 10},
					&VarData{D: D_Uniform, Min: 2, Max: 20},
					&VarData{D: D_Uniform, Min: 3, Max: 30},
				},
			},
		}

		ReportVariables(dirout, fnkey, sets, genPDF)
	}
}
