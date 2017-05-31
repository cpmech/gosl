// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/num"
)

func TestQuadpts01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("quadpts01. quadrature points")

	degreeMax := 5
	glX := make([][]float64, degreeMax+1)
	glW := make([][]float64, degreeMax+1)
	for n := 1; n <= degreeMax; n++ {
		glX[n], glW[n] = num.GaussLegendreXW(-1, 1, n)
	}

	for name, allPts := range IntPoints {

		io.PfYel("\n--------------------------------- %-6s---------------------------------\n", name)

		switch name {
		case "lin":
			for n, pts := range allPts {
				x := make([]float64, n)
				w := make([]float64, n)
				for i := 0; i < n; i++ {
					x[i] = pts[i][0]
					w[i] = pts[i][3]
				}
				io.Pf("\nx = %v\n", x)
				io.Pfgreen("    %v\n", glX[n])
				io.Pf("w = %v\n", w)
				io.Pfgreen("    %v\n", glW[n])
				chk.Vector(tst, io.Sf("lin:%d x", n), 1e-15, x, glX[n])
				chk.Vector(tst, io.Sf("lin:%d w", n), 1e-15, w, glW[n])
			}

		case "qua":
			for n, pts := range allPts {
				io.Pl()
				n1d := int(math.Sqrt(float64(n)))
				x1d := glX[n1d]
				w1d := glW[n1d]
				for j := 0; j < n1d; j++ {
					for i := 0; i < n1d; i++ {
						m := i + n1d*j
						x := pts[m][:2]
						v := pts[m][3]
						y := []float64{x1d[i], x1d[j]}
						w := w1d[i] * w1d[j]
						io.Pf("  %d%d x = %23v  w = %23v\n", i, j, x, v)
						io.Pfgreen("         %23v      %23v\n", y, w)
						chk.Vector(tst, "x", 1e-15, x, y)
						chk.Scalar(tst, "w", 1e-15, v, w)
					}
				}
			}

		case "hex":
			for n, pts := range allPts {
				if n == 14 {
					continue
				}
				io.Pl()
				n1d := int(math.Floor(math.Pow(float64(n), 1.0/3.0) + 0.5))
				x1d := glX[n1d]
				w1d := glW[n1d]
				for k := 0; k < n1d; k++ {
					for j := 0; j < n1d; j++ {
						for i := 0; i < n1d; i++ {
							m := i + n1d*j + (n1d*n1d)*k
							x := pts[m][:3]
							v := pts[m][3]
							y := []float64{x1d[i], x1d[j], x1d[k]}
							w := w1d[i] * w1d[j] * w1d[k]
							io.Pf("%d%d x=%18v w=%18v\n", i, j, x, v)
							io.Pfgreen("     %18v   %18v\n", y, w)
							chk.Vector(tst, "x", 1e-15, x, y)
							chk.Scalar(tst, "w", 1e-14, v, w)
						}
					}
				}
			}
		}

	}
}
