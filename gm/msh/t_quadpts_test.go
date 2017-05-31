// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestQuadpts01(tst *testing.T) {

	// TODO: implement this test

	//verbose()
	chk.PrintTitle("quadpts01. quadrature points")

	/*
		// 1D Gaussian rules: xi and wi
		C13 := math.Sqrt(1.0 / 3.0)
		C35 := math.Sqrt(3.0 / 5.0)
		C37 := math.Sqrt(3.0/7.0 + 2.0*math.Sqrt(6.0/5.0)/7.0)
		C73 := math.Sqrt(3.0/7.0 - 2.0*math.Sqrt(6.0/5.0)/7.0)
		C30 := math.Sqrt(30.0) / 36.0
		C57 := math.Sqrt(5.0+2.0*math.Sqrt(10.0/7.0)) / 3.0
		C75 := math.Sqrt(5.0-2.0*math.Sqrt(10.0/7.0)) / 3.0
		rule2 := [][]float64{
			{-C13, 1.0},
			{+C13, 1.0},
		}
		rule3 := [][]float64{
			{-C35, 5.0 / 9.0},
			{+0.0, 8.0 / 9.0},
			{+C35, 5.0 / 9.0},
		}
		rule4 := [][]float64{
			{-C37, 18.0 - C30},
			{+C37, 18.0 - C30},
			{+C73, 18.0 + C30},
			{-C73, 18.0 + C30},
		}
		rule5 := [][]float64{
			{-C57, (322.0 - 13.0*SQ70) / 900.0},
			{+C57, (322.0 - 13.0*SQ70) / 900.0},
			{+0.0, (322.0 - 13.0*SQ70) / 900.0},
			{+C75, (322.0 + 13.0*SQ70) / 900.0},
			{-C75, (322.0 + 13.0*SQ70) / 900.0},
		}
	*/

	for name, pts := range IntPoints {

		io.Pfyel("--------------------------------- %-6s---------------------------------\n", name)

		for n, p := range pts {
			io.Pforan("%2d: %v\n\n", n, p)
		}
	}
}
