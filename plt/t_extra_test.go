// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func TestWaterfall01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Waterfall01")

	if chk.Verbose {

		// a simple (un-normalized) Gaussian shape with amplitude A.
		fG := func(x, x0, σ, A float64) float64 {
			return A * math.Exp(-math.Pow((x-x0)/σ, 2.0))
		}

		rand.Seed(int64(time.Now().Unix()))

		σ := 0.05
		nx, nt := 401, 11
		X := utl.NonlinSpace(0, 2, nx, 4.0, true)
		T := utl.LinSpace(0, 1, nt)
		Z := utl.Alloc(nt, nx)
		for i := 0; i < nt; i++ {
			for k := 0; k < 4; k++ {
				x0 := rand.Float64() * 2
				A := rand.Float64() * 10.0
				for j := 0; j < nx; j++ {
					Z[i][j] += fG(X[j], x0, σ, A)
				}
			}
		}

		Reset(true, nil)
		Waterfall(X, T, Z, nil, nil)
		err := Save("/tmp/gosl/plt", "waterfall01")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
