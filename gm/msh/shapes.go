// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

// cell kinds
const (
	KindLin = 0 // "lin" cell kind
	KindTri = 1 // "tri" cell kind
	KindQua = 2 // "qua" cell kind
	KindTet = 3 // "tet" cell kind
	KindHex = 4 // "hex" cell kind
)

// cell types
const (
	TypeLin2   = 0  // Lin2 cell type index
	TypeLin3   = 1  // Lin3 cell type index
	TypeLin4   = 2  // Lin4 cell type index
	TypeLin5   = 3  // Lin5 cell type index
	TypeTri3   = 4  // Tri3 cell type index
	TypeTri6   = 5  // Tri6 cell type index
	TypeTri10  = 6  // Tri10 cell type index
	TypeTri15  = 7  // Tri15 cell type index
	TypeQua4   = 8  // Qua4 cell type index
	TypeQua8   = 9  // Qua8 cell type index
	TypeQua9   = 10 // Qua9 cell type index
	TypeQua12  = 11 // Qua12 cell type index
	TypeQua16  = 12 // Qua16 cell type index
	TypeTet4   = 13 // Tet4 cell type index
	TypeTet10  = 14 // Tet10 cell type index
	TypeHex8   = 15 // Hex8 cell type index
	TypeHex20  = 16 // Hex20 cell type index
	TypeNumMax = 17 // max number of types
)

// ShapeFunction computes the shape function and derivatives
type ShapeFunction func(S []float64, dSdR [][]float64, R []float64, derivs bool)

var (
	Functions       []ShapeFunction // shape functions and derivatives [TypeNumMax]
	TypeKeyToIndex  map[string]int  // converts type key (e.g. "lin2") to index (e.g. TypeLin2)
	TypeIndexToKey  []string        // converts type index (e.g. TypeLin2) to key (e.g. "lin2")
	TypeIndexToKind []int           // converts type index (e.g. TypeLin2) to cell kind (e.g. KindLin)
	NumVerts        []int           // number of vertices on shape [TypeNumMax]
	GeomNdim        []int           // geometry number of space dimensions [TypeNumMax]
	EdgeLocalVerts  [][][]int       // local indices of vertices on edges of shape [TypeNumMax][nedges][nverts]
	FaceLocalVerts  [][][]int       // local indices of vertices on faces of shape [TypeNumMax][nfaces][nverts]
	NatCoords       [][][]float64   // natural coordinates of vertices on shape [TypeNumMax][nverts][gndim]
)

