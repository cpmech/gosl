// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Mmatch computes M=q/p and qy0 from c and φ corresponding to the strength that would
// be modelled by the Mohr-Coulomb model matching one of the following cones:
//  cone == "cmp" : compression cone (outer)
//       == "ext" : extension cone (inner)
//       == "psa" : plane-strain
//  Note: p, q, and M are Cambridge (conventional) quantities
func Mmatch(c, φ float64, cone string) (M, qy0 float64) {
	φr := φ * math.Pi / 180.0
	si := math.Sin(φr)
	co := math.Cos(φr)
	var ξ float64
	switch cone {
	case "cmp": // compression cone (outer)
		M = 6.0 * si / (3.0 - si)
		ξ = 6.0 * co / (3.0 - si)
	case "ext": // extension cone (inner)
		M = 6.0 * si / (3.0 + si)
		ξ = 6.0 * co / (3.0 + si)
	case "psa": // plane-strain
		t := si / co
		d := math.Sqrt(3.0 + 4.0*t*t)
		M = 3.0 * t / d
		ξ = 3.0 / d
	default:
		chk.Panic(_fcritcoef_err1, cone)
	}
	qy0 = ξ * c
	return
}

// Phi2M calculates M = max q/p at compression (φ: friction angle at compression (degrees)).
// type = {"oct", "cam", "smp"}
func Phi2M(φ float64, typ string) float64 {
	sφ := math.Sin(φ * math.Pi / 180.0)
	switch typ {
	case "oct":
		return 2.0 * SQ2 * sφ / (3.0 - sφ)
	case "cam":
		return 6.0 * sφ / (3.0 - sφ)
	default:
		chk.Panic(_fcritcoef_err2, "Phi2M", typ)
	}
	return 0
}

// M2Phi calculates φ (friction angle at compression (degrees)) given M (max q/p at compression)
func M2Phi(M float64, typ string) float64 {
	var sφ float64
	switch typ {
	case "oct":
		sφ = 3.0 * M / (M + 2.0*SQ2)
	case "cam":
		sφ = 3.0 * M / (M + 6.0)
	default:
		chk.Panic(_fcritcoef_err2, "M2Phi", typ)
	}
	return math.Asin(sφ) * 180.0 / math.Pi
}

// error messages
var (
	_fcritcoef_err1 = "fcritcoef.go: tsr.Mmatch: 'cone = %s' is incorrect"
	_fcritcoef_err2 = "fcritcoef.go: tsr.%s: 'typ = %s' is unavailable"
)
