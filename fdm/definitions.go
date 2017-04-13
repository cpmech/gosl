// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

// Cb_src defines the callback function type for the "source" term
type Cb_src func(x, y float64, args ...interface{}) float64

// Cb_fxy defnes the callback function type f(x,y)
type Cb_fxy func(x, y float64) float64 // returns f(x,y)
