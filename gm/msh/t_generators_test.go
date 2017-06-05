// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func TestGen01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Gen01. 3x2 region (lower-level)")

	// grid generator
	xmin, xmax := 1.0, 3.0
	ymin, ymax := 0.0, 1.5
	f := func(i, j, nr, ns int) (x, y float64) {
		dx := (xmax - xmin) / float64(nr-1)
		dy := (ymax - ymin) / float64(ns-1)
		x = xmin + float64(i)*dx
		y = ymin + float64(j)*dy
		return
	}

	// constants
	circle := false
	ndivU, ndivV := 3, 2

	// plot arguments
	args := NewArgs()
	args.WithEdges = true
	args.WithCells = true
	args.WithVerts = true
	args.WithIdsCells = true
	args.WithIdsVerts = true
	args.WithTagsVerts = true
	args.WithTagsEdges = true

	// generate many quads
	ctypes := []int{TypeQua4, TypeQua8, TypeQua9, TypeQua12, TypeQua16, TypeQua17}
	for _, ctype := range ctypes {
		mesh, err := GenQuadRegion(ctype, ndivU, ndivV, circle, f)
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// check connectivity
		n := ndivU*ndivV - 1
		switch ctype {
		case TypeQua4:
			chk.Ints(tst, "Qua4:  c0.V", mesh.Cells[0].V, []int{0, 1, 5, 4})
			chk.Ints(tst, "Qua4:  cn.V", mesh.Cells[n].V, []int{6, 7, 11, 10})
		case TypeQua8:
			chk.Ints(tst, "Qua8:  c0.V", mesh.Cells[0].V, []int{0, 2, 13, 11, 1, 8, 12, 7})
			chk.Ints(tst, "Qua8:  cn.V", mesh.Cells[n].V, []int{15, 17, 28, 26, 16, 21, 27, 20})
		case TypeQua9:
			chk.Ints(tst, "Qua9:  c0.V", mesh.Cells[0].V, []int{0, 2, 16, 14, 1, 9, 15, 7, 8})
			chk.Ints(tst, "Qua9:  cn.V", mesh.Cells[n].V, []int{18, 20, 34, 32, 19, 27, 33, 25, 26})
		case TypeQua12:
			chk.Ints(tst, "Qua12: c0.V", mesh.Cells[0].V, []int{0, 3, 21, 18, 1, 11, 20, 14, 2, 15, 19, 10})
			chk.Ints(tst, "Qua12: cn.V", mesh.Cells[n].V, []int{24, 27, 45, 42, 25, 31, 44, 34, 26, 35, 43, 30})
		case TypeQua16:
			chk.Ints(tst, "Qua16: c0.V", mesh.Cells[0].V, []int{0, 3, 33, 30, 1, 13, 32, 20, 2, 23, 31, 10, 11, 12, 22, 21})
			chk.Ints(tst, "Qua16: cn.V", mesh.Cells[n].V, []int{36, 39, 69, 66, 37, 49, 68, 56, 38, 59, 67, 46, 47, 48, 58, 57})
		case TypeQua17:
			chk.Ints(tst, "Qua17: c0.V", mesh.Cells[0].V, []int{0, 4, 32, 28, 1, 14, 31, 24, 2, 19, 30, 17, 3, 25, 29, 13, 18})
			chk.Ints(tst, "Qua17: cn.V", mesh.Cells[n].V, []int{36, 40, 68, 64, 37, 44, 67, 54, 38, 51, 66, 49, 39, 55, 65, 43, 50})
		}

		// compute tags map
		tm, err := mesh.GetTagMaps()
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// vertices along left border
		vset := tm.VertTag2verts[4]
		vset = append(vset, tm.VertTag2verts[41]...)
		vset = append(vset, tm.VertTag2verts[34]...)
		for _, v := range vset {
			chk.Scalar(tst, "x=xmin", 1e-15, v.X[0], xmin)
		}

		// vertices along right border
		vset = tm.VertTag2verts[2]
		vset = append(vset, tm.VertTag2verts[12]...)
		vset = append(vset, tm.VertTag2verts[23]...)
		for _, v := range vset {
			chk.Scalar(tst, "x=xmax", 1e-15, v.X[0], xmax)
		}

		// plot
		if chk.Verbose {
			plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
			mesh.Draw(args)
			plt.HideAllBorders()
			plt.Save("/tmp/gosl/gm", io.Sf("region-%s", TypeIndexToKey[ctype]))
		}
	}
}

