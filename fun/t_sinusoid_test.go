// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"testing"

	"gosl/chk"
	"gosl/utl"
)

func TestSinusoid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Sinusoid01. sinusoid equation 01")

	// data
	π := math.Pi // 3.14159265359...
	T := 1.5     // period [s]
	A0 := 1.7    // mean value
	C1 := 1.0    // amplitude
	θ := π / 3.0 // phase shift [rad]
	sa := NewSinusoidEssential(T, A0, C1, θ)
	sb := NewSinusoidBasis(T, A0, sa.A[1], sa.B[1])

	// check setup data
	chk.Float64(tst, "Period", 1e-15, sa.Period, sb.Period)
	chk.Float64(tst, "Frequency", 1e-15, sa.Frequency, sb.Frequency)
	chk.Float64(tst, "PhaseShift", 1e-15, sa.PhaseShift, sb.PhaseShift)
	chk.Float64(tst, "MeanValue", 1e-15, sa.MeanValue, sb.MeanValue)
	chk.Float64(tst, "Amplitude", 1e-15, sa.Amplitude, sb.Amplitude)
	chk.Float64(tst, "AngularFreq", 1e-15, sa.AngularFreq, sb.AngularFreq)
	chk.Float64(tst, "TimeShift", 1e-15, sa.TimeShift, sb.TimeShift)
	chk.Array(tst, "A", 1e-15, sa.A, sb.A)
	chk.Array(tst, "B", 1e-15, sa.B, sb.B)

	// check essen vs basis
	tt := utl.LinSpace(-sa.TimeShift, 2.5, 7)
	y1 := make([]float64, len(tt))
	y2 := make([]float64, len(tt))
	for i := 0; i < len(tt); i++ {
		y1[i] = sa.Yessen(tt[i])
		y2[i] = sb.Ybasis(tt[i])
	}
	chk.Array(tst, "essen vs basis", 1e-15, y1, y2)

	// check periodicity
	if !sa.TestPeriodicity(0, 4*π, 10) {
		chk.Panic("failed\n")
	}
}
