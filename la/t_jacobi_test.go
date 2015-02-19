// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_jacobi01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("jacobi01")

	A := [][]float64{
		{1, 2, 3},
		{2, 3, 2},
		{3, 2, 2},
	}
	Q := MatAlloc(3, 3)
	v := make([]float64, 3)
	nit, err := Jacobi(Q, v, A)
	if err != nil {
		chk.Panic("Jacobi failed:\n%v", err)
	}
	io.Pforan("number of iterations = %v\n", nit)
	PrintMat("A", A, "%13.8f", false)
	PrintMat("Q", Q, "%13.8f", false)
	PrintVec("v", v, "%13.8f", false)

	chk.Matrix(tst, "Q", 1e-17, Q, [][]float64{
		{7.81993314738381295e-01, 5.26633230856907386e-01, 3.33382506832158143e-01},
		{-7.14394870018381645e-02, 6.07084171793832561e-01, -7.91419742017035133e-01},
		{-6.19179178753124115e-01, 5.95068272145819699e-01, 5.12358171676802088e-01},
	})
	chk.Vector(tst, "v", 1e-17, v, []float64{-1.55809924785903786e+00, 6.69537390404459476e+00, 8.62725343814443657e-01})
}