func TestGen02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Gen02. 3x2 region (higher-level)")

	// constants
	xmin, xmax := 1.0, 3.0
	ymin, ymax := 0.0, 1.5
	ndivU, ndivV := 3, 2

	// generate many quads
	ctypes := []int{TypeQua4, TypeQua8, TypeQua9, TypeQua12, TypeQua16, TypeQua17}
	for _, ctype := range ctypes {
		mesh, err := GenQuadRegionHL(ctype, ndivU, ndivV, xmin, xmax, ymin, ymax)
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// check connectivity
		n := ndivU*ndivV - 1
		switch ctype {
		case TypeQua4:
			chk.Ints(tst, "Qua4:  c0.V", mesh.Cells[0].V, []int{0, 1, 5, 4})
			chk.Ints(tst, "Qua4:  cn.V", mesh.Cells[n].V, []int{6, 7, 11, 10})
		case TypeQua8:
			chk.Ints(tst, "Qua8:  c0.V", mesh.Cells[0].V, []int{0, 2, 13, 11, 1, 8, 12, 7})
			chk.Ints(tst, "Qua8:  cn.V", mesh.Cells[n].V, []int{15, 17, 28, 26, 16, 21, 27, 20})
		case TypeQua9:
			chk.Ints(tst, "Qua9:  c0.V", mesh.Cells[0].V, []int{0, 2, 16, 14, 1, 9, 15, 7, 8})
			chk.Ints(tst, "Qua9:  cn.V", mesh.Cells[n].V, []int{18, 20, 34, 32, 19, 27, 33, 25, 26})
		case TypeQua12:
			chk.Ints(tst, "Qua12: c0.V", mesh.Cells[0].V, []int{0, 3, 21, 18, 1, 11, 20, 14, 2, 15, 19, 10})
			chk.Ints(tst, "Qua12: cn.V", mesh.Cells[n].V, []int{24, 27, 45, 42, 25, 31, 44, 34, 26, 35, 43, 30})
		case TypeQua16:
			chk.Ints(tst, "Qua16: c0.V", mesh.Cells[0].V, []int{0, 3, 33, 30, 1, 13, 32, 20, 2, 23, 31, 10, 11, 12, 22, 21})
			chk.Ints(tst, "Qua16: cn.V", mesh.Cells[n].V, []int{36, 39, 69, 66, 37, 49, 68, 56, 38, 59, 67, 46, 47, 48, 58, 57})
		case TypeQua17:
			chk.Ints(tst, "Qua17: c0.V", mesh.Cells[0].V, []int{0, 4, 32, 28, 1, 14, 31, 24, 2, 19, 30, 17, 3, 25, 29, 13, 18})
			chk.Ints(tst, "Qua17: cn.V", mesh.Cells[n].V, []int{36, 40, 68, 64, 37, 44, 67, 54, 38, 51, 66, 49, 39, 55, 65, 43, 50})
		}
	}
}

