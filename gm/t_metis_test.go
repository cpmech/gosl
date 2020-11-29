// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestMetisShares01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("MetisShares01")

	// graph
	//               0        1
	//           0-------1--------2
	//           |       |        |
	//          2|      3|       4|
	//           |       |        |
	//           3-------4--------5
	//               5       6
	//
	edges := [][2]int{
		{0, 1}, {1, 2},
		{0, 3}, {1, 4}, {2, 5},
		{3, 4}, {4, 5},
	}

	// shares
	shares := MetisShares(edges)
	io.Pf("%v\n", shares)
	chk.Ints(tst, "edges attached to vertex 0", shares[0], []int{0, 2})
	chk.Ints(tst, "edges attached to vertex 1", shares[1], []int{0, 1, 3})
	chk.Ints(tst, "edges attached to vertex 2", shares[2], []int{1, 4})
	chk.Ints(tst, "edges attached to vertex 3", shares[3], []int{2, 5})
	chk.Ints(tst, "edges attached to vertex 4", shares[4], []int{3, 5, 6})
	chk.Ints(tst, "edges attached to vertex 5", shares[5], []int{4, 6})
}

func TestMetisAdjacency01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("MetisAdjacency01")

	// graph
	//
	//          0        1        2        3
	//      0--------1--------2--------3--------4
	//      |        |        |        |        |
	//    12|      13|      14|      15|      16|
	//      |   4    |   5    |   6    |   7    |
	//      5--------6--------7--------8--------9
	//      |        |        |        |        |
	//    17|      18|      19|      20|      21|
	//      |   8    |   9    |   10   |   11   |
	//     10-------11-------12-------13-------14
	//
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4},
		{5, 6}, {6, 7}, {7, 8}, {8, 9},
		{10, 11}, {11, 12}, {12, 13}, {13, 14},
		{0, 5}, {1, 6}, {2, 7}, {3, 8}, {4, 9},
		{5, 10}, {6, 11}, {7, 12}, {8, 13}, {9, 14},
	}

	// shares map
	shares := MetisShares(edges)
	chk.Ints(tst, "edges attached to vertex 0", shares[0], []int{0, 12})
	chk.Ints(tst, "edges attached to vertex 4", shares[4], []int{3, 16})
	chk.Ints(tst, "edges attached to vertex 5", shares[5], []int{4, 12, 17})
	chk.Ints(tst, "edges attached to vertex 8", shares[8], []int{6, 7, 15, 20})
	chk.Ints(tst, "edges attached to vertex 10", shares[10], []int{8, 17})
	chk.Ints(tst, "edges attached to vertex 14", shares[14], []int{11, 21})

	// adjacency list
	xadj, adjncy := MetisAdjacency(edges, shares)
	chk.Int32s(tst, "xadj", xadj, []int32{0, 2, 5, 8, 11, 13, 16, 20, 24, 28, 31, 33, 36, 39, 42, 44})
	chk.Int32s(tst, "adjncy", adjncy, []int32{1, 5, 0, 2, 6, 1, 3, 7, 2, 4, 8, 3, 9, 6, 0, 10, 5, 7, 1, 11, 6, 8, 2, 12, 7, 9, 3, 13, 8, 4, 14, 11, 5, 10, 12, 6, 11, 13, 7, 12, 14, 8, 13, 9})
}

func TestMetisPartitionLowLevel01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("MetisPartitionLowLevel01")

	// graph
	//
	//          0        1        2        3
	//      0--------1--------2--------3--------4
	//      |        |        |        |        |
	//    12|      13|      14|      15|      16|
	//      |   4    |   5    |   6    |   7    |
	//      5--------6--------7--------8--------9
	//      |        |        |        |        |
	//    17|      18|      19|      20|      21|
	//      |   8    |   9    |   10   |   11   |
	//     10-------11-------12-------13-------14
	//
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4},
		{5, 6}, {6, 7}, {7, 8}, {8, 9},
		{10, 11}, {11, 12}, {12, 13}, {13, 14},
		{0, 5}, {1, 6}, {2, 7}, {3, 8}, {4, 9},
		{5, 10}, {6, 11}, {7, 12}, {8, 13}, {9, 14},
	}

	// shares map
	shares := MetisShares(edges)

	// adjacency list
	xadj, adjncy := MetisAdjacency(edges, shares)

	// flags
	npart := 2
	recursive := false

	// partition graph
	nv := len(shares)
	objval, parts := MetisPartitionLowLevel(npart, nv, xadj, adjncy, recursive)
	io.Pf("objval = %v\n", objval)
	io.Pf("parts = %v\n", parts)
	chk.Int32(tst, "objval", objval, 7)
	chk.Int32s(tst, "parts", parts, []int32{0, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1})
}

func TestMetisPartition01(tst *testing.T) {

	// verbose()
	chk.PrintTitle("MetisPartition01")

	// graph
	//
	//          0        1        2        3
	//      0--------1--------2--------3--------4
	//      |        |        |        |        |
	//    12|      13|      14|      15|      16|
	//      |   4    |   5    |   6    |   7    |
	//      5--------6--------7--------8--------9
	//      |        |        |        |        |
	//    17|      18|      19|      20|      21|
	//      |   8    |   9    |   10   |   11   |
	//     10-------11-------12-------13-------14
	//
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4},
		{5, 6}, {6, 7}, {7, 8}, {8, 9},
		{10, 11}, {11, 12}, {12, 13}, {13, 14},
		{0, 5}, {1, 6}, {2, 7}, {3, 8}, {4, 9},
		{5, 10}, {6, 11}, {7, 12}, {8, 13}, {9, 14},
	}

	// flags
	npart := 2
	recursive := false

	// partition graph
	shares, objval, parts := MetisPartition(edges, npart, recursive)
	io.Pf("objval = %v\n", objval)
	io.Pf("parts = %v\n", parts)
	chk.Int32(tst, "objval", objval, 7)
	chk.Int32s(tst, "parts", parts, []int32{0, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1})

	// check shares
	chk.Ints(tst, "edges attached to vertex 0", shares[0], []int{0, 12})
	chk.Ints(tst, "edges attached to vertex 4", shares[4], []int{3, 16})
	chk.Ints(tst, "edges attached to vertex 5", shares[5], []int{4, 12, 17})
	chk.Ints(tst, "edges attached to vertex 8", shares[8], []int{6, 7, 15, 20})
	chk.Ints(tst, "edges attached to vertex 10", shares[10], []int{8, 17})
	chk.Ints(tst, "edges attached to vertex 14", shares[14], []int{11, 21})
}
