// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// TextHist prints a text histogram
//  Input:
//   labels -- labels
//   freqs  -- frequencies
func TextHist(labels []string, freqs []int, barlen int) string {

	// check
	chk.IntAssert(len(labels), len(freqs))
	if len(freqs) < 2 {
		return "freqs slice is too short\n"
	}

	// scale
	fmax := freqs[0]
	lmax := 0
	Lmax := 0
	for i, f := range freqs {
		fmax = imax(fmax, f)
		lmax = imax(lmax, len(labels[i]))
		Lmax = imax(Lmax, len(io.Sf("%d", f)))
	}
	if fmax < 1 {
		return io.Sf("max frequency is too small: fmax=%d\n", fmax)
	}
	scale := float64(barlen) / float64(fmax)

	// print
	sz := io.Sf("%d", lmax+1)
	Sz := io.Sf("%d", Lmax+1)
	l := ""
	total := 0
	for i, f := range freqs {
		l += io.Sf("%"+sz+"s | %"+Sz+"d ", labels[i], f)
		n := int(float64(f) * scale)
		for j := 0; j < n; j++ {
			l += "#"
		}
		l += "\n"
		total += f
	}
	l += io.Sf("%"+sz+"s   %"+Sz+"d\n", "", total)
	return l
}
