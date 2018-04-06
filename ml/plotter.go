// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// Plotter plots results from Machine Learning models
type Plotter struct {

	// input
	data   *Data      // data
	mapper DataMapper // mapper

	// constants
	NumPointsModelY int // number of points for ModelY()
	NumPointsModelC int // nubmer of poitns for ModelC()

	// arguments: data
	ArgsDataY     *plt.A         // args for data y
	ArgsBinClassY map[int]*plt.A // maps y classes [0 or 1] to plot arguments
	ArgsClassesY  map[int]*plt.A // maps y classes [0, 1, 2, ...] to plot arguments

	// arguments: centroids
	ArgsCentroids   *plt.A // args for centroids
	ArgsCentroCirc1 *plt.A // args for circle highlighting centroids
	ArgsCentroCirc2 *plt.A // args for circle highlighting centroids

	// arguments: model
	ArgsModelY *plt.A // arguments for x-y model line
	ArgsModelC *plt.A // arguments for ContourModel
}

// NewPlotter returns a new ploter
//   mapper -- data mapper [may be nil]
func NewPlotter(data *Data, mapper DataMapper) (o *Plotter) {

	// input
	o = new(Plotter)
	o.data = data
	o.mapper = mapper

	// constants
	o.NumPointsModelY = 21
	o.NumPointsModelC = 21

	// arguments: data
	o.ArgsDataY = &plt.A{C: plt.C(2, 0), M: plt.M(0, 0), Ls: "None", NoClip: true}
	o.ArgsBinClassY = map[int]*plt.A{
		0: {C: plt.C(0, 0), M: "o", Ls: "None", NoClip: true},
		1: {C: plt.C(2, 0), M: "*", Ls: "None", NoClip: true, Mec: plt.C(2, 0), Ms: 8},
	}
	nMaxClassesIni := 10
	o.ArgsClassesY = make(map[int]*plt.A)
	for k := 0; k < nMaxClassesIni; k++ {
		o.ArgsClassesY[k] = &plt.A{C: plt.C(k, 0), M: plt.M(k, 2), NoClip: true}
	}
	o.ArgsClassesY[0].M = "o"

	// arguments: centroids
	o.ArgsCentroids = &plt.A{Ls: "None", M: "*", Ms: 10, Mec: "k", NoClip: true}
	o.ArgsCentroCirc1 = &plt.A{M: "o", Void: true, Ms: 13, Mec: "k", Mew: 4.4, NoClip: true}
	o.ArgsCentroCirc2 = &plt.A{M: "o", Void: true, Ms: 13, Mec: "w", Mew: 1.3, NoClip: true}

	// arguments: model
	o.ArgsModelY = &plt.A{C: plt.C(0, 0), M: "None", Ls: "-", NoClip: true}
	o.ArgsModelC = &plt.A{Colors: []string{plt.C(1, 0)}, Levels: []float64{0.5}}
	return
}

// DataY plots data x[iFeature] versus data y values
func (o *Plotter) DataY(iFeature int) {
	u := o.data.X.GetCol(iFeature)
	v := o.data.Y
	plt.Plot(u, v, o.ArgsDataY)
	plt.HideTRborders()
	plt.Gll(io.Sf("$x_{%d}$", iFeature), "$y$", nil)
}

// DataClass plots data classes
//   classes -- use given classes instead of data.Y
func (o *Plotter) DataClass(nClass, iFeature, jFeature int, classes []int) {

	// argsmap
	argsmap := o.ArgsClassesY
	if nClass == 2 {
		argsmap = o.ArgsBinClassY
	}

	// plot points
	var k int
	for iSample := 0; iSample < o.data.Nsamples; iSample++ {
		if classes != nil {
			k = classes[iSample] % len(argsmap)
		} else {
			k = int(o.data.Y[iSample]) % len(argsmap)
		}
		args := argsmap[k]
		ui := o.data.X.Get(iSample, iFeature)
		vi := o.data.X.Get(iSample, jFeature)
		plt.PlotOne(ui, vi, args)
	}

	// setup
	plt.HideTRborders()
	plt.Gll(io.Sf("$x_{%d}$", iFeature), io.Sf("$x_{%d}$", jFeature), nil)
}

