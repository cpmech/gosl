// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

func PrintIndicatorMatrix(x [][]int) (l string) {
	m, n := len(x), len(x[0])
	w := 2*n + 6 + 13
	l = io.Sf("%s i\\j |", io.StrThickLine(w))
	for j := 0; j < n; j++ {
		l += io.Sf("%2d", j)
	}
	l += io.Sf(" |        sum\n%s", io.StrThinLine(w))
	S := make([]int, n)
	for i := 0; i < m; i++ {
		l += io.Sf(" %3d |", i)
		s := 0
		for j := 0; j < n; j++ {
			l += io.Sf("%2d", x[i][j])
			s += x[i][j]
			S[j] += x[i][j]
		}
		l += io.Sf(" | Σ xij = %2d\n", s)
	}
	l += io.Sf("%s sum =", io.StrThinLine(w))
	for j := 0; j < n; j++ {
		l += io.Sf("%2d", S[j])
	}
	l += io.Sf("\n%s", io.StrThickLine(w))
	return
}

func BuildIndicatorMatrix(nv int, pth []int) (x [][]int) {
	x = utl.IntsAlloc(nv, nv)
	for k := 1; k < len(pth); k++ {
		i, j := pth[k-1], pth[k]
		x[i][j] = 1
	}
	return
}

func CheckIndicatorMatrix(source, target int, x [][]int, verbose bool) (errPath, errLoop int) {
	nv := len(x)
	var okPath, okLoop bool
	var sij, sji int
	for i := 0; i < nv; i++ {
		sij, sji = 0, 0
		for j := 0; j < nv; j++ {
			sij += x[i][j]
			sji += x[j][i]
		}
		d := sij - sji
		if i == target {
			okPath = d == -1
			okLoop = sij == 0
		} else {
			if i == source {
				okPath = d == 1
			} else {
				okPath = d == 0
			}
			okLoop = sij <= 1
		}
		if !okPath {
			errPath++
		}
		if !okLoop {
			errLoop++
		}
		if verbose {
			sok := "ok"
			if !okPath || !okLoop {
				sok = "fail"
			}
			io.Pforan("i=%2d : Σxij - Σxji = %2d  (%s)\n", i, d, sok)
		}
	}
	return
}

func CheckIndicatorMatrixRowMaj(source, target, nv int, xmat []int) (errPath, errLoop int) {
	var okPath, okLoop bool
	var sij, sji int
	for i := 0; i < nv; i++ {
		sij, sji = 0, 0
		for j := 0; j < nv; j++ {
			sij += xmat[i*nv+j]
			sji += xmat[j*nv+i]
		}
		d := sij - sji
		if i == target {
			okPath = d == -1
			okLoop = sij == 0
		} else {
			if i == source {
				okPath = d == 1
			} else {
				okPath = d == 0
			}
			okLoop = sij <= 1
		}
		if !okPath {
			errPath++
		}
		if !okLoop {
			errLoop++
		}
	}
	return
}
