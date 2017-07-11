// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import "math"

// GetITout returns indices and output times
//  Input:
//    all_output_times  -- array with all output times. ex: [0,0.1,0.2,0.22,0.3,0.4]
//    time_stations_out -- time stations for output: ex: [0,0.2,0.4]
//    tol               -- tolerance to compare times
//  Output:
//    iout -- indices of times in all_output_times
//    tout -- times corresponding to iout
//  Notes:
//    use -1 in all_output_times to enforce output of last timestep
func GetITout(allOutputTimes, timeStationsOut []float64, tol float64) (I []int, T []float64) {
	lowerIndex := 0                     // lower index in all_output_times
	lenAoTimes := len(allOutputTimes)   // length of a_o_times
	for _, t := range timeStationsOut { // for all requested output times
		if t < 0 { // final time
			k := len(allOutputTimes) - 1     // last output time index
			I = append(I, k)                 // append index to iout
			T = append(T, allOutputTimes[k]) // append last output time
			continue                         // skip search
		}
		for k := lowerIndex; k < lenAoTimes; k++ { // search within a_o_times
			if math.Abs(t-allOutputTimes[k]) < tol { // found near match
				lowerIndex++                     // update index
				I = append(I, k)                 // add index to iout
				T = append(T, allOutputTimes[k]) // add time to tout
				break                            // stop searching for this 't'
			}
			if allOutputTimes[k] > t { // failed to search for 't'
				lowerIndex = k // update idx to start from here on
				break          // skip this 't' and try the next one
			}
		}
	}
	return
}

// GetStrides returns nReq indices from 0 (inclusive) to nTotal (inclusive)
//   Input:
//     nTotal -- total number of intices
//     nReq -- required indices
//   Example:
//     GetStrides(2001, 5) => [0 400 800 1200 1600 2000 2001]
//
//   NOTE: GetStrides will always include nTotal as the last item in I
//
func GetStrides(nTotal, nReq int) (I []int) {
	if nReq > nTotal {
		nReq = nTotal
	}
	lt := nTotal / nReq
	if lt < 1 {
		lt = 1
	}
	I = IntRange3(0, nTotal, lt)
	if I[len(I)-1] != nTotal {
		I = append(I, nTotal)
	}
	return
}
