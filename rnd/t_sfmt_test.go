// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func Test_list01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("sfmt01. integers")

	Init(1234)

	nints := 10
	irange := utl.IntRange(nints) // integers; e.g. 0,1,2,3,4,5,6,7,8,9
	ifreqs := make([]int, nints)  // frequencies of each integer

	labels := make([]string, nints)
	for i := 0; i < nints; i++ {
		labels[i] = io.Sf("%3d", irange[i])
	}

	nsamples := 1000
	for i := 0; i < nsamples; i++ {
		gen := IntRand(0, nints-1)
		for j, val := range irange {
			if gen == val {
				ifreqs[j]++
				break
			}
		}
	}

	io.Pf(TextHist(labels, ifreqs, 60))
}
