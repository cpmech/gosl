// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"os"

	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/ode"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// problem
	io.Pf("Hairer-Wanner VII-p2 Eq.(1.1)")
	p := ode.ProbHwEq11()

	// FwEuler
	io.Pf("\n------------ Forward-Euler ------------------\n")
	stat1, out1, err := p.Solve(ode.FwEulerKind, true, false)
	stat1.Print()
	status(err)

	// BwEuler
	io.Pf("\n------------ Backward-Euler ------------------\n")
	stat2, out2, err := p.Solve(ode.BwEulerKind, true, false)
	stat2.Print()
	status(err)

	// MoEuler
	io.Pf("\n------------ Modified-Euler ------------------\n")
	stat3, out3, err := p.Solve(ode.MoEulerKind, false, false)
	stat3.Print()
	status(err)

	// DoPri5
	io.Pf("\n------------ Dormand-Prince5 -----------------\n")
	stat4, out4, err := p.Solve(ode.DoPri5kind, false, false)
	stat4.Print()
	status(err)

	// Radau5
	io.Pf("\n------------ Radau5 --------------------------\n")
	stat5, out5, err := p.Solve(ode.Radau5kind, false, false)
	stat5.Print()
	status(err)

	// plot
	npts := 201
	plt.Reset(true, nil)
	p.Plot("FwEuler", out1, npts, true, nil, &plt.A{C: "k", M: ".", Ls: ":"})
	p.Plot("BwEuler", out2, npts, false, nil, &plt.A{C: "r", M: ".", Ls: ":"})
	p.Plot("MoEuler", out3, npts, false, nil, &plt.A{C: "c", M: "+", Ls: ":"})
	p.Plot("Dopri5", out4, npts, false, nil, &plt.A{C: "m", M: "^", Ls: "--", Ms: 3})
	p.Plot("Radau5", out5, npts, false, nil, &plt.A{C: "b", M: "s", Ls: "-", Ms: 3})
	plt.Gll("$x$", "$y$", nil)
	plt.Save("/tmp/gosl", "ode_hweq11")
}

func status(err error) {
	if err != nil {
		io.Pf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
