// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"time"
)

// DurGetCol returns a columns with duration values
func DurGetCol(j int, a [][]time.Duration) (col []time.Duration) {
	col = make([]time.Duration, len(a))
	for i := 0; i < len(a); i++ {
		col[i] = a[i][j]
	}
	return
}

// DurMinMaxAve returns statistics data correponding a duration data collected in 'v'
func DurMinMaxAve(v []time.Duration) (min, max, ave time.Duration) {
	min, max, ave = v[0], v[0], v[0]
	for i := 1; i < len(v); i++ {
		if v[i] < min {
			min = v[i]
		}
		if v[i] > max {
			max = v[i]
		}
		ave += v[i]
	}
	ave /= time.Duration(len(v)) * time.Nanosecond
	return
}
