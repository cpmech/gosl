// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "github.com/cpmech/gosl/vtk"

func main() {

	// create a new VTK Scene
	scn := vtk.NewScene()
	scn.HydroLine = true
	scn.FullAxes = true
	scn.AxesLen = 1.5

	// sphere
	sphere := &vtk.Sphere{
		Cen:   []float64{0, 0, 0},
		R:     1.0,
		Color: []float64{85.0 / 255.0, 128.0 / 255.0, 225.0 / 255.0, 1},
	}
	sphere.AddTo(scn)

	// spheres
	sset := vtk.NewSpheresFromFile("data/points.dat")
	sset.AddTo(scn)

	// start interactive mode
	scn.SavePng = true
	scn.Fnk = "/tmp/gosl/vtk_spheres01"
	scn.Run()
}
