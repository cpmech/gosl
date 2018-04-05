// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import (
	"math"

	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/la"
)

// Kmeans implements the K-means model (Observer of Data)
type Kmeans struct {
	data      *Data       // X data
	nClasses  int         // expected number of classes
	Classes   []int       // [nSamples] indices of classes of each sample
	Centroids []la.Vector // [nClasses][nFeatures] coordinates of centroids
	Nmembers  []int       // [nClasses] number of members in each class
	bins      *gm.Bins    // "bins" to speed up searching for data points given their coordinates (2D or 3D only at the moment)
}

// NewKmeans returns a new K-means model
func NewKmeans(data *Data, nClasses int) (o *Kmeans) {

	// input
	o = new(Kmeans)
	o.data = data
	o.data.AddObserver(o) // need to recompute bins upon data changes

	// classes
	o.nClasses = nClasses
	o.Classes = make([]int, data.Nsamples)
	o.Centroids = make([]la.Vector, nClasses)
	o.Nmembers = make([]int, nClasses)

	// bins
	o.bins = new(gm.Bins)
	ndiv := []int{10, 10}                                 // TODO: make this a parameter
	o.bins.Init(o.data.Stat.MinX, o.data.Stat.MaxX, ndiv) // TODO: make sure minX and maxX are 2D or 3D; i.e. nFeatures ≤ 2
	o.Update()                                            // compute first bins
	return
}

// Update perform updates after data has been changed (as an Observer)
func (o *Kmeans) Update() {
	for i := 0; i < o.data.Nsamples; i++ {
		o.bins.Append([]float64{o.data.X.Get(i, 0), o.data.X.Get(i, 1)}, i, nil)
	}
}

// Nclasses returns the number of classes
func (o *Kmeans) Nclasses() int {
	return o.nClasses
}

// SetCentroids sets centroids; e.g. trial centroids
//   Xc -- [nClass][nFeatures]
func (o *Kmeans) SetCentroids(Xc [][]float64) {
	for i := 0; i < o.nClasses; i++ {
		o.Centroids[i] = Xc[i]
	}
}

// FindClosestCentroids finds closest centroids to each sample
func (o *Kmeans) FindClosestCentroids() {

	// loop over all samples
	del := la.NewVector(o.data.Nfeatures)
	for i := 0; i < o.data.Nsamples; i++ {

		// set min distance to max value possible
		distMin := math.MaxFloat64
		xi := la.Vector(o.data.X.GetRow(i)) // optimize here by using a row-major matrix

		// for each class
		for j := 0; j < o.nClasses; j++ {
			xc := o.Centroids[j]
			la.VecAdd(del, 1, xi, -1, xc) // del := xi - xc
			dist := del.Norm()
			if dist < distMin {
				distMin = dist
				o.Classes[i] = j
			}
		}
	}
}

// ComputeCentroids update centroids based on new classes information (from FindClosestCentroids)
func (o *Kmeans) ComputeCentroids() {

	// clear centroids and number of nMembers
	for k := 0; k < o.nClasses; k++ {
		o.Centroids[k].Fill(0.0)
		o.Nmembers[k] = 0
	}

	// add contributions to centroids and nMembers
	for i := 0; i < o.data.Nsamples; i++ {
		xi := la.Vector(o.data.X.GetRow(i)) // optimize here by using a row-major matrix
		k := o.Classes[i]
		la.VecAdd(o.Centroids[k], 1, o.Centroids[k], 1, xi)
		o.Nmembers[k]++
	}

	// scale centroids based on number of members
	for k := 0; k < o.nClasses; k++ {
		den := float64(o.Nmembers[k])
		for j := 0; j < o.data.Nfeatures; j++ {
			o.Centroids[k][j] /= den
		}
	}
}

// Train trains model
func (o *Kmeans) Train(nMaxIt int, tolNormChange float64) (nIter int) {
	for nIter = 0; nIter < nMaxIt; nIter++ {
		o.FindClosestCentroids()
		o.ComputeCentroids()
	}
	return
}