func TestGen03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Gen03. 2d ring (lower-level)")

	// grid generator
	alpha := math.Pi / 2.0
	r, R := 1.0, 3.0
	f := func(i, j, nr, ns int) (x, y float64) {
		dr := (R - r) / float64(nr-1)
		da := alpha / float64(ns-1)
		a := float64(j) * da
		l := r + float64(i)*dr
		x = l * math.Cos(a)
		y = l * math.Sin(a)
		return
	}

	// constants
	circle := false
	ndivU, ndivV := 3, 5

	// plot arguments
	args := NewArgs()
	args.WithEdges = true
	args.WithCells = true
	args.WithVerts = true
	args.WithIdsVerts = true

	// generate many quads
	ctypes := []int{TypeQua4, TypeQua8, TypeQua9, TypeQua12, TypeQua16, TypeQua17}
	for _, ctype := range ctypes {
		mesh, err := GenQuadRegion(ctype, ndivU, ndivV, circle, f)
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// check connectivity
		n := ndivU*ndivV - 1
		switch ctype {
		case TypeQua4:
			chk.Ints(tst, "Qua4:  c0.V", mesh.Cells[0].V, []int{0, 1, 5, 4})
			chk.Ints(tst, "Qua4:  cn.V", mesh.Cells[n].V, []int{18, 19, 23, 22})
		case TypeQua8:
			chk.Ints(tst, "Qua8:  c0.V", mesh.Cells[0].V, []int{0, 2, 13, 11, 1, 8, 12, 7})
			chk.Ints(tst, "Qua8:  cn.V", mesh.Cells[n].V, []int{48, 50, 61, 59, 49, 54, 60, 53})
		case TypeQua9:
			chk.Ints(tst, "Qua9:  c0.V", mesh.Cells[0].V, []int{0, 2, 16, 14, 1, 9, 15, 7, 8})
			chk.Ints(tst, "Qua9:  cn.V", mesh.Cells[n].V, []int{60, 62, 76, 74, 61, 69, 75, 67, 68})
		case TypeQua12:
			chk.Ints(tst, "Qua12: c0.V", mesh.Cells[0].V, []int{0, 3, 21, 18, 1, 11, 20, 14, 2, 15, 19, 10})
			chk.Ints(tst, "Qua12: cn.V", mesh.Cells[n].V, []int{78, 81, 99, 96, 79, 85, 98, 88, 80, 89, 97, 84})
		case TypeQua16:
			chk.Ints(tst, "Qua16: c0.V", mesh.Cells[0].V, []int{0, 3, 33, 30, 1, 13, 32, 20, 2, 23, 31, 10, 11, 12, 22, 21})
			chk.Ints(tst, "Qua16: cn.V", mesh.Cells[n].V, []int{126, 129, 159, 156, 127, 139, 158, 146, 128, 149, 157, 136, 137, 138, 148, 147})
		case TypeQua17:
			chk.Ints(tst, "Qua17: c0.V", mesh.Cells[0].V, []int{0, 4, 32, 28, 1, 14, 31, 24, 2, 19, 30, 17, 3, 25, 29, 13, 18})
			chk.Ints(tst, "Qua17: cn.V", mesh.Cells[n].V, []int{120, 124, 152, 148, 121, 128, 151, 138, 122, 135, 150, 133, 123, 139, 149, 127, 134})
		}

		// compute tags map
		tm, err := mesh.GetTagMaps()
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// vertices along left border
		vset := tm.VertTag2verts[4]
		vset = append(vset, tm.VertTag2verts[41]...)
		vset = append(vset, tm.VertTag2verts[34]...)
		for _, v := range vset {
			rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
			chk.Scalar(tst, "r", 1e-15, r, rm)
		}

		// vertices along right border
		vset = tm.VertTag2verts[2]
		vset = append(vset, tm.VertTag2verts[12]...)
		vset = append(vset, tm.VertTag2verts[23]...)
		for _, v := range vset {
			Rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
			chk.Scalar(tst, "R", 1e-15, R, Rm)
		}

		// plot
		if chk.Verbose {
			plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
			mesh.Draw(args)
			plt.HideAllBorders()
			plt.Save("/tmp/gosl/gm", io.Sf("ring-%s", TypeIndexToKey[ctype]))
		}
	}
}

