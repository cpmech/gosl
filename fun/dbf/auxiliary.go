// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbf

// setvzero sets v := 0
func setvzero(v []float64) {
	for i := 0; i < len(v); i++ {
		v[i] = 0
	}
}
