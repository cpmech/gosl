// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/num"
)

// CheckDerivs check derivatives computed by isotropic function
func (o *IsoFun) CheckDerivs(A []float64, tol, tol2, tolq, tol3 float64, ver bool, args ...interface{}) (err error) {

	// L and invariants
	chk.IntAssert(len(A), o.ncp)
	L := make([]float64, 3)
	err = M_EigenValsNum(L, A)
	if err != nil {
		return
	}
	p, q, err := GenInvs(L, o.n, o.a)
	if err != nil {
		return
	}
	io.Pforan("L     = %v\n", L)
	io.Pforan("p, q  = %v, %v\n", p, q)

	// constants
	h := 1e-6
	var has_error error

	// derivatives of callback functions ///////////////

	// df/dp, df/dq, d²f/dp², d²f/dq², d²f/dpdq
	dfdp, dfdq := o.gfcn(p, q, args...)
	d2fdp2, d2fdq2, d2fdpdq := o.hfcn(p, q, args...)
	if ver {
		io.Pfpink("\nd w.r.t invariants . . . \n")
	}

	// check df/dp
	dfdp_num, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
		return o.ffcn(x, q, args...)
	}, p, h)
	err = chk.PrintAnaNum("df/dp   ", tol, dfdp, dfdp_num, ver)
	if err != nil {
		has_error = err
	}

	// check df/dq
	dfdq_num, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
		return o.ffcn(p, x, args...)
	}, q, h)
	err = chk.PrintAnaNum("df/dq   ", tol, dfdq, dfdq_num, ver)
	if err != nil {
		has_error = err
	}

	// check d²f/dp²
	d2fdp2_num, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
		dfdp_tmp, _ := o.gfcn(x, q, args...)
		return dfdp_tmp
	}, p, h)
	err = chk.PrintAnaNum("d²f/dp² ", tol, d2fdp2, d2fdp2_num, ver)
	if err != nil {
		has_error = err
	}

	// check d²f/dq²
	d2fdq2_num, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
		_, dfdq_tmp := o.gfcn(p, x, args...)
		return dfdq_tmp
	}, q, h)
	err = chk.PrintAnaNum("d²f/dq² ", tol, d2fdq2, d2fdq2_num, ver)
	if err != nil {
		has_error = err
	}

	// check d²f/dpdq
	d2fdpdq_num, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
		dfdp_tmp, _ := o.gfcn(p, x, args...)
		return dfdp_tmp
	}, q, h)
	err = chk.PrintAnaNum("d²f/dpdq", tol, d2fdpdq, d2fdpdq_num, ver)
	if err != nil {
		has_error = err
	}

	// derivatives w.r.t eigenvalues ///////////////////

	// df/dL and d²f/dLdL
	_, err = o.Gp(L, args...)
	if err != nil {
		chk.Panic("Gp failed:\n%v", err)
	}
	err = o.HafterGp(args...)
	if err != nil {
		chk.Panic("HafterGp failed:\n%v", err)
	}
	dfdL := make([]float64, 3)
	d2pdLdL := la.MatAlloc(3, 3)
	d2qdLdL := la.MatAlloc(3, 3)
	d2fdLdL := la.MatAlloc(3, 3)
	copy(dfdL, o.DfdL)
	la.MatCopy(d2pdLdL, 1, o.d2pdLdL)
	la.MatCopy(d2qdLdL, 1, o.d2qdLdL)
	la.MatCopy(d2fdLdL, 1, o.DgdL)

	// check df/dL
	if ver {
		io.Pfpink("\ndf/dL . . . . . . . . \n")
	}
	var fval, tmp float64
	for j := 0; j < 3; j++ {
		dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
			tmp, L[j] = L[j], x
			defer func() { L[j] = tmp }()
			fval, err = o.Fp(L, args...)
			if err != nil {
				chk.Panic("Fp failed:\n%v", err)
			}
			return fval
		}, L[j], h)
		err := chk.PrintAnaNum(io.Sf("df/dL[%d]", j), tol, dfdL[j], dnum, ver)
		if err != nil {
			has_error = err
		}
	}

	// check d²p/dLdL
	if ver {
		io.Pfpink("\nd²p/dLdL . . . . . . . . \n")
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
				tmp, L[j] = L[j], x
				defer func() { L[j] = tmp }()
				_, err = o.Gp(L, args...)
				if err != nil {
					chk.Panic("Gp failed\n%v", err)
				}
				return o.dpdL[i]
			}, L[j], h)
			err := chk.PrintAnaNum(io.Sf("d²p/dL[%d]dL[%d]", i, j), tol2, d2pdLdL[i][j], dnum, ver)
			if err != nil {
				has_error = err
			}
		}
	}

	// check d²q/dLdL
	if ver {
		io.Pfpink("\nd²q/dLdL . . . . . . . . \n")
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
				tmp, L[j] = L[j], x
				defer func() { L[j] = tmp }()
				_, err = o.Gp(L, args...)
				if err != nil {
					chk.Panic("Gp failed\n%v", err)
				}
				return o.dqdL[i]
			}, L[j], h)
			err := chk.PrintAnaNum(io.Sf("d²q/dL[%d]dL[%d]", i, j), tolq, d2qdLdL[i][j], dnum, ver)
			if err != nil {
				has_error = err
			}
		}
	}

	// check d²f/dLdL
	if ver {
		io.Pfpink("\nd²f/dLdL . . . . . . . . \n")
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
				tmp, L[j] = L[j], x
				defer func() { L[j] = tmp }()
				_, err = o.Gp(L, args...)
				if err != nil {
					chk.Panic("Gp failed\n%v", err)
				}
				return o.DfdL[i]
			}, L[j], h)
			err := chk.PrintAnaNum(io.Sf("d²f/dL[%d]dL[%d]", i, j), tol2, d2fdLdL[i][j], dnum, ver)
			if err != nil {
				has_error = err
			}
		}
	}

	// derivatives w.r.t full tensor ///////////////////

	// dfdA and d²f/dAdA
	ncp := len(A)
	dfdA := make([]float64, ncp)
	d2fdAdA := la.MatAlloc(ncp, ncp)
	_, err = o.Ga(dfdA, A, args...) // also computes P, L and Acpy
	if err != nil {
		chk.Panic("Ga failed:\n%v", err)
	}
	err = o.HafterGa(d2fdAdA, args...) // also computes dPdA
	if err != nil {
		chk.Panic("HafterGa failed:\n%v", err)
	}

	// check df/dA
	if ver {
		io.Pfpink("\ndf/dA . . . . . . . . \n")
	}
	for j := 0; j < ncp; j++ {
		dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
			tmp, A[j] = A[j], x
			defer func() { A[j] = tmp }()
			fval, err = o.Fa(A, args...)
			if err != nil {
				chk.Panic("Fa failed:\n%v", err)
			}
			return fval
		}, A[j], h)
		err := chk.PrintAnaNum(io.Sf("df/dA[%d]", j), tol, dfdA[j], dnum, ver)
		if err != nil {
			has_error = err
		}
	}

	// check dP/dA
	if false {
		for k := 0; k < 3; k++ {
			if ver {
				io.Pfpink("\ndP%d/dA . . . . . . . . \n", k)
			}
			for i := 0; i < ncp; i++ {
				for j := 0; j < ncp; j++ {
					dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
						tmp, A[j] = A[j], x
						defer func() { A[j] = tmp }()
						err = M_EigenValsProjsNum(o.P, o.L, o.Acpy)
						if err != nil {
							chk.Panic("M_EigenValsProjsNum failed:\n%v", err)
						}
						return o.P[k][i]
					}, A[j], h)
					err := chk.PrintAnaNum(io.Sf("dP%d/dA[%d]dA[%d]", k, i, j), tol, o.dPdA[k][i][j], dnum, ver)
					if err != nil {
						has_error = err
					}
				}
			}
		}
	}

	// check d²f/dAdA
	if ver {
		io.Pfpink("\nd²f/dAdA . . . . . . . . \n")
	}
	dfdA_tmp := make([]float64, ncp)
	for i := 0; i < ncp; i++ {
		for j := 0; j < ncp; j++ {
			dnum, _ := num.DerivCentral(func(x float64, notused ...interface{}) (res float64) {
				tmp, A[j] = A[j], x
				defer func() { A[j] = tmp }()
				_, err = o.Ga(dfdA_tmp, A, args...)
				if err != nil {
					chk.Panic("Ga failed:\n%v", err)
				}
				return dfdA_tmp[i]
			}, A[j], h)
			dtol := tol2
			if i == j && (i == 3 || i == 4 || i == 5) {
				dtol = tol3
			}
			err := chk.PrintAnaNum(io.Sf("d²f/dA[%d]dA[%d]", i, j), dtol, d2fdAdA[i][j], dnum, ver)
			if err != nil {
				has_error = err
			}
		}
	}

	// any errors?
	if has_error != nil {
		err = has_error
	}
	return
}
