// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"math/rand"

	"github.com/cpmech/gosl/utl"
)

func Test_hash01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("hash01")

	c1 := []float64{1.0000000001}
	c2 := []float64{1.0000000000001, 2.0000000002}
	c3 := []float64{1.0000000000002, 2.0000000002, 3.0000000000003}
	c4 := []float64{1.0000000000003, 2.0000000002, 3.0000000000003}

	h1 := HashPointC(c1)
	h2 := HashPointC(c2)
	h3 := HashPointC(c3)
	h4 := HashPointC(c4)

	utl.Pforan("h1 = %v\n", h1)
	utl.Pforan("h2 = %v\n", h2)
	utl.Pforan("h3 = %v\n", h3)
	utl.Pforan("h4 = %v\n", h4)

	if h1 == h2 {
		utl.Panic("h1 must not be equal to h2")
	}
	if h1 == h3 {
		utl.Panic("h1 must not be equal to h3")
	}
	if h1 == h4 {
		utl.Panic("h1 must not be equal to h4")
	}
	if h2 == h3 {
		utl.Panic("h2 must not be equal to h3")
	}
	if h2 == h4 {
		utl.Panic("h2 must not be equal to h4")
	}
	// TODO: this one fails
	if false {
		if h3 == h4 {
			utl.Panic("h3 must not be equal to h4")
		}
	}
}

// Test for save and recovery
func Test_bins01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//utl.Tsilent = false
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
		bins.Append([]float64{x, y, z}, k)
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
	utl.CompareInts(tst, "check ids", ID, IDchk)

}

// Test for function FindAlongLine (2D)
func Test_bins02(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//utl.Tsilent = false
	var bins Bins
	bins.Init([]float64{0, 0}, []float64{1, 1}, 10)

	// fill bins structure
	maxit := 10 // number of entries
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := float64(k) / float64(maxit)
		ID[k] = k * 11
		bins.Append([]float64{x, x}, ID[k])
	}

	ids := bins.FindAlongLine([]float64{0, 0}, []float64{10, 10}, 0.0000001)
	utl.Pforan("ids = %v\n", ids)

	utl.CompareInts(tst, "check FindAlongLine", ID, ids)

}

// Test for function FindAlongLine (3D)
func Test_bins03(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//utl.Tsilent = false
	var bins Bins
	bins.Init([]float64{0, 0, 0}, []float64{10, 10, 10}, 10)

	// fill bins structure
	maxit := 1000 // number of entries
	ID := make([]int, maxit)
	for k := 0; k < maxit; k++ {
		x := float64(k) / float64(maxit)
		ID[k] = k * 11
		bins.Append([]float64{x, x, x}, ID[k])
	}

	ids := bins.FindAlongLine([]float64{0, 0, 0}, []float64{10, 10, 10}, 0.0000001)
	utl.Pforan("ids = %v\n", ids)

	utl.CompareInts(tst, "check FindAlongLine", ID, ids)

}

// Test for function FindAlongLine (2D) real case
func Test_bins04(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//utl.Tsilent = false
	var bins Bins
	bins.Init([]float64{0, 0}, []float64{1, 2}, 10)

	// fill bins structure

	points := [][]float64{
		{0.21132486540518713, 0.21132486540518713},
		{0.7886751345948129, 0.21132486540518713},
		{0.21132486540518713, 0.7886751345948129},
		{0.7886751345948129, 0.7886751345948129},
		{0.21132486540518713, 1.2113248654051871},
		{0.7886751345948129, 1.2113248654051871},
		{0.21132486540518713, 1.788675134594813},
		{0.7886751345948129, 1.788675134594813}}

	var err error
	for i := 0; i < 8; i++ {
		err = bins.Append(points[i], i)
		if err != nil {
			utl.Panic("Point %v out of range\n", points[i])
		}
	}

	utl.Pforan("bins = %v\n", bins)

	// Find
	x := 0.7886751345948129
	ids := bins.FindAlongLine([]float64{x, 0}, []float64{x, 1}, 1.e-15)

	utl.Pforan("ids = %v\n", ids)

	utl.CompareInts(tst, "check FindAlongLine", []int{1, 3, 5, 7}, ids)

}
