// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

// PointN /////////////////////////////////////////////////////////////////////////////////////////

// PointN implements a point in N-dim space
type PointN struct {

	// esssential
	X []float64 // coordinates

	// optional
	Id    int         // some identification number
	Extra interface{} // some extra data
}

// TODO: test
func NewPointNdim(ndim uint32) (o *PointN) {
	return &PointN{X: make([]float64, ndim)}
}

// TODO: test
func NewPointN(X ...float64) (o *PointN) {
	ndim := len(X)
	o = new(PointN)
	o.X = make([]float64, ndim)
	copy(o.X, X)
	return
}

// GetCloneX returns a new point with X cloned, but not the other data
func (o PointN) GetCloneX() (p *PointN) {
	p = new(PointN)
	p.X = make([]float64, len(o.X))
	copy(p.X, o.X)
	return
}

func (o PointN) ExactlyTheSameX(p *PointN) bool {
	for i := 0; i < len(o.X); i++ {
		if o.X[i] != p.X[i] {
			return false
		}
	}
	return true
}

func (o PointN) AlmostTheSameX(p *PointN, tol float64) bool {
	for i := 0; i < len(o.X); i++ {
		if math.Abs(o.X[i]-p.X[i]) > tol {
			return false
		}
	}
	return true
}

func DistPointPointN(p *PointN, q *PointN) (dist float64) {
	for i := 0; i < len(p.X); i++ {
		dist += math.Pow(q.X[i]-p.X[i], 2.0)
	}
	return math.Sqrt(dist)
}

// BoxN ///////////////////////////////////////////////////////////////////////////////////////////

type BoxN struct {
	// essential
	Lo *PointN
	Hi *PointN

	// auxiliary
	Id int
}

// TODO: test
// L -- limits
func NewBoxN(L ...float64) *BoxN {
	if len(L) == 4 { // 2D
		return &BoxN{Lo: NewPointN(L[0], L[1]), Hi: NewPointN(L[2], L[3])}
	} else if len(L) == 6 { // 3D
		return &BoxN{Lo: NewPointN(L[0], L[1], L[2]), Hi: NewPointN(L[3], L[4], L[5])}
	} else {
		chk.Panic("NewBoxN requires 4 (2D) or 6 (3D) numbers; e.g. xmin,xmax, ymin,ymax, zmin,zmax")
	}
	return nil
}

func (o BoxN) IsInside(p *PointN) bool {
	for i := 0; i < len(p.X); i++ {
		if p.X[i] < o.Lo.X[i] {
			return false
		}
		if p.X[i] > o.Hi.X[i] {
			return false
		}
	}
	return true
}

func (o *BoxN) GetSize() (delta []float64) {
	delta = make([]float64, len(o.Lo.X))
	for i := 0; i < len(o.Lo.X); i++ {
		delta[i] = o.Hi.X[i] - o.Lo.X[i]
	}
	return
}

func (o *BoxN) GetMid() (mid []float64) {
	mid = make([]float64, len(o.Lo.X))
	for i := 0; i < len(o.Lo.X); i++ {
		mid[i] = (o.Lo.X[i] + o.Hi.X[i]) / 2.0
	}
	return
}

func (o BoxN) Draw(withTxt bool, args, argsTxt *plt.A) {
	if argsTxt == nil {
		argsTxt = &plt.A{C: "k", Fsz: 8, Ha: "center", Va: "center"}
	}
	lbl := io.Sf("%d", o.Id)
	ndim := len(o.Lo.X)
	mid := o.GetMid()
	if ndim == 2 {
		if args == nil {
			args = &plt.A{Ec: "#005cd9", Fc: "#dfe7f3", Closed: true, NoClip: true}
		}
		pts := [][]float64{
			[]float64{o.Lo.X[0], o.Lo.X[1]},
			[]float64{o.Hi.X[0], o.Lo.X[1]},
			[]float64{o.Hi.X[0], o.Hi.X[1]},
			[]float64{o.Lo.X[0], o.Hi.X[1]},
		}
		plt.Polyline(pts, args)
		if withTxt {
			plt.Text(mid[0], mid[1], lbl, argsTxt)
		}
		plt.AxisRange(o.Lo.X[0], o.Hi.X[0], o.Lo.X[1], o.Hi.X[1])
	} else {
		if args == nil {
			args = &plt.A{Ec: "#005cd9", Fc: "#dfe7f3", A: 0.7, Lw: 2, Ls: ":"}
		}
		plt.Box(o.Lo.X[0], o.Hi.X[0], o.Lo.X[1], o.Hi.X[1], o.Lo.X[2], o.Hi.X[2], args)
		if withTxt {
			plt.Text3d(mid[0], mid[1], mid[2], lbl, argsTxt)
		}
		plt.AxisRange3d(o.Lo.X[0], o.Hi.X[0], o.Lo.X[1], o.Hi.X[1], o.Lo.X[2], o.Hi.X[2])
	}
}

