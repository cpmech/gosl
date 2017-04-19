// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/gm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	// input data
	fn, fnk := io.ArgToFilename(0, "nurbs01", ".msh", true)
	ctrl := io.ArgToBool(1, true)
	ids := io.ArgToBool(2, true)
	useminmax := io.ArgToBool(3, false)
	axisequal := io.ArgToBool(4, true)
	xmin := io.ArgToFloat(5, 0)
	xmax := io.ArgToFloat(6, 0)
	ymin := io.ArgToFloat(7, 0)
	ymax := io.ArgToFloat(8, 0)
	npts := io.ArgToInt(10, 41)

	// print input table
	io.Pf("\n%s\n", io.ArgsTable("INPUT ARGUMENTS",
		"mesh filename", "fn", fn,
		"show control points", "ctrl", ctrl,
		"show ids", "ids", ids,
		"use xmin,xmax,ymin,ymax", "useminmax", useminmax,
		"enforce axis.equal", "axisequal", axisequal,
		"min(x)", "xmin", xmin,
		"max(x)", "xmax", xmax,
		"min(y)", "ymin", ymin,
		"max(y)", "ymax", ymax,
		"number of divisions", "npts", npts,
	))

	// load nurbss
	B := gm.ReadMsh(fnk)

	// plot
	plt.Reset(false, nil)
	for _, b := range B {
		if ctrl {
			b.DrawCtrl2d(ids, "", "")
		}
		b.DrawElems2d(npts, ids, "", "")
	}
	if axisequal {
		plt.Equal()
	}
	if useminmax {
		plt.AxisRange(xmin, xmax, ymin, ymax)
	}
	plt.Save("/tmp/gosl", fnk)
}
