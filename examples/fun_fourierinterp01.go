// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// function and analytic derivative
	f := func(x float64) (float64, error) { return fun.Boxcar(x-math.Pi/2, 0, math.Pi), nil }

	plt.Reset(true, &plt.A{WidthPt: 500, Prop: 1.7})
	plt.SplotGap(0.0, 0.2)

	plt.Subplot(3, 1, 1)
	plt.Title("f(x) and interpolation", &plt.A{Fsz: 9})

	plt.Subplot(3, 1, 2)
	plt.Title("df/dx(x) and derivative of interpolation", &plt.A{Fsz: 9})

	plt.Subplot(3, 1, 3)
	plt.Title("d2f/dx2(x) and second deriv interpolation", &plt.A{Fsz: 9})

	for k, p := range []uint64{2, 3, 4, 5} {

		N := 2 << p
		fou, err := fun.NewFourierInterp(N, fun.SmoLanczosKind)
		chk.EP(err)
		err = fou.CalcA(f)
		chk.EP(err)

		ff := f
		ll := ""
		if k == 2 {
			ff = nil
			ll = "Lanczos"
		}
		l := io.Sf("%d", N)
		fou.Plot(3, 3, ff, nil, nil, &plt.A{C: "k", L: ""}, &plt.A{C: plt.C(k, 0), L: l}, &plt.A{C: plt.C(k, 0), L: ll}, &plt.A{C: plt.C(k, 0), L: l})
	}

	for k, p := range []uint64{2, 3, 4, 5} {

		N := 2 << p
		fou, err := fun.NewFourierInterp(N, fun.SmoRcosKind)
		chk.EP(err)
		err = fou.CalcA(f)
		chk.EP(err)

		ll := ""
		if k == 2 {
			ll = "Rcos"
		}
		fou.Plot(3, 3, nil, nil, nil, nil, &plt.A{C: plt.C(k, 0), Ls: "--", L: ""}, &plt.A{C: plt.C(k, 0), Ls: "--", L: ll}, &plt.A{C: plt.C(k, 0), Ls: "--", L: ""})
	}

	for k, p := range []uint64{2, 3, 4, 5} {

		N := 2 << p
		fou, err := fun.NewFourierInterp(N, fun.SmoCesaroKind)
		chk.EP(err)
		err = fou.CalcA(f)
		chk.EP(err)

		ll := ""
		if k == 2 {
			ll = "Cesaro"
		}
		fou.Plot(3, 3, nil, nil, nil, nil, &plt.A{C: plt.C(k, 0), Ls: ":", L: ""}, &plt.A{C: plt.C(k, 0), Ls: ":", L: ll}, &plt.A{C: plt.C(k, 0), Ls: ":", L: ""})
	}

	plt.Save("/tmp/gosl", "fourierinterp01")
}
