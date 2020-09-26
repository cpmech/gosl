// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hb

import (
	"gosl/chk"
	"gosl/io"
	"testing"
)

func TestArrayAndInts(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Array and Ints")

	io.Pf(". . . writing . . .\n")

	uSource := []float64{2.895225697183167e-07, 0.7, -1, 8.431314054288291e-10, -6.4544742997839375, -15.060440179324589, -6.454474343732561, 1.4446963710799783e-08, 0.7, -1, 8.431260528668272e-10, -6.454473969747283, -15.060439619456256, -6.454474076761063, 3.919168102695628e-08, 0.7, -1, 8.431207003048254e-10, -6.454473665192271}

	f := Create("/tmp/gosl/h5", "basic01")
	f.PutArray("/u", uSource)
	f.PutArray("/displacements/u", []float64{4, 5, 6})
	f.PutArray("/displacements/v", []float64{40, 50, 60})
	f.PutArray("/time0/ip0/a0/u", []float64{7, 8, 9})
	f.PutArray("/time1/ip0/b0/u", []float64{70, 80, 90})
	f.PutInts("/someints", []int{100, 200, 300, 400})
	f.PutInt("/data/oneint", 123)
	f.PutFloat64("/data/onef64", 123.456)
	f.Close()

	io.Pf(". . . reading . . .\n")

	g := Open("/tmp/gosl/h5", "basic01")
	u := g.GetArray("/u")
	du := g.GetArray("/displacements/u")
	dv := g.GetArray("/displacements/v")
	t0i0a0u := g.GetArray("/time0/ip0/a0/u")
	t1i0b0u := g.GetArray("/time1/ip0/b0/u")
	someints := g.GetInts("/someints")
	oneint := g.GetInt("/data/oneint")
	onef64 := g.GetFloat64("/data/onef64")
	io.Pforan("u          = %.2f\n", u)
	io.Pforan("d_u        = %v\n", du)
	io.Pforan("d_v        = %v\n", dv)
	io.Pforan("t0_i0_a0_u = %v\n", t0i0a0u)
	io.Pforan("t1_i0_b0_u = %v\n", t1i0b0u)
	chk.Array(tst, "u         ", 1e-17, u, uSource)
	chk.Array(tst, "d_u       ", 1e-17, du, []float64{4, 5, 6})
	chk.Array(tst, "d_v       ", 1e-17, dv, []float64{40, 50, 60})
	chk.Array(tst, "t0_i0_a0_u", 1e-17, t0i0a0u, []float64{7, 8, 9})
	chk.Array(tst, "t1_i0_b0_u", 1e-17, t1i0b0u, []float64{70, 80, 90})
	chk.Ints(tst, "someints", someints, []int{100, 200, 300, 400})
	chk.Int(tst, "onevalue", oneint, 123)
	chk.Float64(tst, "onef64", 1e-15, onef64, 123.456)

	io.Pf(". . . reopening file because gob requires same reading order . . .\n")
	g.Close()
	g = Open("/tmp/gosl/h5", "basic01")

	io.Pf(". . . reading again . . .\n")

	intoU := make([]float64, len(uSource))
	intoDu := make([]float64, 3)
	intoDv := make([]float64, 3)
	intoT0i0a0u := make([]float64, 3)
	intoT1i0b0u := make([]float64, 3)

	dimsU := g.ReadArray(intoU, "/u")
	dimsDu := g.ReadArray(intoDu, "/displacements/u")
	dimsDv := g.ReadArray(intoDv, "/displacements/v")
	dimsT0i0a0u := g.ReadArray(intoT0i0a0u, "/time0/ip0/a0/u")
	dimsT1i0b0u := g.ReadArray(intoT1i0b0u, "/time1/ip0/b0/u")
	g.Close()

	chk.Ints(tst, "dims: /u", dimsU, []int{len(uSource)})
	chk.Ints(tst, "dims: /displacements/u", dimsDu, []int{3})
	chk.Ints(tst, "dims: /displacements/v", dimsDv, []int{3})
	chk.Ints(tst, "dims: /time0/ip0/a0/u", dimsT0i0a0u, []int{3})
	chk.Ints(tst, "dims: /time1/ip0/b0/u", dimsT1i0b0u, []int{3})

	chk.Array(tst, "into_u         ", 1e-17, intoU, uSource)
	chk.Array(tst, "into_d_u       ", 1e-17, intoDu, []float64{4, 5, 6})
	chk.Array(tst, "into_d_v       ", 1e-17, intoDv, []float64{40, 50, 60})
	chk.Array(tst, "into_t0_i0_a0_u", 1e-17, intoT0i0a0u, []float64{7, 8, 9})
	chk.Array(tst, "into_t1_i0_b0_u", 1e-17, intoT1i0b0u, []float64{70, 80, 90})
}

