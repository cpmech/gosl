// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// PlotterClass defines a plotter to plot classification data
type PlotterClass struct {

	// input
	data    *Data // x-data
	classes []int // y-data

	// constants
	MgridNpts int // number of points for meshgrid (for contours)

	// arguments
	ArgsYclasses  map[int]*plt.A // maps y classes [0, 1, 2, ...] to plot arguments
	ArgsCentroids *plt.A         // args for centroids
	ArgsCircle1   *plt.A         // args for centroids
	ArgsCircle2   *plt.A         // args for centroids
}

// NewPlotterClass returns a new ploter
func NewPlotterClass(data *Data, classes []int, nClasses int) (o *PlotterClass) {

	// input
	o = new(PlotterClass)
	o.data = data
	o.classes = classes

	// constants
	o.MgridNpts = 21

	// arguments
	o.ArgsYclasses = make(map[int]*plt.A)
	for k := 0; k < nClasses; k++ {
		o.ArgsYclasses[k] = &plt.A{C: plt.C(k, 0), M: plt.M(k, 2), NoClip: true}
	}
	o.ArgsCentroids = &plt.A{Ls: "None", M: "*", Ms: 10, Mec: "k", NoClip: true}
	o.ArgsCircle1 = &plt.A{M: "o", Void: true, Ms: 13, Mec: "k", Mew: 4.4, NoClip: true}
	o.ArgsCircle2 = &plt.A{M: "o", Void: true, Ms: 13, Mec: "w", Mew: 1.3, NoClip: true}
	return
}

// Data plots data classes
func (o *PlotterClass) Data(iFeature, jFeature int, binary bool) {
	for iSample := 0; iSample < o.data.Nsamples; iSample++ {
		k := o.classes[iSample] % len(o.ArgsYclasses)
		args := o.ArgsYclasses[k]
		ui := o.data.X.Get(iSample, iFeature)
		vi := o.data.X.Get(iSample, jFeature)
		plt.PlotOne(ui, vi, args)
	}
	plt.HideTRborders()
	plt.Gll(io.Sf("$x_{%d}$", iFeature), io.Sf("$x_{%d}$", jFeature), nil)
}

// Centroids plots centroids of classes
func (o *PlotterClass) Centroids(centroids []la.Vector) {
	nClasses := len(centroids)
	for i := 0; i < nClasses; i++ {
		k := i % len(o.ArgsYclasses)
		o.ArgsCentroids.C = o.ArgsYclasses[k].C
		u, v := centroids[i][0], centroids[i][1]
		plt.PlotOne(u, v, o.ArgsCentroids)
		plt.PlotOne(u, v, o.ArgsCircle1)
		plt.PlotOne(u, v, o.ArgsCircle2)
		plt.Text(u, v, io.Sf("%d", i), &plt.A{Fsz: 8})
	}
}
