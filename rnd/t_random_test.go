// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

const NSAMPLES = 1000

func Test_int01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("int01. integers")

	Init(1234)

	nints := 10
	irange := utl.IntRange(nints) // integers; e.g. 0,1,2,3,4,5,6,7,8,9
	ifreqs := make([]int, nints)  // frequencies of each integer

	labels := make([]string, nints)
	for i := 0; i < nints; i++ {
		labels[i] = io.Sf("%3d", irange[i])
	}

	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		gen := Int(0, nints-1)
		for j, val := range irange {
			if gen == val {
				ifreqs[j]++
				break
			}
		}
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	io.Pf(TextHist(labels, ifreqs, 60))
}

func Test_int02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("int02. integers")

	Init(1234)

	nints := 10
	irange := utl.IntRange(nints) // integers; e.g. 0,1,2,3,4,5,6,7,8,9
	ifreqs := make([]int, nints)  // frequencies of each integer

	labels := make([]string, nints)
	for i := 0; i < nints; i++ {
		labels[i] = io.Sf("%3d", irange[i])
	}

	t0 := time.Now()
	samples := make([]int, NSAMPLES)
	Ints(samples, 0, nints-1)
	for i := 0; i < NSAMPLES; i++ {
		for j, val := range irange {
			if samples[i] == val {
				ifreqs[j]++
				break
			}
		}
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	io.Pf(TextHist(labels, ifreqs, 60))
}

func Test_MTint01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTint01. integers (Mersenne Twister)")

	Init(1234)

	nints := 10
	irange := utl.IntRange(nints) // integers; e.g. 0,1,2,3,4,5,6,7,8,9
	ifreqs := make([]int, nints)  // frequencies of each integer

	labels := make([]string, nints)
	for i := 0; i < nints; i++ {
		labels[i] = io.Sf("%3d", irange[i])
	}

	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		gen := MTint(0, nints-1)
		for j, val := range irange {
			if gen == val {
				ifreqs[j]++
				break
			}
		}
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	io.Pf(TextHist(labels, ifreqs, 60))
}

func Test_MTint02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MTint02. integers (Mersenne Twister)")

	Init(1234)

	nints := 10
	irange := utl.IntRange(nints) // integers; e.g. 0,1,2,3,4,5,6,7,8,9
	ifreqs := make([]int, nints)  // frequencies of each integer

	labels := make([]string, nints)
	for i := 0; i < nints; i++ {
		labels[i] = io.Sf("%3d", irange[i])
	}

	t0 := time.Now()
	samples := make([]int, NSAMPLES)
	MTints(samples, 0, nints-1)
	for i := 0; i < NSAMPLES; i++ {
		for j, val := range irange {
			if samples[i] == val {
				ifreqs[j]++
				break
			}
		}
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

	io.Pf(TextHist(labels, ifreqs, 60))
}

func Test_bins01(tst *testing.T) {

	verbose()
	chk.PrintTitle("bins01")
}

func Test_MTflt01(tst *testing.T) {

	verbose()
	chk.PrintTitle("MTflt01. float64 (Mersenne Twister)")

	Init(1234)

	xmin := 10.0
	xmax := 20.0

	t0 := time.Now()
	for i := 0; i < NSAMPLES; i++ {
		gen := MTfloat64(xmin, xmax)
		io.Pforan("gen = %v\n", gen)
	}
	io.Pforan("time elapsed = %v\n", time.Now().Sub(t0))

}

func Test_flip01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("flip01")

	Init(1234)

	p := 0.5
	nsamples := 100
	ntrue := 0
	nfalse := 0
	for i := 0; i < nsamples; i++ {
		if FlipCoin(p) {
			ntrue++
		} else {
			nfalse++
		}
	}

	io.Pforan("ntrue  = %v (42)\n", ntrue)
	io.Pforan("nfalse = %v (58)\n", nfalse)
}

func Test_shuffle01(tst *testing.T) {

	verbose()
	chk.PrintTitle("shuffle01")

	Init(0)

	nums := utl.IntRange(10)
	io.Pfgreen("before = %v\n", nums)
	IntShuffle(nums)
	io.Pfcyan("after  = %v\n", nums)
}
