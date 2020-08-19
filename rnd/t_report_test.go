// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"gosl/chk"
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
				Vars: []*Variable{
					{D: "N", M: 1, S: 0.1},
					{D: "N", M: 2, S: 0.2},
					{D: "N", M: 3, S: 0.3},
				},
			},
			&SetOfVars{
				Name: "problem:2",
				Vars: []*Variable{
					{D: "G", M: 1, S: 0.1},
					{D: "G", M: 2, S: 0.2},
					{D: "G", M: 3, S: 0.3},
				},
			},
			&SetOfVars{
				Name: "problem_3",
				Vars: []*Variable{
					{D: "L", M: 1, S: 0.1},
					{D: "L", M: 2, S: 0.2},
					{D: "L", M: 3, S: 0.3},
				},
			},
			&SetOfVars{
				Name: "problem-4",
				Vars: []*Variable{
					{D: "U", Min: 1, Max: 10},
					{D: "U", Min: 2, Max: 20},
					{D: "U", Min: 3, Max: 30},
				},
			},
		}

		ReportVariables(dirout, fnkey, sets, genPDF)
	}
}
