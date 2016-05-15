// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func Test_delaunay01(tst *testing.T) {

	// header
	//verbose()
	chk.PrintTitle("delaunay01")

	// points
	X := []float64{0, 1, 1, 0, 0.5}
	Y := []float64{0, 0, 1, 1, 0.5}

	// generate
	M, err := Delaunay2d(X, Y, chk.Verbose)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}

	// derived
	err = M.CalcDerived(0)
	if err != nil {
		tst.Errorf("derived failed\n", err)
	}

	// check
	xout := make([]float64, 5)
	yout := make([]float64, 5)
	for i, v := range M.Verts {
		io.Pforan("vert %2d : coords = %v\n", v.Id, v.C)
		xout[i] = v.C[0]
		yout[i] = v.C[1]
	}
	chk.Vector(tst, "X", 1e-15, xout, X)
	chk.Vector(tst, "Y", 1e-15, yout, Y)
	for _, c := range M.Cells {
		io.Pforan("cell %2d : verts = %v\n", c.Id, c.Verts)
	}
	chk.Ints(tst, "verts of cell 0", M.Cells[0].Verts, []int{3, 0, 4})
	chk.Ints(tst, "verts of cell 1", M.Cells[1].Verts, []int{4, 1, 2})
	chk.Ints(tst, "verts of cell 2", M.Cells[2].Verts, []int{1, 4, 0})
	chk.Ints(tst, "verts of cell 3", M.Cells[3].Verts, []int{4, 2, 3})

	if chk.Verbose {
		plt.SetForEps(1, 300)
		M.Draw2d(false, false, nil, 1)
		plt.Plot(X, Y, "'ro', clip_on=0")
		plt.Gll("$x$", "$y$", "")
		plt.Equal()
		plt.SaveD("/tmp/gosl", "delaunay01.eps")
	}
}