// Centroids plots centroids of classes
func (o *Plotter) Centroids(centroids []la.Vector) {
	nClasses := len(centroids)
	for i := 0; i < nClasses; i++ {
		k := i % len(o.ArgsClassesY)
		o.ArgsCentroids.C = o.ArgsClassesY[k].C
		u, v := centroids[i][0], centroids[i][1]
		plt.PlotOne(u, v, o.ArgsCentroids)
		plt.PlotOne(u, v, o.ArgsCentroCirc1)
		plt.PlotOne(u, v, o.ArgsCentroCirc2)
		plt.Text(u, v, io.Sf("%d", i), &plt.A{Fsz: 8})
	}
}

// ModelY plots model y values
func (o *Plotter) ModelY(model fun.Sv, iFeature int, xmin, xmax float64) {

	// x vectors
	x := la.NewVector(o.data.Nfeatures) // TODO: set x with xmean
	var xRaw la.Vector
	if o.mapper != nil {
		xRaw = la.NewVector(o.mapper.NumOriginalFeatures())
	}

	// compute points
	u := utl.LinSpace(xmin, xmax, o.NumPointsModelY)
	v := utl.GetMapped(u, func(s float64) float64 {
		if o.mapper == nil {
			x[iFeature] = s
		} else {
			xRaw[iFeature] = s
			o.mapper.Map(x, xRaw)
		}
		return model(x)
	})

	// plot line
	plt.Plot(u, v, o.ArgsModelY)
}

// ModelC plots contour defined by the model f({x} with varying x[iFeature] and x[jFeature]
func (o *Plotter) ModelC(model fun.Sv, iFeature, jFeature int, level float64, ximin, ximax, xjmin, xjmax float64) {

	// x vectors
	x := la.NewVector(o.data.Nfeatures) // TODO: set x with xmean
	var xRaw la.Vector
	if o.mapper != nil {
		xRaw = la.NewVector(o.mapper.NumOriginalFeatures())
	}

	// create meshgrid
	U, V, W := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, o.NumPointsModelC, o.NumPointsModelC, func(s, t float64) (w float64) {
		if o.mapper == nil {
			x[iFeature] = s
			x[jFeature] = t
		} else {
			xRaw[iFeature] = s
			xRaw[jFeature] = t
			o.mapper.Map(x, xRaw)
		}
		w = model(x)
		return
	})

	// plot contour
	plt.ContourL(U, V, W, o.ArgsModelC)
}

// ModelClass plots contour indicating model Classes
func (o *Plotter) ModelClass(model fun.Sv, nClass, iFeature, jFeature int, ximin, ximax, xjmin, xjmax float64) {

	// args map
	argsmap := o.ArgsClassesY
	if nClass == 2 {
		argsmap = o.ArgsBinClassY
	}

	// x vectors
	x := la.NewVector(o.data.Nfeatures) // TODO: set x with xmean
	var xRaw la.Vector
	if o.mapper != nil {
		xRaw = la.NewVector(o.mapper.NumOriginalFeatures())
	}

	// create meshgrid
	U, V, W := utl.MeshGrid2dF(ximin, ximax, xjmin, xjmax, o.NumPointsModelC, o.NumPointsModelC, func(s, t float64) (w float64) {
		if o.mapper == nil {
			x[iFeature] = s
			x[jFeature] = t
		} else {
			xRaw[iFeature] = s
			xRaw[jFeature] = t
			o.mapper.Map(x, xRaw)
		}
		w = model(x)
		return
	})

	// get colors
	colors := make([]string, nClass+1)
	for k := 0; k < nClass; k++ {
		colors[k] = argsmap[k].C
	}
	colors[nClass] = "white"

	// plot contour
	levels := utl.LinSpace(0, float64(nClass-1), nClass+1)
	plt.ContourF(U, V, W, &plt.A{Colors: colors, Levels: levels, NoLines: true, NoLabels: true})
}

// ModelClassOneVsAll plots each Model prediction using 1 = this model, 0 = other models
func (o *Plotter) ModelClassOneVsAll(models []fun.Sv, iFeature, jFeature int, ximin, ximax, xjmin, xjmax float64) {

	// set classes
	classes := make([]int, o.data.Nsamples)
	for i := 0; i < o.data.Nsamples; i++ {
		classes[i] = int(o.data.Y[i])
	}

	// plot contour
	for _, ffcn := range models {
		//pp.ArgsModelC.Colors = []string{pc.ArgsYclasses[k].C}
		o.ModelC(ffcn, iFeature, jFeature, 0.5, ximin, ximax, xjmin, xjmax)
	}
}
