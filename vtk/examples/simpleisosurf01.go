// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/cpmech/gosl/vtk"
)

func main() {

	// vtk scene
	scene := vtk.NewScene()

	// setting the scene up
	scene.WithPlanes = false
	scene.HydroLine = false
	scene.AxesLen = 2

	// function
	r := 1.0 // radius
	fcn := func(x []float64) (f, vx, vy, vz float64) {
		f = x[0]*x[0] + x[1]*x[1] + x[2]*x[2] - r*r
		return
	}

	// iso-surface object
	surface := vtk.NewIsoSurf(fcn)
	surface.AddTo(scene) // add surface to scene

	// add arrow
	arrow := vtk.NewArrow()
	arrow.V = []float64{2, 2, 2}
	arrow.AddTo(scene) // add arrow to scene

	// filename for saving figure on exit
	scene.Fnk = "/tmp/example01"
	scene.SaveOnExit = true

	// add to scene and show in interactive mode
	scene.Run()
}
