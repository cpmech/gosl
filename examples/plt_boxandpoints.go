// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/plt"
)

func main() {
	plt.Reset(true, nil)
	plt.Plot3dPoint(0, -1.5, -0.9, &plt.A{C: "grey", M: "o", Ms: 3000, Mec: "orange"})
	plt.Plot3dPoint(-2, 0.5, -0.9, &plt.A{C: "r", M: "*", Ms: 5000, Mec: "green", Void: true})
	plt.Plot3dPoint(2.5, 1.5, -0.9, &plt.A{C: "r", M: "s", Ms: 1000, Mec: "k", Void: false})
	plt.Box(-0.5, 1, -1, 2, -3, 0, &plt.A{A: 0.5, Lw: 3, Fc: "#5294ed", Ec: "#ffec4f", Wire: true})
	plt.Triad(1.5, "x", "y", "z", nil, nil)
	plt.Default3dView(-1, 1.5, -1.5, 2.5, -3.5, 0.5, true)
	plt.Save("/tmp/gosl", "plt_boxandpoints")
}