func TestGen04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Gen04. 2d ring (higher-level)")

	// constants
	alpha := math.Pi / 2.0
	r, R := 1.0, 3.0
	ndivR, ndivA := 3, 5

	// plot arguments
	args := NewArgs()
	args.WithEdges = true
	args.WithCells = true
	args.WithVerts = true

	// generate many quads
	ctypes := []int{TypeQua4, TypeQua8, TypeQua9, TypeQua12, TypeQua16, TypeQua17}
	for _, ctype := range ctypes {
		mesh, err := GenRing2d(ctype, ndivR, ndivA, r, R, alpha)
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// check connectivity
		n := ndivR*ndivA - 1
		switch ctype {
		case TypeQua4:
			chk.Ints(tst, "Qua4:  c0.V", mesh.Cells[0].V, []int{0, 1, 5, 4})
			chk.Ints(tst, "Qua4:  cn.V", mesh.Cells[n].V, []int{18, 19, 23, 22})
		case TypeQua8:
			chk.Ints(tst, "Qua8:  c0.V", mesh.Cells[0].V, []int{0, 2, 13, 11, 1, 8, 12, 7})
			chk.Ints(tst, "Qua8:  cn.V", mesh.Cells[n].V, []int{48, 50, 61, 59, 49, 54, 60, 53})
		case TypeQua9:
			chk.Ints(tst, "Qua9:  c0.V", mesh.Cells[0].V, []int{0, 2, 16, 14, 1, 9, 15, 7, 8})
			chk.Ints(tst, "Qua9:  cn.V", mesh.Cells[n].V, []int{60, 62, 76, 74, 61, 69, 75, 67, 68})
		case TypeQua12:
			chk.Ints(tst, "Qua12: c0.V", mesh.Cells[0].V, []int{0, 3, 21, 18, 1, 11, 20, 14, 2, 15, 19, 10})
			chk.Ints(tst, "Qua12: cn.V", mesh.Cells[n].V, []int{78, 81, 99, 96, 79, 85, 98, 88, 80, 89, 97, 84})
		case TypeQua16:
			chk.Ints(tst, "Qua16: c0.V", mesh.Cells[0].V, []int{0, 3, 33, 30, 1, 13, 32, 20, 2, 23, 31, 10, 11, 12, 22, 21})
			chk.Ints(tst, "Qua16: cn.V", mesh.Cells[n].V, []int{126, 129, 159, 156, 127, 139, 158, 146, 128, 149, 157, 136, 137, 138, 148, 147})
		case TypeQua17:
			chk.Ints(tst, "Qua17: c0.V", mesh.Cells[0].V, []int{0, 4, 32, 28, 1, 14, 31, 24, 2, 19, 30, 17, 3, 25, 29, 13, 18})
			chk.Ints(tst, "Qua17: cn.V", mesh.Cells[n].V, []int{120, 124, 152, 148, 121, 128, 151, 138, 122, 135, 150, 133, 123, 139, 149, 127, 134})
		}

		// compute tags map
		tm, err := mesh.GetTagMaps()
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// vertices along left border
		vset := tm.VertTag2verts[4]
		vset = append(vset, tm.VertTag2verts[41]...)
		vset = append(vset, tm.VertTag2verts[34]...)
		for _, v := range vset {
			rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
			chk.Scalar(tst, "r", 1e-15, r, rm)
		}

		// vertices along right border
		vset = tm.VertTag2verts[2]
		vset = append(vset, tm.VertTag2verts[12]...)
		vset = append(vset, tm.VertTag2verts[23]...)
		for _, v := range vset {
			Rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
			chk.Scalar(tst, "R", 1e-15, R, Rm)
		}

		// plot
		if chk.Verbose {
			plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
			mesh.Draw(args)
			plt.HideAllBorders()
			plt.Save("/tmp/gosl/gm", io.Sf("ring-%s-hl", TypeIndexToKey[ctype]))
		}
	}
}

