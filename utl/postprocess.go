// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math"
)

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
func GetITout(all_output_times, time_stations_out []float64, tol float64) (I []int, T []float64) {
	lower_index := 0                      // lower index in all_output_times
	len_aotimes := len(all_output_times)  // length of a_o_times
	for _, t := range time_stations_out { // for all requested output times
		if t < 0 { // final time
			k := len(all_output_times) - 1     // last output time index
			I = append(I, k)                   // append index to iout
			T = append(T, all_output_times[k]) // append last output time
			continue                           // skip search
		}
		for k := lower_index; k < len_aotimes; k++ { // search within a_o_times
			if math.Abs(t-all_output_times[k]) < tol { // found near match
				lower_index += 1                   // update index
				I = append(I, k)                   // add index to iout
				T = append(T, all_output_times[k]) // add time to tout
				break                              // stop searching for this 't'
			}
			if all_output_times[k] > t { // failed to search for 't'
				lower_index = k // update idx to start from here on
				break           // skip this 't' and try the next one
			}
		}
	}
	return
}
