// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_octree01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("octree01. PointN")

	p := &PointN{X: []float64{1, 2, 3}}

	io.Pforan("p = %+v\n", p)
	chk.Vector(tst, "p.X", 1e-15, p.X, []float64{1, 2, 3})
	q := &PointN{X: []float64{2, 2, 1}}
	io.Pforan("q = %+v\n", q)
	chk.Vector(tst, "q.X", 1e-15, q.X, []float64{2, 2, 1})
	if p.ExactlyTheSameX(q) {
		tst.Errorf("ExactlyTheSame should return false because points are indeed different")
		return
	}
	if p.AlmostTheSameX(q, 1e-15) {
		tst.Errorf("AlmostTheSame should return false because points are different within given tolerance (1e-15)")
		return
	}
	if p.AlmostTheSameX(q, 1.0) {
		tst.Errorf("AlmostTheSame should return false because points are different within given tolerance (1.0)")
		return
	}
	if !p.AlmostTheSameX(q, 2.0) {
		tst.Errorf("AlmostTheSame should return true because points are different within given tolerance (2.0)")
		return
	}

	a := p.GetCloneX()
	chk.Vector(tst, "a == p", 1e-15, a.X, []float64{1, 2, 3})

	dap := DistPointPointN(a, p)
	chk.Scalar(tst, "dist(a,p)", 1e-15, dap, 0)

	dpq := DistPointPointN(p, q)
	chk.Scalar(tst, "dist(p,q)", 1e-15, dpq, math.Sqrt(5.0))
}