func init() {

	// set Functions
	Functions = make([]ShapeFunction, TypeNumMax)
	Functions[TypeLin2] = FuncLin2
	Functions[TypeLin3] = FuncLin3
	Functions[TypeLin4] = FuncLin4
	Functions[TypeLin5] = FuncLin5
	Functions[TypeTri3] = FuncTri3
	Functions[TypeTri6] = FuncTri6
	Functions[TypeTri10] = FuncTri10
	Functions[TypeTri15] = FuncTri15
	Functions[TypeQua4] = FuncQua4
	Functions[TypeQua8] = FuncQua8
	Functions[TypeQua9] = FuncQua9
	Functions[TypeQua12] = FuncQua12
	Functions[TypeQua16] = FuncQua16
	Functions[TypeTet4] = FuncTet4
	Functions[TypeTet10] = FuncTet10
	Functions[TypeHex8] = FuncHex8
	Functions[TypeHex20] = FuncHex20

	// set TypeKeyToIndex
	TypeKeyToIndex = make(map[string]int)
	TypeKeyToIndex["lin2"] = TypeLin2
	TypeKeyToIndex["lin3"] = TypeLin3
	TypeKeyToIndex["lin4"] = TypeLin4
	TypeKeyToIndex["lin5"] = TypeLin5
	TypeKeyToIndex["tri3"] = TypeTri3
	TypeKeyToIndex["tri6"] = TypeTri6
	TypeKeyToIndex["tri10"] = TypeTri10
	TypeKeyToIndex["tri15"] = TypeTri15
	TypeKeyToIndex["qua4"] = TypeQua4
	TypeKeyToIndex["qua8"] = TypeQua8
	TypeKeyToIndex["qua9"] = TypeQua9
	TypeKeyToIndex["qua12"] = TypeQua12
	TypeKeyToIndex["qua16"] = TypeQua16
	TypeKeyToIndex["tet4"] = TypeTet4
	TypeKeyToIndex["tet10"] = TypeTet10
	TypeKeyToIndex["hex8"] = TypeHex8
	TypeKeyToIndex["hex20"] = TypeHex20

	// set TypeIndexToKey
	TypeIndexToKey = make([]string, TypeNumMax)
	TypeIndexToKey[TypeLin2] = "lin2"
	TypeIndexToKey[TypeLin3] = "lin3"
	TypeIndexToKey[TypeLin4] = "lin4"
	TypeIndexToKey[TypeLin5] = "lin5"
	TypeIndexToKey[TypeTri3] = "tri3"
	TypeIndexToKey[TypeTri6] = "tri6"
	TypeIndexToKey[TypeTri10] = "tri10"
	TypeIndexToKey[TypeTri15] = "tri15"
	TypeIndexToKey[TypeQua4] = "qua4"
	TypeIndexToKey[TypeQua8] = "qua8"
	TypeIndexToKey[TypeQua9] = "qua9"
	TypeIndexToKey[TypeQua12] = "qua12"
	TypeIndexToKey[TypeQua16] = "qua16"
	TypeIndexToKey[TypeTet4] = "tet4"
	TypeIndexToKey[TypeTet10] = "tet10"
	TypeIndexToKey[TypeHex8] = "hex8"
	TypeIndexToKey[TypeHex20] = "hex20"

	// set TypeIndexToKind
	TypeIndexToKind = make([]int, TypeNumMax)
	TypeIndexToKind[TypeLin2] = KindLin
	TypeIndexToKind[TypeLin3] = KindLin
	TypeIndexToKind[TypeLin4] = KindLin
	TypeIndexToKind[TypeLin5] = KindLin
	TypeIndexToKind[TypeTri3] = KindTri
	TypeIndexToKind[TypeTri6] = KindTri
	TypeIndexToKind[TypeTri10] = KindTri
	TypeIndexToKind[TypeTri15] = KindTri
	TypeIndexToKind[TypeQua4] = KindQua
	TypeIndexToKind[TypeQua8] = KindQua
	TypeIndexToKind[TypeQua9] = KindQua
	TypeIndexToKind[TypeQua12] = KindQua
	TypeIndexToKind[TypeQua16] = KindQua
	TypeIndexToKind[TypeTet4] = KindTet
	TypeIndexToKind[TypeTet10] = KindTet
	TypeIndexToKind[TypeHex8] = KindHex
	TypeIndexToKind[TypeHex20] = KindHex

	// set NumVerts
	NumVerts = make([]int, TypeNumMax)
	NumVerts[TypeLin2] = 2
	NumVerts[TypeLin3] = 3
	NumVerts[TypeLin4] = 4
	NumVerts[TypeLin5] = 5
	NumVerts[TypeTri3] = 3
	NumVerts[TypeTri6] = 6
	NumVerts[TypeTri10] = 10
	NumVerts[TypeTri15] = 15
	NumVerts[TypeQua4] = 4
	NumVerts[TypeQua8] = 8
	NumVerts[TypeQua9] = 9
	NumVerts[TypeQua12] = 12
	NumVerts[TypeQua16] = 16
	NumVerts[TypeTet4] = 4
	NumVerts[TypeTet10] = 10
	NumVerts[TypeHex8] = 8
	NumVerts[TypeHex20] = 20

	// set GeomNdim
	GeomNdim = make([]int, TypeNumMax)
	GeomNdim[TypeLin2] = 1
	GeomNdim[TypeLin3] = 1
	GeomNdim[TypeLin4] = 1
	GeomNdim[TypeLin5] = 1
	GeomNdim[TypeTri3] = 2
	GeomNdim[TypeTri6] = 2
	GeomNdim[TypeTri10] = 2
	GeomNdim[TypeTri15] = 2
	GeomNdim[TypeQua4] = 2
	GeomNdim[TypeQua8] = 2
	GeomNdim[TypeQua9] = 2
	GeomNdim[TypeQua12] = 2
	GeomNdim[TypeQua16] = 2
	GeomNdim[TypeTet4] = 3
	GeomNdim[TypeTet10] = 3
	GeomNdim[TypeHex8] = 3
	GeomNdim[TypeHex20] = 3

	// set EdgeLocalVerts
	EdgeLocalVerts = make([][][]int, TypeNumMax)
	EdgeLocalVerts[TypeTri3] = [][]int{{0, 1}, {1, 2}, {2, 0}}
	EdgeLocalVerts[TypeTri6] = [][]int{{0, 1, 3}, {1, 2, 4}, {2, 0, 5}}
	EdgeLocalVerts[TypeTri10] = [][]int{{0, 1, 3, 6}, {1, 2, 4, 7}, {2, 0, 5, 8}}
	EdgeLocalVerts[TypeTri15] = [][]int{{0, 1, 3, 6, 7}, {1, 2, 4, 8, 9}, {2, 0, 5, 10, 11}}
	EdgeLocalVerts[TypeQua4] = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}}
	EdgeLocalVerts[TypeQua8] = [][]int{{0, 1, 4}, {1, 2, 5}, {2, 3, 6}, {3, 0, 7}}
	EdgeLocalVerts[TypeQua9] = [][]int{{0, 1, 4}, {1, 2, 5}, {2, 3, 6}, {3, 0, 7}}
	EdgeLocalVerts[TypeQua12] = [][]int{{0, 1, 4, 8}, {1, 2, 5, 9}, {2, 3, 6, 10}, {3, 0, 7, 11}}
	EdgeLocalVerts[TypeQua16] = [][]int{{0, 1, 4, 8}, {1, 2, 5, 9}, {2, 3, 6, 10}, {3, 0, 7, 11}}
	EdgeLocalVerts[TypeTet4] = [][]int{{0, 1}, {1, 2}, {2, 0}, {0, 3}, {1, 3}, {2, 3}}
	EdgeLocalVerts[TypeTet10] = [][]int{{0, 1, 4}, {1, 2, 5}, {2, 0, 6}, {0, 3, 7}, {1, 3, 8}, {2, 3, 9}}
	EdgeLocalVerts[TypeHex8] = [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}, {4, 5}, {5, 6}, {6, 7}, {7, 4}, {0, 4}, {1, 5}, {2, 6}, {3, 7}}
	EdgeLocalVerts[TypeHex20] = [][]int{{0, 1, 8}, {1, 2, 9}, {2, 3, 10}, {3, 0, 11}, {4, 5, 12}, {5, 6, 13}, {6, 7, 14}, {7, 4, 15}, {0, 4, 16}, {1, 5, 17}, {2, 6, 18}, {3, 7, 19}}

	// set FaceLocalVerts
	FaceLocalVerts = make([][][]int, TypeNumMax)
	FaceLocalVerts[TypeTet4] = [][]int{{0, 3, 2}, {0, 1, 3}, {0, 2, 1}, {1, 2, 3}}
	FaceLocalVerts[TypeTet10] = [][]int{{0, 3, 2, 7, 9, 6}, {0, 1, 3, 4, 8, 7}, {0, 2, 1, 6, 5, 4}, {1, 2, 3, 5, 9, 8}}
	FaceLocalVerts[TypeHex8] = [][]int{{0, 4, 7, 3}, {1, 2, 6, 5}, {0, 1, 5, 4}, {2, 3, 7, 6}, {0, 3, 2, 1}, {4, 5, 6, 7}}
	FaceLocalVerts[TypeHex20] = [][]int{{0, 4, 7, 3, 16, 15, 19, 11}, {1, 2, 6, 5, 9, 18, 13, 17}, {0, 1, 5, 4, 8, 17, 12, 16}, {2, 3, 7, 6, 10, 19, 14, 18}, {0, 3, 2, 1, 11, 10, 9, 8}, {4, 5, 6, 7, 12, 13, 14, 15}}

	// set NatCoords
	NatCoords = make([][][]float64, TypeNumMax)
	NatCoords[TypeLin2] = [][]float64{
		{-1, 1},
	}
	NatCoords[TypeLin3] = [][]float64{
		{-1, 1, 0},
	}
	NatCoords[TypeLin4] = [][]float64{
		{-1, 1, -1.0 / 3.0, 1.0 / 3.0},
	}
	NatCoords[TypeLin5] = [][]float64{
		{-1, 1, 0, -0.5, 0.5},
	}
	NatCoords[TypeTri3] = [][]float64{
		{0, 1, 0},
		{0, 0, 1},
	}
	NatCoords[TypeTri6] = [][]float64{
		{0, 1, 0, 0.5, 0.5, 0},
		{0, 0, 1, 0, 0.5, 0.5},
	}
	NatCoords[TypeTri10] = [][]float64{
		{0, 1, 0, 1.0 / 3.0, 2.0 / 3.0, 0, 2.0 / 3.0, 1.0 / 3.0, 0, 1.0 / 3.0},
		{0, 0, 1, 0, 1.0 / 3.0, 2.0 / 3.0, 0, 2.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0},
	}
	NatCoords[TypeTri15] = [][]float64{
		{0, 1, 0, 0.5, 0.5, 0, 0.25, 0.75, 0.75, 0.25, 0, 0, 0.25, 0.5, 0.25},
		{0, 0, 1, 0, 0.5, 0.5, 0, 0, 0.25, 0.75, 0.75, 0.25, 0.25, 0.25, 0.5},
	}
	NatCoords[TypeQua4] = [][]float64{
		{-1, 1, 1, -1},
		{-1, -1, 1, 1},
	}
	NatCoords[TypeQua8] = [][]float64{
		{-1, 1, 1, -1, 0, 1, 0, -1},
		{-1, -1, 1, 1, -1, 0, 1, 0},
	}
	NatCoords[TypeQua9] = [][]float64{
		{-1, 1, 1, -1, 0, 1, 0, -1, 0},
		{-1, -1, 1, 1, -1, 0, 1, 0, 0},
	}
	NatCoords[TypeQua12] = [][]float64{
		{-1, 1, 1, -1, -1.0 / 3.0, 1, 1.0 / 3.0, -1, 1.0 / 3.0, 1, -1.0 / 3.0, -1},
		{-1, -1, 1, 1, -1, -1.0 / 3.0, 1, 1.0 / 3.0, -1, 1.0 / 3.0, 1, -1.0 / 3.0},
	}
	NatCoords[TypeQua16] = [][]float64{
		{-1, 1, 1, -1, -1.0 / 3.0, 1, 1.0 / 3.0, -1, 1.0 / 3.0, 1, -1.0 / 3.0, -1, -1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0, -1.0 / 3.0},
		{-1, -1, 1, 1, -1, -1.0 / 3.0, 1, 1.0 / 3.0, -1, 1.0 / 3.0, 1, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0},
	}
	NatCoords[TypeTet4] = [][]float64{
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	NatCoords[TypeTet10] = [][]float64{
		{0, 1, 0, 0, 0.5, 0.5, 0, 0, 0.5, 0},
		{0, 0, 1, 0, 0, 0.5, 0.5, 0, 0, 0.5},
		{0, 0, 0, 1, 0, 0, 0, 0.5, 0.5, 0.5},
	}
	NatCoords[TypeHex8] = [][]float64{
		{-1, 1, 1, -1, -1, 1, 1, -1},
		{-1, -1, 1, 1, -1, -1, 1, 1},
		{-1, -1, -1, -1, 1, 1, 1, 1},
	}
	NatCoords[TypeHex20] = [][]float64{
		{-1, 1, 1, -1, -1, 1, 1, -1, 0, 1, 0, -1, 0, 1, 0, -1, -1, 1, 1, -1},
		{-1, -1, 1, 1, -1, -1, 1, 1, -1, 0, 1, 0, -1, 0, 1, 0, -1, -1, 1, 1},
		{-1, -1, -1, -1, 1, 1, 1, 1, -1, -1, -1, -1, 1, 1, 1, 1, 0, 0, 0, 0},
	}
}
