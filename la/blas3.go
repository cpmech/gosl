// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/la/oblas"
)

// MatMatMul returns the matrix multiplication (scaled)
//
//  c := α⋅a⋅b    ⇒    cij := α * aik * bkj
//
func MatMatMul(c *Matrix, α float64, a, b *Matrix) {
	err := oblas.Dgemm(false, false, a.M, b.N, a.N, α, a.Data, a.M, b.Data, b.M, 0.0, c.Data, c.M)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}

// MatMatMulAdd returns the matrix multiplication (scaled) with addition
//
//  c += α⋅a⋅b    ⇒    cij += α * aik * bkj
//
func MatMatMulAdd(c *Matrix, α float64, a, b *Matrix) {
	err := oblas.Dgemm(false, false, a.M, b.N, a.N, α, a.Data, a.M, b.Data, b.M, 1.0, c.Data, c.M)
	if err != nil {
		chk.Panic("%v\n", err)
	}
}
