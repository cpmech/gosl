// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "math"

// LatinIHS implements the improved distributed hypercube sampling algorithm.
// Note: code developed by John Burkardt (GNU LGPL license) --  see source code
// for further information.
//  Input:
//   dim -- spatial dimension
//   n   -- number of points to be generated
//   d   -- duplication factor ≥ 1 (~ 5 is reasonable)
//  Output:
//   x   -- [dim*n] points
func LatinIHS(dim, n, d int) (x []int) {

	//  Discussion:
	//
	//    N Points in a DIM_NUM dimensional Latin hypercube are to be selected.
	//
	//    Each of the DIM_NUM coordinate dimensions is discretized to the values
	//    1 through N.  The points are to be chosen in such a way that
	//    no two points have any coordinate value in common.  This is
	//    a standard Latin hypercube requirement, and there are many
	//    solutions.
	//
	//    This algorithm differs in that it tries to pick a solution
	//    which has the property that the points are "spread out"
	//    as evenly as possible.  It does this by determining an optimal
	//    even spacing, and using the duplication factor D to allow it
	//    to choose the best of the various options available to it.
	//
	//  Licensing:
	//
	//    This code is distributed under the GNU LGPL license.
	//
	//  Modified:
	//
	//    10 April 2003
	//
	//  Author:
	//
	//    John Burkardt
	//
	//  Reference:
	//
	//    Brian Beachkofski, Ramana Grandhi,
	//    Improved Distributed Hypercube Sampling,
	//    American Institute of Aeronautics and Astronautics Paper 2002-1274.

	// auxiliary variables
	var i, j, k, count, point_index, best int
	var min_all, min_can, dist float64

	// constant
	r8_huge := 1.0E+30

	// slices
	avail := make([]int, dim*n)
	list := make([]int, d*n)
	point := make([]int, dim*d*n)
	x = make([]int, dim*n)

	opt := float64(n) / math.Pow(float64(n), float64(1.0/float64(dim)))

	// pick the first point
	for i = 0; i < dim; i++ {
		x[i+(n-1)*dim] = Int(1, n) // i4_uniform_ab(1, n, seed)
	}

	// initialize avail and set an entry in a random row of each column of avail to n
	for j = 0; j < n; j++ {
		for i = 0; i < dim; i++ {
			avail[i+j*dim] = j + 1
		}
	}
	for i = 0; i < dim; i++ {
		avail[i+(x[i+(n-1)*dim]-1)*dim] = n
	}

	// main loop: assign a value to X(1:M,COUNT) for COUNT = N-1 down to 2.
	for count = n - 1; 2 <= count; count-- {

		// generate valid points.
		for i = 0; i < dim; i++ {
			for k = 0; k < d; k++ {
				for j = 0; j < count; j++ {
					list[j+k*count] = avail[i+j*dim]
				}
			}

			for k = count*d - 1; 0 <= k; k-- {
				point_index = Int(0, k) // i4_uniform_ab(0, k, seed)
				point[i+k*dim] = list[point_index]
				list[point_index] = list[k]
			}
		}

		// for each candidate, determine the distance to all the
		// points that have already been selected, and save the minimum value.
		min_all = r8_huge
		best = 0
		for k = 0; k < d*count; k++ {
			min_can = r8_huge

			for j = count; j < n; j++ {

				dist = 0.0
				for i = 0; i < dim; i++ {
					dist = dist + math.Pow(float64(point[i+k*dim])-float64(x[i+j*dim]), 2.0)
				}
				dist = math.Sqrt(dist)

				if dist < min_can {
					min_can = dist
				}
			}

			if math.Abs(min_can-opt) < min_all {
				min_all = math.Abs(min_can - opt)
				best = k
			}

		}
		for i = 0; i < dim; i++ {
			x[i+(count-1)*dim] = point[i+best*dim]
		}

		// having chosen x[:,count], update avail.
		for i = 0; i < dim; i++ {
			for j = 0; j < n; j++ {
				if avail[i+j*dim] == x[i+(count-1)*dim] {
					avail[i+j*dim] = avail[i+(count-1)*dim]
				}
			}
		}
	}

	// for the last point, there's only one choice.
	for i = 0; i < dim; i++ {
		x[i+0*dim] = avail[i+0*dim]
	}
	return
}