// If point p lies outside box b, the distance to the nearest point on b is returned.
// If p is inside b or on its surface, zero is returned.
func DistPointBoxN(p *PointN, b *BoxN) (dist float64) {
	for i := 0; i < len(p.X); i++ {
		if p.X[i] < b.Lo.X[i] {
			dist += math.Pow(p.X[i]-b.Lo.X[i], 2.0)
		}
		if p.X[i] > b.Hi.X[i] {
			dist += math.Pow(p.X[i]-b.Hi.X[i], 2.0)
		}
	}
	return math.Sqrt(dist)
}

// Octree /////////////////////////////////////////////////////////////////////////////////////////

type Entity interface {
}

type Octree struct {

	// constants
	DIM  uint32 // dimension
	PMAX uint32 // roughly how many levels fit in 32 bits
	QO   uint32 // 4 for quadtree, 8 for octree
	QL   uint32 // offset constant to leftmost daughter

	// internal
	maxd    uint32            // max depth: number of levels to be represented
	blo     []float64         // [DIM]
	bscale  []float64         // [DIM]
	elhash  map[uint32]Entity // contains stored elements hashed by box #
	pophash map[uint32]uint32 // contains node population info
}

// NewOctree creates an Octree
//   ndim -- dimension
func NewOctree(L ...float64) (o *Octree) {

	// allocate object
	o = new(Octree)

	// dimension
	var ndim uint32
	if len(L) == 4 { // 2D
		ndim = 2
		o.blo = []float64{L[0], L[2]}
		o.bscale = []float64{L[1] - L[0], L[3] - L[2]}
	} else if len(L) == 6 { // 3D
		ndim = 3
		o.blo = []float64{L[0], L[2], L[4]}
		o.bscale = []float64{L[1] - L[0], L[3] - L[2], L[5] - L[4]}
	} else {
		chk.Panic("NewOctree requires 4 (2D) or 6 (3D) numbers; e.g. xmin,xmax, ymin,ymax, zmin,zmax")
	}

	// constants
	o.DIM = uint32(ndim)
	o.PMAX = 32 / ndim
	o.QO = 4
	if ndim == 3 {
		o.QO = 8
	}
	o.QL = o.QO - 2

	// internal
	o.maxd = o.PMAX
	return
}

// qobox creates new box whose index is k. The root box is k==1
func (o *Octree) qobox(k uint32) (box *BoxN) {
	box = &BoxN{Lo: NewPointNdim(o.DIM), Hi: NewPointNdim(o.DIM), Id: int(k)}
	var j, kb uint32
	offset := make([]float64, o.DIM)
	del := 1.0
	for k > 1 { // up through ancestors until get to root.
		kb = (k + o.QL) % o.QO // which daughter is k? Add its offset.
		for j = 0; j < o.DIM; j++ {
			if kb&(1<<j) != 0 {
				offset[j] += del
			}
		}
		k = (k + o.QL) >> o.DIM // replace k by its mother,
		del *= 2.0              // where offsets will be twice as big.
	}
	for j = 0; j < o.DIM; j++ { // At the end, scale the offsets by the final del to
		box.Lo.X[j] = o.blo[j] + o.bscale[j]*offset[j]/del // make them metrically correct.
		box.Hi.X[j] = o.blo[j] + o.bscale[j]*(offset[j]+1.0)/del
	}
	return
}
