// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	//SAVE_FIG = true
	SAVE_FIG = false
)

func Test_noncteM01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("noncteM01")

	prms := []string{"φ", "Mfix"}
	vals := []float64{32, 0}
	var o NcteM
	o.Init(prms, vals)

	// check
	if math.Abs(o.M(1)-o.Mcs) > 1e-17 {
		chk.Panic("M(+1) failed. err = %v", o.M(1)-o.Mcs)
	}
	if o.Mfix {
		if math.Abs(o.M(-1)-o.Mcs) > 1e-17 {
			chk.Panic("M(-1) failed. err = %v", o.M(-1)-o.Mcs)
		}
	} else {
		Mext := 6.0 * math.Sin(32*math.Pi/180) / (3 + math.Sin(32*math.Pi/180))
		if math.Abs(o.M(-1)-Mext) > 1e-15 {
			chk.Panic("M(-1) failed. err = %v", o.M(-1)-Mext)
		}
	}

	ver, tol := chk.Verbose, 1e-9
	for _, w := range utl.LinSpace(-1, 1, 11) {
		dMdw := o.DMdw(w)
		d2Mdw2 := o.D2Mdw2(w)
		chk.DerivScaSca(tst, "dM/dw  ", tol, dMdw, w, 1e-1, ver, func(x float64) (float64, error) {
			return o.M(x), nil
		})
		chk.DerivScaSca(tst, "dM²/dw²", tol, d2Mdw2, w, 1e-1, ver, func(x float64) (float64, error) {
			return o.DMdw(x), nil
		})
	}

	ver, tol = chk.Verbose, 1e-9
	nd := test_nd
	for m := 0; m < len(test_nd)-3; m++ {
		//for m := 0; m < 3; m++ {
		A := test_AA[m]
		σ := M_Alloc2(nd[m])
		Ten2Man(σ, A)
		s := M_Alloc2(nd[m])
		dMdσ := M_Alloc2(nd[m])
		d2Mdσdσ := M_Alloc4(nd[m])
		p, q, w := M_pqws(s, σ)
		o.Deriv2(d2Mdσdσ, dMdσ, σ, s, p, q, w)
		io.Pforan("σ = %v\n", σ)
		io.Pforan("tr(dMdσ) = %v\n", M_Tr(dMdσ))
		if math.Abs(M_Tr(dMdσ)) > 1e-16 {
			chk.Panic("tr(dMdσ)=%v failed", M_Tr(dMdσ))
		}
		I_dc_d2Mdσdσ := M_Alloc2(nd[m]) // I:d²M/dσdσ
		for j := 0; j < len(σ); j++ {
			for k := 0; k < len(σ); k++ {
				I_dc_d2Mdσdσ[j] += Im[k] * d2Mdσdσ[k][j]
			}
		}
		//io.Pfblue2("I_dc_d2Mdσdσ = %v\n", I_dc_d2Mdσdσ)
		chk.Vector(tst, "I_dc_d2Mdσdσ", 1e-15, I_dc_d2Mdσdσ, nil)

		// check dMdσ
		chk.DerivScaVec(tst, "dM/dσ", tol, dMdσ, σ, 1e-1, ver, func(x []float64) (float64, error) {
			return o.M(M_w(x)), nil
		})

		// check d²Mdσdσ
		s_tmp := M_Alloc2(nd[m])
		chk.DerivVecVec(tst, "d²M/dσdσ", tol, d2Mdσdσ, σ, 1e-1, ver, func(f, x []float64) error {
			p_tmp, q_tmp, w_tmp := M_pqws(s_tmp, x)
			o.Deriv1(f, x, s_tmp, p_tmp, q_tmp, w_tmp) // f := dMdσ
			return nil
		})
	}
}

func Test_Mw02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Mw02")

	prms := []string{"φ", "Mfix"}
	vals := []float64{32, 0}
	var o NcteM
	o.Init(prms, vals)

	if SAVE_FIG {
		// rosette
		full, ref := false, true
		r := 1.1 * SQ2 * o.M(1) / 3.0
		PlotRosette(r, full, ref, true, 7)

		// NcteM
		npts := 201
		X := make([]float64, npts)
		Y := make([]float64, npts)
		W := utl.LinSpace(-1, 1, npts)
		for i, w := range W {
			θ := math.Asin(w) / 3.0
			r := SQ2 * o.M(w) / 3.0
			X[i] = -r * math.Sin(math.Pi/6.0-θ)
			Y[i] = r * math.Cos(math.Pi/6.0-θ)
			//plt.Text(X[i], Y[i], io.Sf("$\\\\theta=%.2f$", θ*180.0/math.Pi), "size=8, ha='center', color='red'")
			//plt.Text(X[i], Y[i], io.Sf("$w=%.2f$", w), "size=8, ha='center', color='red'")
		}
		plt.Plot(X, Y, &plt.A{C: "b"})

		// MC
		g := func(θ float64) float64 {
			return SQ2 * o.Sinφ / (SQ3*math.Cos(θ) - o.Sinφ*math.Sin(θ))
		}
		io.Pforan("M( 1) = %v\n", SQ2*o.M(1)/3.0)
		io.Pforan("g(30) = %v\n", g(math.Pi/6.0))
		for i, w := range W {
			θ := math.Asin(w) / 3.0
			r := g(θ)
			X[i] = -r * math.Sin(math.Pi/6.0-θ)
			Y[i] = r * math.Cos(math.Pi/6.0-θ)
		}
		plt.Plot(X, Y, &plt.A{C: "k"})

		// save
		plt.Equal()
		plt.Save("/tmp/gosl", "mw02")
	}
}
