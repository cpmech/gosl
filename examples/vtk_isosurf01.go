// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/vtk"
)

func calc_p(x []float64) float64 {
	return -(x[0] + x[1] + x[2]) / 3.0
}

func calc_q(x []float64) float64 {
	return math.Sqrt(((x[0]-x[1])*(x[0]-x[1]) +
		(x[1]-x[2])*(x[1]-x[2]) +
		(x[2]-x[0])*(x[2]-x[0])) / 2.0)
}

func main() {

	// create a new VTK Scene
	scn := vtk.NewScene()
	scn.Reverse = true // start viewing the negative side of the x-y-z Cartesian system

	// parameters
	M := 1.0  // slope of line in p-q graph
	pt := 0.0 // tensile p
	a0 := 0.8 // size of surface

	// limits and divisions for grid generation
	pqth := []float64{pt, a0, 0, M * a0, 0, 360}
	ndiv := []int{21, 21, 41}

	// cone symbolising the Drucker-Prager criterion
	cone := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
		p, q := calc_p(x), calc_q(x)
		f = q - M*p
		return
	})
	cone.Limits = pqth
	cone.Ndiv = ndiv
	cone.OctRotate = true
	cone.GridShowPts = false
	cone.Color = []float64{0, 1, 1, 1}
	cone.CmapNclrs = 0 // use this to use specified color
	cone.AddTo(scn)    // remember to add to Scene

	// ellipsoid symbolising the Cam-clay yield surface
	ellipsoid := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
		p, q := calc_p(x), calc_q(x)
		f = q*q + M*M*(p-pt)*(p-a0)
		return
	})
	ellipsoid.Limits = pqth
	cone.Ndiv = ndiv
	ellipsoid.OctRotate = true
	ellipsoid.Color = []float64{1, 1, 0, 0.5}
	ellipsoid.CmapNclrs = 0 // use this to use specified color
	ellipsoid.AddTo(scn)    // remember to add to Scene

	// illustrate use of Arrow
	arr := vtk.NewArrow() // X0 is equal to origin
	arr.V = []float64{-1, -1, -1}
	arr.AddTo(scn)

	// illustrate use of Sphere
	sph := vtk.NewSphere()
	sph.Cen = []float64{-a0, -a0, -a0}
	sph.R = 0.05
	sph.AddTo(scn)

	// start interactive mode
	scn.SavePng = true
	scn.Fnk = "/tmp/gosl/vtk_isosurf01"
	scn.Run()
}
