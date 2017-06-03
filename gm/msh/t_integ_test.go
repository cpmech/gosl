// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func TestInteg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Integ01. integration of scalar function")

	// load mesh
	m, err := Read("data/mesh01.msh")
	if err != nil {
		tst.Errorf("Read failed:\n%v\n", err)
		return
	}

	o, err := NewIntegrator(4, m)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	chk.Int(tst, "Nv", o.Nv, 4)
	chk.Int(tst, "Dim", o.Dim, 2)
	chk.Matrix(tst, "X", 1e-15, o.X, [][]float64{
		{0.5, 0.5},
		{1.0, 0.5},
		{1.0, 1.0},
		{0.5, 1.0},
	})
	io.Pforan("o = %+v\n", o)

	// TODO

	if chk.Verbose && false {
		args := NewArgs()
		args.WithIdsCells = true
		args.WithIdsVerts = true
		plt.Reset(true, nil)
		m.Draw(args)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/gm", "integ01")
	}
}
