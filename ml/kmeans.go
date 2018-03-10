// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import "github.com/cpmech/gosl/la"

type Kmeans struct {
	clusters *la.Matrix
}

/*
func (o *Kmeans) Run() {
	// Kmax = max number of clusters
	for K := 0; K < Kmax; K++ {
		for e := 0; e < nTrials; e++ {
			go func() { // run in parallel
			}()
		}
	}
}
*/
