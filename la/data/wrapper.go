// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package data wraps test_mat from https://people.sc.fsu.edu/~jburkardt/f_src/test_mat/test_mat.html
// This package should be used in tests only
package data

/*
#cgo LDFLAGS: -lm

#include <complex.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "test_mat.h"

static void cleanup(
    double* a01,
    double* a02,
    double* a03,
    double* a04,
    double* a05,
    double* a06,
    double* a07,
    double* a08,
    double* a09,
    double* a10,
    double* a11,
    double* a12,
    double* a13,
    double* a14,
    double* a15,
    double* a16,
    double* a17
) {
    if (a01 != NULL) { free (a01); }
    if (a02 != NULL) { free (a02); }
    if (a03 != NULL) { free (a03); }
    if (a04 != NULL) { free (a04); }
    if (a05 != NULL) { free (a05); }
    if (a06 != NULL) { free (a06); }
    if (a07 != NULL) { free (a07); }
    if (a08 != NULL) { free (a08); }
    if (a09 != NULL) { free (a09); }
    if (a10 != NULL) { free (a10); }
    if (a11 != NULL) { free (a11); }
    if (a12 != NULL) { free (a12); }
    if (a13 != NULL) { free (a13); }
    if (a14 != NULL) { free (a14); }
    if (a15 != NULL) { free (a15); }
    if (a16 != NULL) { free (a16); }
    if (a17 != NULL) { free (a17); }
}

static double get(int i, double* v) { return v[i]; }
*/
import "C"

type A123 struct {

	// size
	M, N int // size

	// matrices
	A  []float64 // matrix data
	EL []float64 // eigen data
	ER []float64 // eigen data
	AI []float64 // inverse
	P  []float64 // plu data
	LL []float64 // plu data
	UU []float64 // plu data
	Q  []float64 // qr data
	R  []float64 // qr data
	U  []float64 // svd data
	S  []float64 // svd data
	V  []float64 // svd data

	// vectors
	EV  []float64 // eigen data
	NL  []float64 // null vector
	NR  []float64 // null vector
	RHS []float64 // rhs
	SOL []float64 // solution

	// scalars
	Det float64 // determinant
}

func (o *A123) Generate() {

	var a *C.double
	var el *C.double
	var er *C.double
	var ai *C.double
	var p *C.double
	var ll *C.double
	var uu *C.double
	var q *C.double
	var r *C.double
	var u *C.double
	var s *C.double
	var v *C.double

	var ev *C.double
	var nl *C.double
	var nr *C.double
	var rhs *C.double
	var sol *C.double

	a = C.a123()
	el = C.a123_eigen_left()
	er = C.a123_eigen_right()
	ev = C.a123_eigenvalues()
	ai = C.a123_inverse()
	nl = C.a123_null_left()
	nr = C.a123_null_right()
	C.a123_plu(p, ll, uu)
	C.a123_qr(q, r)
	rhs = C.a123_rhs()
	sol = C.a123_solution()
	C.a123_svd(u, s, v)

	o.M, o.N = 3, 3
	l := o.M * o.N
	o.A = make([]float64, l)
	o.EL = make([]float64, l)
	o.ER = make([]float64, l)
	o.AI = make([]float64, l)
	o.P = make([]float64, l)
	o.LL = make([]float64, l)
	o.UU = make([]float64, l)
	o.Q = make([]float64, l)
	o.R = make([]float64, l)
	o.U = make([]float64, l)
	o.S = make([]float64, l)
	o.V = make([]float64, l)

	for k := 0; k < l; k++ {
		o.A[k] = float64(C.get(C.int(k), a))
		o.EL[k] = float64(C.get(C.int(k), el))
		o.ER[k] = float64(C.get(C.int(k), er))
		o.AI[k] = float64(C.get(C.int(k), ai))
		o.P[k] = float64(C.get(C.int(k), p))
		o.LL[k] = float64(C.get(C.int(k), ll))
		o.UU[k] = float64(C.get(C.int(k), uu))
		o.Q[k] = float64(C.get(C.int(k), q))
		o.R[k] = float64(C.get(C.int(k), r))
		o.U[k] = float64(C.get(C.int(k), u))
		o.S[k] = float64(C.get(C.int(k), s))
		o.V[k] = float64(C.get(C.int(k), v))
	}

	o.EV = make([]float64, o.M)
	o.NL = make([]float64, o.M)
	o.NR = make([]float64, o.M)
	o.RHS = make([]float64, o.M)
	o.SOL = make([]float64, o.M)

	for i := 0; i < o.M; i++ {
		o.EV[i] = float64(C.get(C.int(i), ev))
		o.NL[i] = float64(C.get(C.int(i), nl))
		o.NR[i] = float64(C.get(C.int(i), nr))
		o.RHS[i] = float64(C.get(C.int(i), rhs))
		o.SOL[i] = float64(C.get(C.int(i), sol))
	}

	C.cleanup(a, el, er, ai, p, ll, uu, q, r, u, s, v, ev, nl, nr, rhs, sol)

	o.Det = float64(C.a123_determinant())
}
