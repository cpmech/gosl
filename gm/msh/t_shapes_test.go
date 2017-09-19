// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func TestShp01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Shp01")

	r := []float64{0, 0, 0}

	for ctypeindex := range Functions {

		io.Pfyel("--------------------------------- %-6s---------------------------------\n", TypeIndexToKey[ctypeindex])

		// check S
		tol := 1e-17
		if ctypeindex == TypeTri10 {
			tol = 1e-14
		}
		checkShape(tst, ctypeindex, tol, chk.Verbose)

		// check dSdR
		tol = 1e-14
		if ctypeindex == TypeLin5 || ctypeindex == TypeLin4 || ctypeindex == TypeTri10 || ctypeindex == TypeQua12 || ctypeindex == TypeQua16 {
			tol = 1e-10
		}
		if ctypeindex == TypeTri15 {
			tol = 1e-9
		}
		checkDerivs(tst, ctypeindex, r, tol, chk.Verbose)

		io.PfGreen("OK\n")
	}
}

// checkShape checks that shape functions result in 1.0 @ nodes
func checkShape(tst *testing.T, ctypeindex int, tol float64, verbose bool) {

	// information
	fcn := Functions[ctypeindex]
	ndim := GeomNdim[ctypeindex]
	nverts := NumVerts[ctypeindex]
	coords := NatCoords[ctypeindex]

	// allocate slices
	S := la.NewVector(nverts)
	dSdR := la.NewMatrix(nverts, ndim)

	// loop over all vertices
	errS := 0.0
	r := []float64{0, 0, 0}
	for n := 0; n < nverts; n++ {

		// natural coordinates @ vertex
		for i := 0; i < ndim; i++ {
			r[i] = coords[i][n]
		}

		// compute function
		fcn(S, dSdR, r, false)

		// check
		if verbose {
			for _, val := range S {
				if math.Abs(val) < 1e-15 {
					val = 0
				}
				io.Pf("%3v", val)
			}
			io.Pf("\n")
		}
		for m := 0; m < nverts; m++ {
			if n == m {
				errS += math.Abs(S[m] - 1.0)
			} else {
				errS += math.Abs(S[m])
			}
		}
	}

	// error
	if errS > tol {
		tst.Errorf("%s failed with err = %g\n", TypeIndexToKey[ctypeindex], errS)
		return
	}
}

// checkDerivs checks dSdR derivatives of shape structures
func checkDerivs(tst *testing.T, ctypeindex int, r []float64, tol float64, verbose bool) {

	// information
	fcn := Functions[ctypeindex]
	ndim := GeomNdim[ctypeindex]
	nverts := NumVerts[ctypeindex]

	// allocate slices
	S := la.NewVector(nverts)
	dSdR := la.NewMatrix(nverts, ndim)

	// analytical
	fcn(S, dSdR, r, true)

	// check
	M := dSdR.GetDeep2()
	chk.DerivVecVec(tst, "dSdR", tol, M, r[:ndim], 1e-1, verbose, func(f, x []float64) {
		fcn(f, nil, x, false)
	})
}
