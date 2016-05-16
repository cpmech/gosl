// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

var (
	EdgeLocalVerts map[string][][]int
	FaceLocalVerts map[string][][]int
)

func init() {

	EdgeLocalVerts = map[string][][]int{
		"lin2":  [][]int{},
		"lin3":  [][]int{},
		"lin4":  [][]int{},
		"lin5":  [][]int{},
		"tri3":  [][]int{{0, 1}, {1, 2}, {2, 0}},
		"tri6":  [][]int{{0, 1, 3}, {1, 2, 4}, {2, 0, 5}},
		"tri10": [][]int{{0, 1, 3, 6}, {1, 2, 4, 7}, {2, 0, 5, 8}},
		"tri15": [][]int{{0, 1, 3, 6, 7}, {1, 2, 4, 8, 9}, {2, 0, 5, 10, 11}},
		"qua4":  [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}},
		"qua8":  [][]int{{0, 1, 4}, {1, 2, 5}, {2, 3, 6}, {3, 0, 7}},
		"qua9":  [][]int{{0, 1, 4}, {1, 2, 5}, {2, 3, 6}, {3, 0, 7}},
		"qua12": [][]int{{0, 1, 4, 8}, {1, 2, 5, 9}, {2, 3, 6, 10}, {3, 0, 7, 11}},
		"qua16": [][]int{{0, 1, 4, 8}, {1, 2, 5, 9}, {2, 3, 6, 10}, {3, 0, 7, 11}},
	}

	FaceLocalVerts = map[string][][]int{
		"tet4":  [][]int{{0, 3, 2}, {0, 1, 3}, {0, 2, 1}, {1, 2, 3}},
		"tet10": [][]int{{0, 3, 2, 7, 9, 6}, {0, 1, 3, 4, 8, 7}, {0, 2, 1, 6, 5, 4}, {1, 2, 3, 5, 9, 8}},
		"hex8":  [][]int{{0, 4, 7, 3}, {1, 2, 6, 5}, {0, 1, 5, 4}, {2, 3, 7, 6}, {0, 3, 2, 1}, {4, 5, 6, 7}},
		"hex20": [][]int{{0, 4, 7, 3, 16, 15, 19, 11}, {1, 2, 6, 5, 9, 18, 13, 17}, {0, 1, 5, 4, 8, 17, 12, 16}, {2, 3, 7, 6, 10, 19, 14, 18}, {0, 3, 2, 1, 11, 10, 9, 8}, {4, 5, 6, 7, 12, 13, 14, 15}},
	}
}
