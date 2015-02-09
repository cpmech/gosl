// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
    "math"
    "code.google.com/p/gosl/utl"
)

// NcteM implements Argyris-Sheng et al M(w) non-constant M coefficient
type NcteM struct {
    // input
    φ    float64
    Mfix bool
    // derived
    Sinφ float64
    Tanφ float64
    Mcs  float64
    μcs  float64
}

// Init initialises this object
func (o *NcteM) Init(prms []string, vals []float64) {
    for i, p := range prms {
        switch p {
        case "φ":    o.φ    = vals[i]
        case "Mfix": o.Mfix = vals[i] > 0
        }
    }
    φmin, φmax := 1e-7, 89.99
    if o.φ < φmin || o.φ > φmax {
        utl.Panic(_nonctem_err1, φmin, φmax, o.φ)
    }
    o.Sinφ = math.Sin(o.φ * math.Pi / 180.0)
    o.Tanφ = math.Tan(o.φ * math.Pi / 180.0)
    o.Mcs  = 6.0 * o.Sinφ / (3.0 - o.Sinφ)
    o.μcs  = math.Pow((3.0 - o.Sinφ)/(3.0 + o.Sinφ), 4.0)
}

// String returns a string representing this structure
func (o *NcteM) String() (s string) {
    return utl.Sf("  φ=%v Mcs=%v Mfix=%v\n", o.φ, o.Mcs, o.Mfix)
}

// M implements M(w)
func (o *NcteM) M(w float64) float64 {
    if o.Mfix {
        return o.Mcs
    }
    return o.Mcs * math.Pow(2.0*o.μcs / (o.μcs + 1.0 + (o.μcs - 1.0) * w), 0.25)
}

// DMdw implements dM/dw
func (o *NcteM) DMdw(w float64) float64 {
    if o.Mfix {
        return 0
    }
    return 0.25 * o.M(w) * (1.0 - o.μcs) / (o.μcs + 1.0 + (o.μcs - 1.0) * w)
}

// D2Mdw2 implements d²M/dw²
func (o *NcteM) D2Mdw2(w float64) float64 {
    if o.Mfix {
        return 0
    }
    Dcs := (o.μcs + 1.0 + (o.μcs - 1.0) * w)
    return 0.3125 * o.M(w) * (1.0 - o.μcs) * (1.0 - o.μcs) / (Dcs * Dcs)
}

// Deriv1 returns the first derivative of M w.r.t σ
//  Note: only dMdσ is output
func (o *NcteM) Deriv1(dMdσ, σ, s []float64, p, q, w float64) {
    if o.Mfix {
        for i := 0; i < len(σ); i++ {
            dMdσ[i] = 0
        }
        return
    }
    M_LodeDeriv1(dMdσ, σ, s, p, q, w) // dMdσ := dwdσ
    dMdw := o.DMdw(w)
    for i := 0; i < len(σ); i++ {
        dMdσ[i] *= dMdw
    }
}

// Deriv2 returns the first and second derivatives of M w.r.t σ
//  Note: d2Mdσdσ and dMdσ output
func (o *NcteM) Deriv2(d2Mdσdσ [][]float64, dMdσ, σ, s []float64, p, q, w float64) {
    if o.Mfix {
        for i := 0; i < len(σ); i++ {
            for j := 0; j < len(σ); j++ {
                d2Mdσdσ[i][j] = 0
            }
            dMdσ[i] = 0
        }
        return
    }
    M_LodeDeriv2(d2Mdσdσ, dMdσ, σ, s, p, q, w) // d2Mdσdσ := d2wdσdσ, dMdσ := dwdσ
    dMdw   := o.DMdw(w)
    d2Mdw2 := o.D2Mdw2(w)
    for i := 0; i < len(σ); i++ {
        for j := 0; j < len(σ); j++ {
            d2Mdσdσ[i][j] = d2Mdw2 * dMdσ[i] * dMdσ[j] + dMdw * d2Mdσdσ[i][j]
        }
    }
    for i := 0; i < len(σ); i++ {
        dMdσ[i] *= dMdw
    }
}

// error messages
var (
    _nonctem_err1 = "nonctem.go: NcteM.Init: φ [deg] must be in [%g,%g]. φ==%g is incorrect"
)
