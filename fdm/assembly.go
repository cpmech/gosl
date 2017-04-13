// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdm

import (
	"github.com/cpmech/gosl/la"
)

// InitK11andK12 initialises the two triplets
func InitK11andK12(K11, K12 *la.Triplet, e *Equations) {
	K11.Init(e.N1, e.N1, e.N1*5)
	K12.Init(e.N1, e.N2, e.N1*5)
}

// AssemblePoisson2d assembles K11 and K12 corresponding to the Poisson problem in 2D
//   Solving:
//                 ∂u²        ∂u²
//            - kx ———  -  ky ———  =  s(x,y)
//                 ∂x²        ∂y²
//  Input:
//    kx and ky -- the diffusion coefficients
//    src -- the source term function s(x,y) (may be nil)
//    g -- the 2D grid
//    e -- the Equation numbers
//  Output:
//    K11, K12 and F1 are assembled (must be pre-allocated)
func AssemblePoisson2d(K11, K12 *la.Triplet, F1 []float64, kx, ky float64, src Cb_src, g *Grid2d, e *Equations) {
	K11.Start()
	K12.Start()
	la.VecFill(F1, 0.0)
	alp, bet, gam := 2.0*(kx/g.Dxx+ky/g.Dyy), -kx/g.Dxx, -ky/g.Dyy
	mol := []float64{alp, bet, bet, gam, gam}
	for i, I := range e.RF1 {
		col, row := I%g.Nx, I/g.Nx
		nodes := []int{I, I - 1, I + 1, I - g.Nx, I + g.Nx} // I, left, right, bottom, top
		if col == 0 {
			nodes[1] = nodes[2]
		}
		if col == g.Nx-1 {
			nodes[2] = nodes[1]
		}
		if row == 0 {
			nodes[3] = nodes[4]
		}
		if row == g.Ny-1 {
			nodes[4] = nodes[3]
		}
		for k, J := range nodes {
			j1, j2 := e.FR1[J], e.FR2[J] // 1 or 2?
			if j1 > -1 {                 // 11
				K11.Put(i, j1, mol[k])
			} else { // 12
				K12.Put(i, j2, mol[k])
			}
		}
		if src != nil {
			x := float64(col) * g.Dx
			y := float64(row) * g.Dy
			F1[i] += src(x, y)
		}
	}
}

// JoinVecs joins U1 and U2 by placing their components at the right place in U
func JoinVecs(U, U1, U2 []float64, e *Equations) {
	for I := 0; I < e.N; I++ {
		i1, i2 := e.FR1[I], e.FR2[I] // 1 or 2?
		if i1 > -1 {
			U[I] = U1[i1]
		} else {
			U[I] = U2[i2]
		}
	}
}
