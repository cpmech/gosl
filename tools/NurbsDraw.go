// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// default input data
	fn := "nurbs01.msh"
	ctrl := true
	ids := true
	npts := 41

	// parse flags
	flag.Parse()
	if len(flag.Args()) > 0 {
		fn = flag.Arg(0)
	}
	if len(flag.Args()) > 1 {
		ctrl = io.Atob(flag.Arg(1))
	}
	if len(flag.Args()) > 2 {
		ids = io.Atob(flag.Arg(2))
	}
	if len(flag.Args()) > 3 {
		npts = io.Atoi(flag.Arg(3))
	}

	// print input data
	io.Pforan("Input data\n")
	io.Pforan("==========\n")
	io.Pfblue2("  fn   = %v\n", fn)
	io.Pfblue2("  ctrl = %v\n", ctrl)
	io.Pfblue2("  ids  = %v\n", ids)
	io.Pfblue2("  npts = %v\n", npts)

	// load nurbss
	fnk := io.FnKey(fn)
	B := gm.ReadMsh(fnk)

	// plot
	plt.SetForEps(0.75, 500)
	for _, b := range B {
		if ctrl {
			b.DrawCtrl2d(ids, "", "")
		}
		b.DrawElems2d(npts, ids, "", "")
	}
	plt.Equal()
	plt.Save(fnk + ".eps")
}