func Test_octree02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("octree02. BoxN and DistPointPointN")

	b := &BoxN{&PointN{X: []float64{-1, -2, -3}}, &PointN{X: []float64{3, 2, 1}}, 0}
	delta := b.GetSize()
	chk.Vector(tst, "delta", 1e-15, delta, []float64{4, 4, 4})

	p := &PointN{X: []float64{-2, 0, 0}}
	dist := DistPointBoxN(p, b)
	io.Pforan("dist = %v\n", dist)
	chk.Scalar(tst, "dist(p,b)", 1e-15, dist, 1.0)

	if b.IsInside(p) {
		tst.Errorf("is inside box failed")
		return
	}

	q := &PointN{X: []float64{-2, 3, 0}}
	dist = DistPointBoxN(q, b)
	io.Pforan("dist = %v\n", dist)
	chk.Scalar(tst, "dist(q,b)", 1e-15, dist, math.Sqrt2)

	if b.IsInside(q) {
		tst.Errorf("is inside box failed")
		return
	}

	r := &PointN{X: []float64{-2, 3, 2}}
	dist = DistPointBoxN(r, b)
	io.Pforan("dist = %v\n", dist)
	chk.Scalar(tst, "dist(r,b)", 1e-15, dist, math.Sqrt(3.0))

	if b.IsInside(r) {
		tst.Errorf("is inside box failed")
		return
	}

	s := &PointN{X: []float64{0, 0, 0}}
	dist = DistPointBoxN(s, b)
	io.Pforan("dist = %v\n", dist)
	chk.Scalar(tst, "dist(s,b)", 1e-15, dist, 0)

	if !b.IsInside(s) {
		tst.Errorf("is inside box failed")
		return
	}

	s.X[0] = 1.0
	dist = DistPointBoxN(s, b)
	io.Pforan("dist = %v\n", dist)
	chk.Scalar(tst, "dist(s,b)", 1e-15, dist, 0)

	if !b.IsInside(s) {
		tst.Errorf("is inside box failed")
		return
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		b.Draw(true, nil, nil)
		plt.Plot3dPoint(p.X[0], p.X[1], p.X[2], &plt.A{C: "r", Ec: "r"})
		plt.Plot3dPoint(q.X[0], q.X[1], q.X[2], &plt.A{C: "g", Ec: "g"})
		plt.Plot3dPoint(r.X[0], r.X[1], r.X[2], &plt.A{C: "y", Ec: "y"})
		plt.Triad(1, "x", "y", "z", nil, nil)
		plt.Default3dView(-3, 3, -3, 3, -3, 3, true)
		//err := plt.ShowSave("/tmp/gosl", "octree02")
		err := plt.Save("/tmp/gosl", "octree02")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_octree03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("octree03. qobox")

	o := NewOctree(-1, 1, -2, 1) // 4 numbers => 2D
	chk.Vector(tst, "blo", 1e-15, o.blo, []float64{-1, -2})
	chk.Vector(tst, "bscale", 1e-15, o.bscale, []float64{2, 3})

	b1 := o.qobox(1)
	b2 := o.qobox(2)
	b3 := o.qobox(3)
	b4 := o.qobox(4)
	b5 := o.qobox(5)
	b6 := o.qobox(6)
	b7 := o.qobox(7)
	b15 := o.qobox(15)
	b21 := o.qobox(21)
	b22 := o.qobox(22)
	b41 := o.qobox(41)
	b45 := o.qobox(45)
	b52 := o.qobox(52)
	chk.Vector(tst, "1: lo", 1e-15, b1.Lo.X, []float64{-1, -2})
	chk.Vector(tst, "1: hi", 1e-15, b1.Hi.X, []float64{1, 1})
	chk.Vector(tst, "2: lo", 1e-15, b2.Lo.X, []float64{-1, -2})
	chk.Vector(tst, "2: hi", 1e-15, b2.Hi.X, []float64{0, -0.5})
	chk.Vector(tst, "3: lo", 1e-15, b3.Lo.X, []float64{0, -2})
	chk.Vector(tst, "3: hi", 1e-15, b3.Hi.X, []float64{1, -0.5})
	chk.Vector(tst, "4: lo", 1e-15, b4.Lo.X, []float64{-1, -0.5})
	chk.Vector(tst, "4: hi", 1e-15, b4.Hi.X, []float64{0, 1.0})
	chk.Vector(tst, "5: lo", 1e-15, b5.Lo.X, []float64{0, -0.5})
	chk.Vector(tst, "5: hi", 1e-15, b5.Hi.X, []float64{1, 1})
	chk.Vector(tst, "6: lo", 1e-15, b6.Lo.X, []float64{-1, -2})
	chk.Vector(tst, "6: hi", 1e-15, b6.Hi.X, []float64{-0.5, -1.25})
	chk.Vector(tst, "7: lo", 1e-15, b7.Lo.X, []float64{-0.5, -2})
	chk.Vector(tst, "7: hi", 1e-15, b7.Hi.X, []float64{0, -1.25})
	chk.Vector(tst, "15: lo", 1e-15, b15.Lo.X, []float64{-0.5, -0.5})
	chk.Vector(tst, "15: hi", 1e-15, b15.Hi.X, []float64{0, 0.25})
	chk.Vector(tst, "21: lo", 1e-15, b21.Lo.X, []float64{0.5, 0.25})
	chk.Vector(tst, "21: hi", 1e-15, b21.Hi.X, []float64{1, 1})
	chk.Vector(tst, "22: lo", 1e-15, b22.Lo.X, []float64{-1, -2})
	chk.Vector(tst, "22: hi", 1e-15, b22.Hi.X, []float64{-0.75, -1.625})
	chk.Vector(tst, "41: lo", 1e-15, b41.Lo.X, []float64{0.25, -1.625})
	chk.Vector(tst, "41: hi", 1e-15, b41.Hi.X, []float64{0.5, -1.25})
	chk.Vector(tst, "45: lo", 1e-15, b45.Lo.X, []float64{0.75, -1.625})
	chk.Vector(tst, "45: hi", 1e-15, b45.Hi.X, []float64{1, -1.25})
	chk.Vector(tst, "52: lo", 1e-15, b52.Lo.X, []float64{0.5, -0.875})
	chk.Vector(tst, "52: hi", 1e-15, b52.Hi.X, []float64{0.75, -0.5})

	if chk.Verbose {
		b85 := o.qobox(85)
		b200 := o.qobox(200)
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		b1.Draw(true, nil, nil)
		b2.Draw(true, &plt.A{C: "red", A: 0.2}, nil)
		b3.Draw(true, &plt.A{C: "blue", A: 0.2}, nil)
		b4.Draw(true, &plt.A{C: "yellow", A: 0.2}, nil)
		b5.Draw(true, &plt.A{C: "green", A: 0.2}, nil)
		b6.Draw(true, &plt.A{C: "green", A: 0.2}, nil)
		b7.Draw(true, &plt.A{C: "yellow", A: 0.2}, nil)
		b15.Draw(true, &plt.A{C: "red", A: 0.2}, nil)
		b21.Draw(true, &plt.A{C: "red", A: 0.2}, nil)
		b22.Draw(true, &plt.A{C: "yellow", A: 0.2}, nil)
		b41.Draw(true, &plt.A{C: "yellow", A: 0.2}, nil)
		b45.Draw(true, &plt.A{C: "blue", A: 0.2}, nil)
		b52.Draw(true, &plt.A{C: "red", A: 0.2}, nil)
		b85.Draw(true, &plt.A{C: "blue", A: 0.2}, nil)
		b200.Draw(true, &plt.A{C: "green", A: 0.2}, nil)
		plt.AxisRange(-3, 3, -3, 3)
		plt.Equal()
		plt.HideAllBorders()
		plt.Grid(nil)
		err := plt.Save("/tmp/gosl", "octree03")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
