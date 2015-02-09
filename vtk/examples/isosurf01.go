// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"github.com/cpmech/gosl/utl"
	"github.com/cpmech/gosl/vtk"
)

func calc_I1(x []float64) float64 {
	return x[0] + x[1] + x[2]
}

func calc_sqJ2d(x []float64) float64 {
	return math.Sqrt(((x[0]-x[1])*(x[0]-x[1]) +
		(x[1]-x[2])*(x[1]-x[2]) +
		(x[2]-x[0])*(x[2]-x[0])) / 6.0)
}

func calc_p(x []float64) float64 {
	return (x[0] + x[1] + x[2]) / 3.0
}

func calc_q(x []float64) float64 {
	return math.Sqrt(((x[0]-x[1])*(x[0]-x[1]) +
		(x[1]-x[2])*(x[1]-x[2]) +
		(x[2]-x[0])*(x[2]-x[0])) / 2.0)
}

func main() {

	α := 0.45
	κ := 2000.0
	SQ3 := math.Sqrt(3.0)
	M := 3.0 * SQ3 * α
	qy0 := SQ3 * κ

	sinφ := 3.0 * M / (M + 6.0)
	φrad := math.Asin(sinφ)
	φdeg := φrad * 180.0 / math.Pi
	cosφ := math.Cos(φrad)
	tanφ := math.Tan(φrad)
	c := qy0 * tanφ / M
	pt := c / tanφ
	κ_ := 6.0 * c * cosφ / (SQ3 * (3.0 - sinφ))
	utl.Pforan("α   = %v\n", α)
	utl.Pforan("κ   = %v  (%v)\n", κ, κ_)
	utl.Pforan("φ   = %v\n", φdeg)
	utl.Pforan("c   = %v\n", c)
	utl.Pforan("M   = %v\n", M)
	utl.Pforan("qy0 = %v\n", qy0)

	scn := vtk.NewScene()
	scn.AxesLen = c * SQ3

	pqth := []float64{-pt, pt, 0, pt * M * 1.2, 0, 360}
	ndiv := []int{21, 21, 41}

	cone1 := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
		I1, sqJ2d := calc_I1(x), calc_sqJ2d(x)
		f = sqJ2d - α*I1 - κ
		return
	})
	cone1.Limits = pqth
	cone1.Ndiv = ndiv
	cone1.OctRotate = true
	cone1.GridShowPts = false
	cone1.Color = []float64{1, 0, 0, 1.0}
	cone1.CmapNclrs = 0
	cone1.AddTo(scn)

	cone2 := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
		p, q := calc_p(x), calc_q(x)
		f = q - M*p - qy0
		return
	})
	cone2.Limits = pqth
	cone2.Ndiv = ndiv
	cone2.OctRotate = true
	cone2.GridShowPts = false
	cone2.Color = []float64{0, 0, 1, 0.5}
	cone2.ShowWire = true
	cone2.CmapNclrs = 0
	cone2.AddTo(scn)

	arr := vtk.NewArrow()
	arr.X0 = []float64{-pt / SQ3, -pt / SQ3, -pt / SQ3}
	arr.V = []float64{pt, pt, pt}
	arr.Color = []float64{0, 1, 0, 1}
	arr.CyliRad = κ / 10
	arr.ConeRad = 1.2 * κ / 10
	arr.ConePct = 0.5
	arr.AddTo(scn)

	scn.Run()
}
