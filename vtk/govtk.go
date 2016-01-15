// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vtk

/*
#cgo CXXFLAGS: -O3 -I. -I/usr/include/eigen3 -I/usr/include/vtk-5.10 -DVTK_EXCLUDE_STRSTREAM_HEADERS
#cgo LDFLAGS: -lstdc++ -lm -lvtkRendering -lvtkHybrid -lvtkParallel -lvtkGraphics -lvtkImaging -lvtkIO -lvtkFiltering -lvtkCommon -lvtksys
#include "connectgovtk.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// Scene holds essential data to draw and interect with the x-y-z Cartesian system
type Scene struct {

	// options
	AxesLen    float64 // length of x-y-z axes
	HydroLine  bool    // show hydrostatic line
	Reverse    bool    // reverse direction for default camera
	FullAxes   bool    // show negative and positive portions of axes
	WithPlanes bool    // show transparent auxiliary planes
	Interact   bool    // run interactive mode
	SaveOnExit bool    // save figure upon exit
	Fnk        string  // file name key (without .png)

	// vtk objects
	arrows     []*Arrow
	spheres    []*Sphere
	spheresSet []*Spheres
	isosurfs   []*IsoSurf

	// c data
	win unsafe.Pointer // GoslVTK::Win
}

// Arrow adds an arrow to Scene
type Arrow struct {

	// options
	X0         []float64 // origin of arrow
	V          []float64 // vector defining arrow
	ConePct    float64   // percentage of length to draw tip cone
	ConeRad    float64   // radius of cone
	CyliRad    float64   // cylinder radius
	Resolution int       // resolution of cross-section
	Color      []float64 // {red, green, blue, opacity}

	// c data
	arr unsafe.Pointer // GoslVTK::Arrow
}

// Sphere adds a sphere to Scene
type Sphere struct {

	// options
	Cen   []float64 // centre x-y-z coordinates
	R     float64   // radius
	Color []float64 // {red, green, blue, opacity}

	// c data
	sph unsafe.Pointer // GoslVTK::Sphere
}

// Spheres adds a set of spheres (e.g. particles) to Scene
type Spheres struct {

	// options
	X     []float64 // x coordinates
	Y     []float64 // y coordinates
	Z     []float64 // z coordinates
	R     []float64 // radii
	Color []float64 // {red, green, blue, opacity}

	// c data
	sset unsafe.Pointer // GoslVTK::Sphere
}

// Cb_fcn is a callback function to compute f, v := f(x) where v = dfdx
type Cb_fcn func(x []float64) (fval, vx, vy, vz float64)

// IsoSurf holds data to generate isosurfaces
//  CmapNclrs:
//      0 => use fixed color
//  CmapRangeType:
//      0 => [default] use sgrid range values (automatic)
//      1 => use Frange
//      2 => use CmapFrange
type IsoSurf struct {

	// options
	Limits        []float64 // {xmin,xmax, ymin,ymax, zmin,zmax}
	Ndiv          []int     // {nx, ny, nz}. all must be >= 2
	Frange        []float64 // {fmin, fmax}. min and max values of isosurface levels
	OctRotate     bool      // apply rotation to octahedral reference system
	Nlevels       int       // number of isosurface levels (0 or 1 => just one @ Levelmin)
	CmapType      string    // colormap type. e.g. "warm"
	CmapNclrs     int       // colormap number of colors
	CmapRangeType int       // colormap range type
	CmapFrange    []float64 // colormap fmin and fmax
	Color         []float64 // {red, green, blue, opacity}
	ShowWire      bool      // show wireframe of main object
	GridShowPts   bool      // show underlying grid
	fcn           Cb_fcn    // function

	// c data
	isf unsafe.Pointer // GoslVTK::IsoSurf
}

// NewScene allocates a new Scene structure
func NewScene() *Scene {
	return &Scene{
		AxesLen:    1.0,
		HydroLine:  true,
		FullAxes:   true,
		WithPlanes: true,
		Interact:   true,
		Fnk:        "tmp_gosl_vtk",
	}
}

// NewArrow allocates a new Arrow structure
func NewArrow() *Arrow {
	return &Arrow{
		X0:         []float64{0, 0, 0},
		V:          []float64{1, 1, 1},
		ConePct:    0.1,
		ConeRad:    0.03,
		CyliRad:    0.015,
		Resolution: 20,
		Color:      []float64{1, 0, 0, 1},
	}
}

// NewSphere allocates a new Sphere structure
func NewSphere() *Sphere {
	return &Sphere{
		Cen:   []float64{0, 0, 0},
		R:     1.0,
		Color: []float64{0, 1, 1, 1},
	}
}

// NewSpheres allocates a new set of spheres structure
func NewSpheres() *Spheres {
	return &Spheres{
		X:     []float64{0, 1, 1, 0, 0, 1, 1, 0},
		Y:     []float64{0, 0, 1, 1, 0, 0, 1, 1},
		Z:     []float64{0, 0, 0, 0, 1, 1, 1, 1},
		R:     []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		Color: []float64{1, 0, 0, 1},
	}
}

// NewSpheresFromFile add spheres (e.g. particles) by reading a file in the following format
//   x    y    z    r
//  0.0  0.0  0.0  0.1
//  1.0  0.0  0.0  0.1
//   ...
func NewSpheresFromFile(filename string) *Spheres {
	_, dat, err := io.ReadTable(filename)
	if err != nil {
		chk.Panic("NewSpheresFromFile failed:\n%v", err)
	}
	return &Spheres{
		X:     dat["x"],
		Y:     dat["y"],
		Z:     dat["z"],
		R:     dat["r"],
		Color: []float64{1, 0, 0, 1},
	}
}

// NewIsoSurf allocates a new IsoSurf structure
func NewIsoSurf(f Cb_fcn) *IsoSurf {
	return &IsoSurf{
		Limits:     []float64{-1, 1, -1, 1, -1, 1},
		Ndiv:       []int{21, 21, 21},
		Frange:     []float64{0, 1},
		Nlevels:    1,
		CmapType:   "warm",
		CmapNclrs:  16,
		CmapFrange: []float64{0, 1},
		Color:      []float64{0, 0, 1, 1},
		fcn:        f,
	}
}

// AddTo adds Arrow to Scene
func (o *Arrow) AddTo(scn *Scene) (err error) {
	scn.arrows = append(scn.arrows, o)
	return
}

// AddTo adds Sphere to Scene
func (o *Sphere) AddTo(scn *Scene) (err error) {
	scn.spheres = append(scn.spheres, o)
	return
}

// AddTo adds Spheres to Scene
func (o *Spheres) AddTo(scn *Scene) (err error) {
	scn.spheresSet = append(scn.spheresSet, o)
	return
}

// AddTo adds IsoSurf to Scene
func (o *IsoSurf) AddTo(scn *Scene) (err error) {
	scn.isosurfs = append(scn.isosurfs, o)
	return
}

// b2i converts bool to int
func b2i(b bool) (i int) {
	if b {
		return 1
	}
	return 0
}

// Run shows Scene in interactive mode or saving a .png file
func (o *Scene) Run() (err error) {

	// input data
	axeslen := (C.double)(o.AxesLen)
	hydroline := (C.long)(b2i(o.HydroLine))
	reverse := (C.long)(b2i(o.Reverse))
	fullaxes := (C.long)(b2i(o.FullAxes))
	withplanes := (C.long)(b2i(o.WithPlanes))
	interact := (C.long)(b2i(o.Interact))
	saveonexit := (C.long)(b2i(o.SaveOnExit))
	fnk := C.CString(o.Fnk)
	defer C.free(unsafe.Pointer(fnk))

	// connect Go and C
	GOVTK_X = make([]float64, 3)
	C.GOVTK_F = (*C.double)(unsafe.Pointer(&GOVTK_F))
	C.GOVTK_VX = (*C.double)(unsafe.Pointer(&GOVTK_VX))
	C.GOVTK_VY = (*C.double)(unsafe.Pointer(&GOVTK_VX))
	C.GOVTK_VZ = (*C.double)(unsafe.Pointer(&GOVTK_VX))
	C.GOVTK_X = (*C.double)(unsafe.Pointer(&GOVTK_X[0]))
	C.GOVTK_I = (*C.long)(unsafe.Pointer(&GOVTK_I))

	// alloc win
	o.win = C.win_alloc(reverse)
	defer C.win_dealloc(o.win)
	if o.win == nil {
		return chk.Err("C.scene_begin failed\n")
	}

	// arrows
	for _, O := range o.arrows {
		x0 := (*C.double)(unsafe.Pointer(&O.X0[0]))
		v := (*C.double)(unsafe.Pointer(&O.V[0]))
		cone_pct := (C.double)(O.ConePct)
		cone_rad := (C.double)(O.ConeRad)
		cyli_rad := (C.double)(O.CyliRad)
		resolution := (C.long)(O.Resolution)
		color := (*C.double)(unsafe.Pointer(&O.Color[0]))
		O.arr = C.arrow_addto(o.win, x0, v, cone_pct, cone_rad, cyli_rad, resolution, color)
		defer C.arrow_dealloc(O.arr)
	}

	// spheres
	for _, O := range o.spheres {
		cen := (*C.double)(unsafe.Pointer(&O.Cen[0]))
		r := (C.double)(O.R)
		color := (*C.double)(unsafe.Pointer(&O.Color[0]))
		O.sph = C.sphere_addto(o.win, cen, r, color)
		defer C.sphere_dealloc(O.sph)
	}

	// spheres set
	for _, O := range o.spheresSet {
		n := len(O.X)
		if n < 1 {
			continue
		}
		if len(O.Y) != n || len(O.Z) != n || len(O.R) != n {
			return chk.Err("cannot add set of spheres because X,Y,Z,R have different dimensions")
		}
		x := (*C.double)(unsafe.Pointer(&O.X[0]))
		y := (*C.double)(unsafe.Pointer(&O.Y[0]))
		z := (*C.double)(unsafe.Pointer(&O.Z[0]))
		r := (*C.double)(unsafe.Pointer(&O.R[0]))
		color := (*C.double)(unsafe.Pointer(&O.Color[0]))
		O.sset = C.spheres_addto(o.win, (C.long)(n), x, y, z, r, color)
		defer C.spheres_dealloc(O.sset)
	}

	// isosurfs
	for _, O := range o.isosurfs {

		// input data
		limits := (*C.double)(unsafe.Pointer(&O.Limits[0]))
		ndiv := (*C.long)(unsafe.Pointer(&O.Ndiv[0]))
		frange := (*C.double)(unsafe.Pointer(&O.Frange[0]))
		octrotate := (C.long)(b2i(O.OctRotate))
		nlevels := (C.long)(O.Nlevels)
		cmaptype := C.CString(O.CmapType)
		cmapnclrs := (C.long)(O.CmapNclrs)
		cmaprangetype := (C.long)(O.CmapRangeType)
		cmapfrange := (*C.double)(unsafe.Pointer(&O.CmapFrange[0]))
		color := (*C.double)(unsafe.Pointer(&O.Color[0]))
		showwire := (C.long)(b2i(O.ShowWire))
		gridshowpts := (C.long)(b2i(O.GridShowPts))
		defer C.free(unsafe.Pointer(cmaptype))

		// connect Go and C
		idx := len(GOVTK_FCN)
		GOVTK_FCN = append(GOVTK_FCN, O.fcn)

		// call C routine: add isosurf
		index := (C.long)(idx)
		O.isf = C.isosurf_addto(o.win, index, limits, ndiv, frange, octrotate, nlevels, cmaptype, cmapnclrs, cmaprangetype, cmapfrange, color, showwire, gridshowpts)
		defer C.isosurf_dealloc(O.isf)
	}

	// call C routine: end
	status := C.scene_run(o.win, axeslen, hydroline, reverse, fullaxes, withplanes, interact, saveonexit, fnk)
	if status != 0 {
		return chk.Err("C.scene_end failed\n")
	}
	return
}

// global variables for communication with C
var (
	GOVTK_FCN []Cb_fcn
	GOVTK_F   float64
	GOVTK_VX  float64
	GOVTK_VY  float64
	GOVTK_VZ  float64
	GOVTK_X   []float64
	GOVTK_I   int
)

//export govtk_isosurf_fcn
func govtk_isosurf_fcn() {
	GOVTK_F, GOVTK_VX, GOVTK_VY, GOVTK_VZ = GOVTK_FCN[GOVTK_I](GOVTK_X)
}
