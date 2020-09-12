// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"gosl/io"
)

// Sinusoid implements the sinusoid equation:
//
//   y(t) = A0 + C1⋅cos(ω0⋅t + θ)             [essential-form]
//
//   y(t) = A0 + A1⋅cos(ω0⋅t) + B1⋅sin(ω0⋅t)  [basis-form]
//
//   A1 =  C1⋅cos(θ)
//   B1 = -C1⋅sin(θ)
//   θ  = arctan(-B1 / A1)   if A1<0, θ += π
//   C1 = sqrt(A1² + B1²)
//
type Sinusoid struct {

	// input
	Period     float64 // T: period; e.g. [s]
	MeanValue  float64 // A0: mean value; e.g. [m]
	Amplitude  float64 // C1: amplitude; e.g. [m]
	PhaseShift float64 // θ: phase shift; e.g. [rad]

	// derived
	Frequency   float64 // f: frequency; e.g. [Hz] or [1 cycle/s]
	AngularFreq float64 // ω0 = 2⋅π⋅f: angular frequency; e.g. [rad⋅s⁻¹]
	TimeShift   float64 // ts = θ / ω0: time shift; e.g. [s]

	// derived: coefficients
	A []float64 // A0, A1, A2, ... (if series mode)
	B []float64 // B0, B1, B2, ... (if series mode)
}

// NewSinusoidEssential creates a new Sinusoid object with the "essential" parameters set
//   T  -- period; e.g. [s]
//   A0 -- mean value; e.g. [m]
//   C1 -- amplitude; e.g. [m]
//   θ  -- phase shift; e.g. [rad]
func NewSinusoidEssential(T, A0, C1, θ float64) (o *Sinusoid) {

	// input
	o = new(Sinusoid)
	o.Period = T
	o.MeanValue = A0
	o.Amplitude = C1
	o.PhaseShift = θ

	// derived
	o.Frequency = 1.0 / T
	o.AngularFreq = 2.0 * math.Pi * o.Frequency
	o.TimeShift = θ / o.AngularFreq

	// derived: coefficients
	A1 := +C1 * math.Cos(θ)
	B1 := -C1 * math.Sin(θ)
	o.A = []float64{A0, A1}
	o.B = []float64{+0, B1}
	return
}

// NewSinusoidBasis creates a new Sinusoid object with the "basis" parameters set
//   T  -- period; e.g. [s]
//   A0 -- mean value; e.g. [m]
//   A1 -- coefficient of the cos term
//   B1 -- coefficient of the sin term
func NewSinusoidBasis(T, A0, A1, B1 float64) (o *Sinusoid) {

	// coefficients
	C1 := math.Sqrt(A1*A1 + B1*B1)
	θ := math.Atan2(-B1, A1)

	// input
	o = new(Sinusoid)
	o.Period = T
	o.MeanValue = A0
	o.Amplitude = C1
	o.PhaseShift = θ

	// derived
	o.Frequency = 1.0 / T
	o.AngularFreq = 2.0 * math.Pi * o.Frequency
	o.TimeShift = θ / o.AngularFreq

	// derived: coefficients
	o.A = []float64{A0, A1}
	o.B = []float64{+0, B1}
	return
}

// Yessen computes y(t) = A0 + C1⋅cos(ω0⋅t + θ [essential-form]
func (o *Sinusoid) Yessen(t float64) float64 {
	ω0 := o.AngularFreq
	return o.MeanValue + o.Amplitude*math.Cos(ω0*t+o.PhaseShift)
}

// Ybasis computes y(t) = A0 + A1⋅cos(ω0⋅t) + B1⋅sin(ω0⋅t) [basis-form]
func (o *Sinusoid) Ybasis(t float64) (res float64) {
	ω0 := o.AngularFreq
	res = o.A[0]
	for i := 1; i < len(o.A); i++ {
		k := float64(i)
		res += o.A[i]*math.Cos(k*ω0*t) + o.B[i]*math.Sin(k*ω0*t)
	}
	return
}

// ApproxSquareFourier approximates sinusoid using Fourier series with N terms
func (o *Sinusoid) ApproxSquareFourier(N int) {
	o.A = make([]float64, 1+N) // all even a's are 0 (except A0)
	o.B = make([]float64, 1+N) // all b's are 0
	o.A[0] = o.MeanValue
	for k := 1; k <= N; k += 4 {
		o.A[k] = 4.0 / (float64(k) * math.Pi)
	}
	for k := 3; k <= N; k += 4 {
		o.A[k] = -4.0 / (float64(k) * math.Pi)
	}
}

// TestPeriodicity tests that f(t) = f(T + t)
func (o *Sinusoid) TestPeriodicity(tmin, tmax float64, npts int) bool {
	dt := (tmax - tmin) / float64(npts-1)
	for i := 0; i < npts; i++ {
		t := tmin + float64(i)*dt
		ya := o.Yessen(t)
		yb := o.Yessen(t + o.Period)
		if math.Abs(ya-yb) > 1e-14 {
			io.PfRed("error = %v\n", math.Abs(ya-yb))
			return false
		}
		ya = o.Ybasis(t)
		yb = o.Ybasis(t + o.Period)
		if math.Abs(ya-yb) > 1e-14 {
			io.PfRed("error = %v\n", math.Abs(ya-yb))
			return false
		}
	}
	return true
}
