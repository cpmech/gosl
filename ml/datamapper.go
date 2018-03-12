// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ml

import "github.com/cpmech/gosl/la"

// DataMapper maps features into an expanded set of features
type DataMapper interface {
	Map(x, xRaw la.Vector)                        // maps xRaw into x
	GetMapped(XYraw [][]float64, useY bool) *Data // returns new data with mapped X values
	NumOriginalFeatures() int                     // returns the number of original features
	NumExtraFeatures() int                        // returns the number of added features
}
