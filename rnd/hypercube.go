// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// LatinIHS implements the improved distributed hypercube sampling algorithm.
// Note: code developed by John Burkardt (GNU LGPL license) --  see source code
// for further information.
//  Input:
//   dim -- spatial dimension
//   n   -- number of points to be generated
//   d   -- duplication factor ≥ 1 (~ 5 is reasonable)
//  Output:
//   x   -- [dim][n] points
func LatinIHS(dim, n, d int) (x [][]int) {

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
	var i, j, k, count, pointIndex, best int
	var minAll, minCan, dist float64

	// constant
	r8huge := 1.0E+30

	// slices
	avail := make([]int, dim*n)
	list := make([]int, d*n)
	point := make([]int, dim*d*n)
	x = utl.IntAlloc(dim, n)

	opt := float64(n) / math.Pow(float64(n), float64(1.0/float64(dim)))

	// pick the first point
	for i = 0; i < dim; i++ {
		x[i][n-1] = Int(1, n)
	}

	// initialize avail and set an entry in a random row of each column of avail to n
	for j = 0; j < n; j++ {
		for i = 0; i < dim; i++ {
			avail[i+j*dim] = j + 1
		}
	}
	for i = 0; i < dim; i++ {
		avail[i+(x[i][n-1]-1)*dim] = n
	}

	// main loop: assign a value to x[1:m,count] for count = n-1 down to 2
	for count = n - 1; 2 <= count; count-- {

		// generate valid points.
		for i = 0; i < dim; i++ {
			for k = 0; k < d; k++ {
				for j = 0; j < count; j++ {
					list[j+k*count] = avail[i+j*dim]
				}
			}

			for k = count*d - 1; 0 <= k; k-- {
				pointIndex = Int(0, k)
				point[i+k*dim] = list[pointIndex]
				list[pointIndex] = list[k]
			}
		}

		// for each candidate, determine the distance to all the
		// points that have already been selected, and save the minimum value
		minAll = r8huge
		best = 0
		for k = 0; k < d*count; k++ {
			minCan = r8huge

			for j = count; j < n; j++ {

				dist = 0.0
				for i = 0; i < dim; i++ {
					dist = dist + math.Pow(float64(point[i+k*dim])-float64(x[i][j]), 2.0)
				}
				dist = math.Sqrt(dist)

				if dist < minCan {
					minCan = dist
				}
			}

			if math.Abs(minCan-opt) < minAll {
				minAll = math.Abs(minCan - opt)
				best = k
			}

		}
		for i = 0; i < dim; i++ {
			x[i][count-1] = point[i+best*dim]
		}

		// having chosen x[:,count], update avail
		for i = 0; i < dim; i++ {
			for j = 0; j < n; j++ {
				if avail[i+j*dim] == x[i][count-1] {
					avail[i+j*dim] = avail[i+(count-1)*dim]
				}
			}
		}
	}

	// for the last point, there's only one choice
	for i = 0; i < dim; i++ {
		x[i][0] = avail[i+0*dim]
	}
	return
}

// HypercubeCoords computes the coordinates in the hypercube
//   Input:
//     sample -- the hypercube sampling indices; e.g. from LatinIHS [ndim][npoints]
//     xmin -- min limit of coordinates [ndim]
//     xmax -- max limit of coordinates [ndim]
//   Output:
//     X -- coordinates [ndim][npoints]
func HypercubeCoords(sample [][]int, xmin, xmax []float64) (X [][]float64) {
	m := len(sample)
	if m < 1 {
		return
	}
	n := len(sample[0])
	X = utl.Alloc(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			X[i][j] = xmin[i] + float64(sample[i][j]-1)*(xmax[i]-xmin[i])/float64(n-1)
		}
	}
	return
}