func TestGen05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Gen05. 2d ring (full ring)")

	// constants
	alpha := 2.0 * math.Pi
	r, R := 1.0, 3.0
	ndivR, ndivA := 2, 3

	// plot arguments
	args := NewArgs()
	args.WithEdges = true
	args.WithCells = true
	args.WithVerts = true
	args.WithIdsVerts = true

	// generate many quads
	ctypes := []int{TypeQua4, TypeQua8, TypeQua9, TypeQua12, TypeQua16, TypeQua17}
	for _, ctype := range ctypes {
		mesh, err := GenRing2d(ctype, ndivR, ndivA, r, R, alpha)
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// check connectivity
		n := ndivR*ndivA - 1
		switch ctype {
		case TypeQua4:
			chk.Ints(tst, "Qua4:  c0.V", mesh.Cells[0].V, []int{0, 1, 4, 3})
			chk.Ints(tst, "Qua4:  cn.V", mesh.Cells[n].V, []int{7, 8, 2, 1})
		case TypeQua8:
			chk.Ints(tst, "Qua8:  c0.V", mesh.Cells[0].V, []int{0, 2, 10, 8, 1, 6, 9, 5})
			chk.Ints(tst, "Qua8:  cn.V", mesh.Cells[n].V, []int{18, 20, 4, 2, 19, 23, 3, 22})
		case TypeQua9:
			chk.Ints(tst, "Qua9:  c0.V", mesh.Cells[0].V, []int{0, 2, 12, 10, 1, 7, 11, 5, 6})
			chk.Ints(tst, "Qua9:  cn.V", mesh.Cells[n].V, []int{22, 24, 4, 2, 23, 29, 3, 27, 28})
		case TypeQua12:
			chk.Ints(tst, "Qua12: c0.V", mesh.Cells[0].V, []int{0, 3, 16, 13, 1, 8, 15, 10, 2, 11, 14, 7})
			chk.Ints(tst, "Qua12: cn.V", mesh.Cells[n].V, []int{29, 32, 6, 3, 30, 35, 5, 37, 31, 38, 4, 34})
		case TypeQua16:
			chk.Ints(tst, "Qua16: c0.V", mesh.Cells[0].V, []int{0, 3, 24, 21, 1, 10, 23, 14, 2, 17, 22, 7, 8, 9, 16, 15})
			chk.Ints(tst, "Qua16: cn.V", mesh.Cells[n].V, []int{45, 48, 6, 3, 46, 55, 5, 59, 47, 62, 4, 52, 53, 54, 61, 60})
		case TypeQua17:
			chk.Ints(tst, "Qua17: c0.V", mesh.Cells[0].V, []int{0, 4, 24, 20, 1, 10, 23, 17, 2, 14, 22, 12, 3, 18, 21, 9, 13})
			chk.Ints(tst, "Qua17: cn.V", mesh.Cells[n].V, []int{44, 48, 8, 4, 45, 51, 7, 58, 46, 56, 6, 54, 47, 59, 5, 50, 55})
		}

		// compute tags map
		tm, err := mesh.GetTagMaps()
		if err != nil {
			tst.Errorf("%v", err)
			return
		}

		// vertices along left border
		vset := tm.VertTag2verts[4]
		vset = append(vset, tm.VertTag2verts[41]...)
		vset = append(vset, tm.VertTag2verts[34]...)
		for _, v := range vset {
			rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
			chk.Scalar(tst, "r", 1e-15, r, rm)
		}

		// vertices along right border
		vset = tm.VertTag2verts[2]
		vset = append(vset, tm.VertTag2verts[12]...)
		vset = append(vset, tm.VertTag2verts[23]...)
		for _, v := range vset {
			Rm := math.Sqrt(v.X[0]*v.X[0] + v.X[1]*v.X[1])
			chk.Scalar(tst, "R", 1e-15, R, Rm)
		}

		// plot
		if chk.Verbose {
			plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
			mesh.Draw(args)
			plt.HideAllBorders()
			plt.Save("/tmp/gosl/gm", io.Sf("ring-%s-full", TypeIndexToKey[ctype]))
		}
	}
}
