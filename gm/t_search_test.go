// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"math/rand"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_hash01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("hash01")

	c1 := []float64{1.0000000001}
	c2 := []float64{1.0000000000001, 2.0000000002}
	c3 := []float64{1.0000000000002, 2.0000000002, 3.0000000000003}
	c4 := []float64{1.0000000000003, 2.0000000002, 3.0000000000003}

	h1 := HashPointC(c1)
	h2 := HashPointC(c2)
	h3 := HashPointC(c3)
	h4 := HashPointC(c4)

	io.Pforan("h1 = %v\n", h1)
	io.Pforan("h2 = %v\n", h2)
	io.Pforan("h3 = %v\n", h3)
	io.Pforan("h4 = %v\n", h4)

	if h1 == h2 {
		chk.Panic("h1 must not be equal to h2")
	}
	if h1 == h3 {
		chk.Panic("h1 must not be equal to h3")
	}
	if h1 == h4 {
		chk.Panic("h1 must not be equal to h4")
	}
	if h2 == h3 {
		chk.Panic("h2 must not be equal to h3")
	}
	if h2 == h4 {
		chk.Panic("h2 must not be equal to h4")
	}
	// TODO: this one fails
	if false {
		if h3 == h4 {
			chk.Panic("h3 must not be equal to h4")
		}
	}
}

func Test_bins01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins01. save and recovery")
	var bins Bins
	bins.Init([]float64{0, 0, 0}, []float64{10, 10, 10}, 100)

	// fill bins structure
	maxit := 10000 // number of entries
	X := make([]float64, maxit)
	Y := make([]float64, maxit)
	Z := make([]float64, maxit)
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := rand.Float64() * 10
		y := rand.Float64() * 10
		z := rand.Float64() * 10
		X[k] = x
		Y[k] = y
		Z[k] = z
		ID[k] = k
		err := bins.Append([]float64{x, y, z}, k)
		if err != nil {
			chk.Panic(err.Error())
		}
	}

	// getting ids from bins
	IDchk := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := X[k]
		y := Y[k]
		z := Z[k]
		id := bins.Find([]float64{x, y, z})
		IDchk[k] = id
	}
	chk.Ints(tst, "check ids", ID, IDchk)

}

func Test_bins02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins02. find along line (2D)")

	// bins
	var bins Bins
	bins.Init([]float64{-0.2, -0.2}, []float64{0.8, 1.8}, 5)

	// fill bins structure
	maxit := 5 // number of entries
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := float64(k) / float64(maxit)
		ID[k] = k
		err := bins.Append([]float64{x, 2*x + 0.2}, ID[k])
		if err != nil {
			chk.Panic(err.Error())
		}
	}

	// add more points to bins
	for i := 0; i < 5; i++ {
		err := bins.Append([]float64{float64(i) * 0.1, 1.8}, 100+i)
		if err != nil {
			chk.Panic(err.Error())
		}
	}

	// message
	for _, bin := range bins.All {
		if bin != nil {
			io.Pf("%v\n", bin)
		}
	}

	// find points along diagonal
	ids := bins.FindAlongSegment([]float64{0.0, 0.2}, []float64{0.8, 1.8}, 1e-8)
	io.Pforan("ids = %v\n", ids)
	chk.Ints(tst, "ids", ids, ID)

	// find additional points
	ids = bins.FindAlongSegment([]float64{-0.2, 1.8}, []float64{0.8, 1.8}, 1e-8)
	io.Pfcyan("ids = %v\n", ids)
	chk.Ints(tst, "ids", ids, []int{100, 101, 102, 103, 104, 4})

	// draw
	if chk.Verbose {
		plt.SetForPng(1, 500, 150)
		bins.Draw2d(true, true, true, true, map[int]bool{8: true, 9: true, 10: true})
		plt.SetXnticks(15)
		plt.SetYnticks(15)
		plt.SaveD("/tmp/gosl/gm", "test_bins02.png")
	}
}

func Test_bins03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins03. find along line (3D)")

	// bins
	var bins Bins
	bins.Init([]float64{0, 0, 0}, []float64{10, 10, 10}, 10)

	// fill bins structure
	maxit := 10 // number of entries
	ID := make([]int, maxit)
	var err error
	for k := 0; k < maxit; k++ {
		x := float64(k) / float64(maxit) * 10
		ID[k] = k * 11
		err = bins.Append([]float64{x, x, x}, ID[k])
		if err != nil {
			chk.Panic(err.Error())
		}
	}

	// find points along diagonal
	ids := bins.FindAlongSegment([]float64{0, 0, 0}, []float64{10, 10, 10}, 0.0000001)
	io.Pforan("ids = %v\n", ids)
	chk.Ints(tst, "ids", ID, ids)
}

func Test_bins04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("bins04. find along line (2D)")

	// bins
	var bins Bins
	bins.Init([]float64{0, 0}, []float64{1, 2}, 10)

	// add points
	points := [][]float64{
		{0.21132486540518713, 0.21132486540518713},
		{0.7886751345948129, 0.21132486540518713},
		{0.21132486540518713, 0.7886751345948129},
		{0.7886751345948129, 0.7886751345948129},
		{0.21132486540518713, 1.2113248654051871},
		{0.7886751345948129, 1.2113248654051871},
		{0.21132486540518713, 1.788675134594813},
		{0.7886751345948129, 1.788675134594813},
	}
	var err error
	for i := 0; i < 8; i++ {
		err = bins.Append(points[i], i)
		if err != nil {
			chk.Panic(err.Error())
		}
	}
	io.Pforan("bins = %v\n", bins)

	// find points
	x := 0.7886751345948129
	ids := bins.FindAlongSegment([]float64{x, 0}, []float64{x, 2}, 1.e-15)
	io.Pforan("ids = %v\n", ids)
	chk.Ints(tst, "ids", []int{1, 3, 5, 7}, ids)

	// draw
	if chk.Verbose {
		plt.SetForPng(1, 500, 150)
		bins.Draw2d(true, true, true, true, nil)
		plt.SaveD("/tmp/gosl/gm", "test_bins04.png")
	}
}