func TestDeep2(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Deep2")

	io.Pf(". . . writing . . .\n")

	f := Create("/tmp/gosl/h5", "basic02")
	f.PutDeep2("/deep2/a", [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	})
	f.PutDeep2("/deep2/b", [][]float64{
		{10, 20, 30, 11},
		{40, 50, 60, 21},
		{70, 80, 90, 31},
	})
	f.PutDeep2("/deep2/c", [][]float64{
		{10, 20, 11},
		{40, 50, 21},
		{70, 80, 31},
	})
	f.Close()

	io.Pf(". . . reading . . .\n")

	g := Open("/tmp/gosl/h5", "basic02")
	a := g.GetDeep2("/deep2/a")
	b := g.GetDeep2("/deep2/b")
	c := g.GetDeep2("/deep2/c")
	io.Pforan("a = %v\n", a)
	io.Pforan("b = %v\n", b)
	io.Pforan("c = %v\n", c)
	chk.Deep2(tst, "a", 1e-17, a, [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	})
	chk.Deep2(tst, "b", 1e-17, b, [][]float64{
		{10, 20, 30, 11},
		{40, 50, 60, 21},
		{70, 80, 90, 31},
	})
	chk.Deep2(tst, "c", 1e-17, c, [][]float64{
		{10, 20, 11},
		{40, 50, 21},
		{70, 80, 31},
	})

	io.Pf(". . . reopening file because gob requires same reading order . . .\n")
	g.Close()
	g = Open("/tmp/gosl/h5", "basic02")

	io.Pf(". . . reading again . . .\n")

	m, n, araw := g.GetDeep2raw("/deep2/a")
	g.Close()

	chk.Int(tst, "m", m, 4)
	chk.Int(tst, "n", n, 3)
	chk.Array(tst, "araw", 1e-15, araw, []float64{1, 4, 7, 10, 2, 5, 8, 11, 3, 6, 9, 12})
}

func TestDeep3(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Deep3")

	io.Pf(". . . writing . . .\n")

	data := [][][]float64{
		{{1, 2, 3}, {4}, {5, 6}},
		{{7}, {8, 9}, {10, 11, 12}},
		{{-1, -2}, {-3, -4}, {-5, -6, -7}, {-8}},
	}

	f := Create("/tmp/gosl/h5", "basic03")
	f.PutDeep3("/a", data)
	f.Close()

	io.Pf(". . . reading . . .\n")

	g := Open("/tmp/gosl/h5", "basic03")
	a := g.GetDeep3("/a")
	g.Close()
	io.Pfpink("a = %v\n", a)
	chk.Deep3(tst, "a", 1e-5, a, data)
}

func TestAttributes(tst *testing.T) {

	// verbose()
	chk.PrintTitle("Attributes")

	io.Pf(". . . writing . . .\n")

	f := Create("/tmp/gosl/h5", "basic04")
	f.SetStringAttribute("/", "summary", "simulation went well")
	f.SetIntAttribute("/", "nverts", 666)
	f.SetIntsAttribute("/", "someints", []int{111, 222, 333})
	f.Close()

	io.Pf(". . . reading . . .\n")

	g := Open("/tmp/gosl/h5", "basic04")
	res := g.GetStringAttribute("/", "summary")
	nverts := g.GetIntAttribute("/", "nverts")
	vals := g.GetIntsAttribute("/", "someints")
	g.Close()
	io.Pf("summary  = %v\n", res)
	io.Pf("nverts   = %v\n", nverts)
	io.Pf("someints = %v\n", vals)
	chk.String(tst, res, "simulation went well")
	if nverts != 666 {
		chk.Panic("error with nverts. %d != 666", nverts)
	}
	chk.Ints(tst, "someints", vals, []int{111, 222, 333})
}
