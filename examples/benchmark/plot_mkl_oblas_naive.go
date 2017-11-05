// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func main() {

	nSamplesA := 1000
	nSamplesB := 100

	_, oblasA := io.ReadTable(io.Sf("/tmp/gosl/oblas-dgemm01a-%dsamples.res", nSamplesA))

	_, oblasB := io.ReadTable(io.Sf("/tmp/gosl/oblas-dgemm01b-%dsamples.res", nSamplesB))

	_, mklA := io.ReadTable(io.Sf("/tmp/gosl/mkl-dgemm01a-%dsamples.res", nSamplesA))

	_, mklB := io.ReadTable(io.Sf("/tmp/gosl/mkl-dgemm01b-%dsamples.res", nSamplesB))

	nNaive := 0
	for _, gflops := range oblasB["naiveGflops"] {
		if gflops < 1e-14 {
			break
		}
		nNaive++
	}

	plt.Reset(true, &plt.A{WidthPt: 500, Prop: 0.7})
	plt.Plot(oblasA["m"], oblasA["Gflops"], &plt.A{C: "#1549bd", L: "OpenBLAS"})
	plt.Plot(oblasA["m"], oblasA["naiveGflops"], &plt.A{C: "r", L: "Naive"})
	plt.Plot(mklA["m"], mklA["Gflops"], &plt.A{C: "g", L: "MKL"})
	plt.Title(io.Sf("Single-threaded. nSamples=%d", nSamplesA), &plt.A{Fsz: 9})
	plt.SetTicksXlist(mklA["m"])
	plt.SetTicksRotationX(45)
	plt.SetYnticks(16)
	plt.Gll("$m$", "$GFlops$", &plt.A{LegOut: false, LegNcol: 1})
	plt.Save("/tmp/gosl", "mkl-oblas-comparison-small")

	plt.Reset(true, &plt.A{WidthPt: 500, Prop: 0.7})
	plt.Plot(oblasB["m"], oblasB["Gflops"], &plt.A{C: "#1549bd", L: "OpenBLAS"})
	plt.Plot(oblasB["m"][:nNaive], oblasB["naiveGflops"][:nNaive], &plt.A{C: "r", L: "Naive"})
	plt.Plot(mklB["m"], mklB["Gflops"], &plt.A{C: "g", L: "MKL"})
	plt.Title(io.Sf("Single-threaded. nSamples=%d", nSamplesB), &plt.A{Fsz: 9})
	plt.SetTicksXlist(mklB["m"])
	plt.SetTicksRotationX(45)
	plt.SetYnticks(16)
	plt.Gll("$m$", "$GFlops$", &plt.A{LegOut: false, LegNcol: 1})
	plt.Save("/tmp/gosl", "mkl-oblas-comparison-large")
}
