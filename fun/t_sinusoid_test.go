// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestSinusoid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Sinusoid01. sinusoid equation 01")

	// data
	π := math.Pi // 3.14159265359...
	T := 1.5     // period [s]
	θ := π / 3.0 // phase shift [rad]
	A0 := 1.7    // mean value
	C1 := 1.0    // amplitude
	sa := NewSinusoidEssential(T, A0, C1, θ)
	sb := NewSinusoidBasis(T, A0, sa.A[1], sa.B[1])

	// check setup data
	chk.Scalar(tst, "Period", 1e-15, sa.Period, sb.Period)
	chk.Scalar(tst, "Frequency", 1e-15, sa.Frequency, sb.Frequency)
	chk.Scalar(tst, "PhaseShift", 1e-15, sa.PhaseShift, sb.PhaseShift)
	chk.Scalar(tst, "MeanValue", 1e-15, sa.MeanValue, sb.MeanValue)
	chk.Scalar(tst, "Amplitude", 1e-15, sa.Amplitude, sb.Amplitude)
	chk.Scalar(tst, "AngularFreq", 1e-15, sa.AngularFreq, sb.AngularFreq)
	chk.Scalar(tst, "TimeShift", 1e-15, sa.TimeShift, sb.TimeShift)
	chk.Vector(tst, "A", 1e-15, sa.A, sb.A)
	chk.Vector(tst, "B", 1e-15, sa.B, sb.B)

	// check essen vs basis
	tt := utl.LinSpace(-sa.TimeShift, 2.5, 7)
	y1 := make([]float64, len(tt))
	y2 := make([]float64, len(tt))
	for i := 0; i < len(tt); i++ {
		y1[i] = sa.Yessen(tt[i])
		y2[i] = sb.Ybasis(tt[i])
	}
	chk.Vector(tst, "essen vs basis", 1e-15, y1, y2)

	// check periodicity
	if !sa.TestPeriodicity(0, 4*π, 10) {
		chk.Panic("failed\n")
	}

	// plot
	if chk.Verbose {
		tt = utl.LinSpace(-sa.TimeShift, 2.5, 201)
		yy := utl.GetMapped(tt, func(t float64) float64 { return sa.Yessen(t) })
		plt.Reset(true, nil)
		plt.AxVline(-sa.TimeShift, &plt.A{C: "grey"})
		plt.AxVline(0.0, &plt.A{C: "grey"})
		plt.AxHline(1.7, &plt.A{C: "grey"})
		plt.Text(-sa.TimeShift, A0, io.Sf("%g", -sa.TimeShift), &plt.A{C: "r", Fsz: 9, Rot: 90, Ha: "center", Va: "center"})
		plt.Plot(tt, yy, &plt.A{C: "r", NoClip: true})
		plt.AxisYmin(0)
		plt.Gll("t", "y(t)", nil)
		plt.HideTRborders()
		plt.Save("/tmp/gosl/fun", "sinusoid01")
	}
}

func TestSinusoid02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Sinusoid02. Fourier approx of square wave")

	T := 1.0  // period [s]
	θ := 0.0  // phase shift [rad]
	A0 := 0.5 // mean value
	C1 := 1.0 // amplitude
	ss := NewSinusoidEssential(T, A0, C1, θ)

	square := func(t float64) float64 {
		t0 := -T / 4.0
		return A0 + C1*math.Pow(-1, math.Floor(2*(t-t0)/T))
	}

	// plot
	if chk.Verbose {
		ss.ApproxSquareFourier(5)
		tt := utl.LinSpace(-ss.Period, ss.Period, 201)
		y1 := make([]float64, len(tt))
		y2 := make([]float64, len(tt))
		y3 := make([]float64, len(tt))
		for i := 0; i < len(tt); i++ {
			y1[i] = square(tt[i])
			y2[i] = ss.Yessen(tt[i])
			y3[i] = ss.Ybasis(tt[i])
		}
		plt.Reset(true, nil)
		plt.Plot(tt, y1, &plt.A{C: "b", NoClip: true, L: "analytic"})
		plt.Plot(tt, y2, &plt.A{C: "r", NoClip: true, L: "essential"})
		plt.Plot(tt, y3, &plt.A{C: "g", NoClip: true, L: "Fourier"})
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Gll("t", "y(t)", &plt.A{LegOut: true, LegNcol: 3})
		plt.Save("/tmp/gosl/fun", "sinusoid02")
	}
}
