// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/vtk"
)

func cone_angle(s []float64) float64 {
	den := s[0] + s[1] + s[2]
	if den < 1e-14 {
		return 1e30
	}
	return math.Sqrt(math.Pow(s[0]-s[1], 2.0)+math.Pow(s[1]-s[2], 2.0)+math.Pow(s[2]-s[0], 2.0)) / den
}

func plot_cone(α float64, preservePrev bool) {
	//pp := 0
	//if preservePrev {
	//pp = 1
	//}
	//plt.Wireframe(X, Y, Z, io.Sf("color='k', zmin=0, zmax=0.5, preservePrev=%d", pp))
}

func main() {

	// create a new VTK Scene
	scn := vtk.NewScene()
	scn.HydroLine = true
	scn.FullAxes = true
	scn.AxesLen = 1.5

	// parameters
	α := 180.0 * math.Atan(1.0/math.Sqrt2) / math.Pi // <<< touches lower plane
	α = 90.0 - α                                     // <<< touches upper plane
	α = 15.0
	kα := math.Tan(α * math.Pi / 180.0)

	// cone
	cone := vtk.NewIsoSurf(func(x []float64) (f, vx, vy, vz float64) {
		f = cone_angle(x) - kα
		return
	})
	cone.Limits = []float64{0, -1, 0, 1, 0, 360}
	cone.Ndiv = []int{21, 21, 41}
	cone.OctRotate = true
	cone.GridShowPts = false
	cone.Color = []float64{0, 1, 1, 1}
	//cone.CmapNclrs = 0 // use this to use specified color
	cone.AddTo(scn) // remember to add to Scene

	// sphere
	sph := vtk.NewSphere()
	sph.Cen = []float64{0, 0, 0}
	sph.R = 1
	sph.AddTo(scn)

	// start interactive mode
	scn.SaveOnExit = false
	scn.Fnk = "/tmp/gosl/vtk_cone01"
	scn.Run()
}
