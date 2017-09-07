// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
)

func main() {

	// fix seed
	rnd.Init(1111)

	// generate data
	π := math.Pi             // 3.14159265359...
	Fs := 1000.0             // Sampling frequency
	T := 1.0 / Fs            // Sampling period
	L := 1500                // Length of signal
	t := make([]float64, L)  // Time vector
	xo := make([]float64, L) // Original signal containing a 50 Hz sinusoid of amplitude 0.7 and a 120 Hz sinusoid of amplitude 1.
	xc := make([]float64, L) // Corrupted signal with zero-mean noise a std-dev of 2.
	for i := 0; i < L; i++ {
		t[i] = float64(i) * T
		xo[i] = 0.7*math.Sin(2*π*50*t[i]) + math.Sin(2*π*120*t[i])
		xc[i] = xo[i] + 2.0*rnd.Normal(0, 2)
	}

	// allocate data arrays
	oData := make([]complex128, L)
	cData := make([]complex128, L)
	for i := 0; i < L; i++ {
		oData[i] = complex(xo[i], 0)
		cData[i] = complex(xc[i], 0)
	}

	// compute the Fourier transform of original signal
	fun.Dft1d(oData, false)

	// compute the Fourier transform of corrupted signal
	fun.Dft1d(cData, false)

	// process results
	P := make([]float64, L/2+1) // single-sided spectrum of the original signal
	Q := make([]float64, L/2+1) // single-sided spectrum of the corrupted signal
	F := make([]float64, L/2+1) // frequency domain f
	for i := 0; i < L/2+1; i++ {
		P[i] = 2 * cmplx.Abs(oData[i]) / float64(L)
		Q[i] = 2 * cmplx.Abs(cData[i]) / float64(L)
		F[i] = Fs * float64(i) / float64(L)
	}

	// plot
	plt.Reset(true, &plt.A{WidthPt: 450, Dpi: 150, Prop: 1.5})

	plt.Subplot(3, 1, 1)
	plt.Plot(t[:50], xo[:50], &plt.A{C: "b", Ls: "-", L: "signal", NoClip: true})
	plt.Plot(t[:50], xc[:50], &plt.A{C: "r", Ls: "-", L: "corrupted", NoClip: true})
	plt.Gll("$t\\quad[\\mu s]$", "$x(t)$", nil)
	plt.HideTRborders()

	plt.Subplot(3, 1, 2)
	plt.AxHline(0.7, &plt.A{C: "green", Ls: "--", NoClip: true})
	plt.AxHline(1.0, &plt.A{C: "green", Ls: "--", NoClip: true})
	plt.Plot(F, P, &plt.A{C: "#0052b8"})
	plt.Gll("$f\\quad[Hz]$", "$P(f)$", nil)
	plt.HideTRborders()

	plt.Subplot(3, 1, 3)
	plt.AxHline(0.7, &plt.A{C: "green", Ls: "--", NoClip: true})
	plt.AxHline(1.0, &plt.A{C: "green", Ls: "--", NoClip: true})
	plt.Plot(F, Q, &plt.A{C: "#ed670d"})
	plt.Gll("$f\\quad[Hz]$", "$Q(f)$", nil)
	plt.HideTRborders()

	plt.Save("/tmp/gosl", "fun_fft01")
}
