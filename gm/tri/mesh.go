// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tri

// Definitions
//
//   2D:
//              Nodes                 Edges
//
//    y           2
//    |           @                     @
//    +--x       / \                   / \
//            5 /   \ 4               /   \
//             @     @             2 /     \ 1
//            /       \             /       \
//           /         \           /         \
//          @-----@-----@         @-----------@
//         0      3      1              0
//

// Vertex holds vertex data
type Vertex struct {
	ID  int       // identifier
	Tag int       // tag
	X   []float64 // coordinates [2]
}

// Cell holds cell data
type Cell struct {
	ID       int   // identifier
	Tag      int   // tag
	V        []int // vertices
	EdgeTags []int // edge tags (2D or 3D)
}

// Mesh defines mesh data
type Mesh struct {
	Verts []*Vertex // vertices
	Cells []*Cell   // cells
}
