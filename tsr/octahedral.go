// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
)

// L2O converts principal values (σ1,σ2,σ3) to octahedral values (σa,σb,σc)
func L2O(σ1, σ2, σ3 float64) (σa, σb, σc float64) {
	σa = (σ3 - σ2) / SQ2
	σb = (σ3 + σ2 - 2.0*σ1) / SQ6
	σc = -(σ1 + σ2 + σ3) / SQ3
	return
}

// O2L converts octahedral values (σa,σb,σc) to principal values (σ1,σ2,σ3)
func O2L(σa, σb, σc float64) (σ1, σ2, σ3 float64) {
	σ1 = -(σc + SQ2*σb) / SQ3
	σ2 = -(SQ2*σc - σb + SQ3*σa) / SQ6
	σ3 = -(SQ2*σc - σb - SQ3*σa) / SQ6
	return
}

// O2Lmat computes L[I:{1,2,3}][A:{a,b,c}] = dσI/dσA => σI = L * σA
func O2Lmat() (L [][]float64) {
	return [][]float64{
		{0.0, -SQ2 / SQ3, -1.0 / SQ3},
		{-1.0 / SQ2, 1.0 / SQ6, -1.0 / SQ3},
		{1.0 / SQ2, 1.0 / SQ6, -1.0 / SQ3},
	}
}

// PQW2O converts p,q,w to octahedral values (σa,σb,σc)
func PQW2O(p, q, w float64) (σa, σb, σc float64) {
	r := q * SQ2by3
	θ := math.Asin(w) / 3.0
	α := 2.0*math.Pi/3.0 - θ
	σa = r * math.Cos(α)
	σb = r * math.Sin(α)
	σc = p * SQ3
	return
}

// M_oct computes octahedral values of a 2nd order symmetric tensor with Mandel components
//  Note: the (p,q,w) stress invariants are firstly computed;
//        thus, it is more efficient to compute (p,q,w) first and then use 'PQW2O'
func M_oct(a []float64) (σa, σb, σc float64) {
	p, q, w := M_pqw(a)
	σa, σb, σc = PQW2O(p, q, w)
	return
}

// PlotRosette plots rosette in octahedral plane
func PlotRosette(r float64, full, ref bool, withtext bool, fsz float64) {
	// constants
	cr := 1.0
	cf := 0.2
	if full {
		cf = cr
	}
	l1 := []float64{0.0, cr * r}                                                      // line: 1 end points
	l2 := []float64{-cf * r * math.Cos(math.Pi/6.0), -cf * r * math.Sin(math.Pi/6.0)} // line: 2 end points
	l3 := []float64{cf * r * math.Cos(math.Pi/6.0), -cf * r * math.Sin(math.Pi/6.0)}  // line: 3 end points
	l4 := []float64{-cr * r * math.Cos(math.Pi/6.0), cr * r * math.Sin(math.Pi/6.0)}  // line: 4 = neg 1 end points
	lo := []float64{-cr * r * math.Cos(math.Pi/3.0), cr * r * math.Sin(math.Pi/3.0)}  // line: origin of cylindrical system

	// main lines
	plt.Plot([]float64{0.0, l1[0]}, []float64{0.0, l1[1]}, "'k-', color='grey', zorder=0")
	plt.Plot([]float64{0.0, l2[0]}, []float64{0.0, l2[1]}, "'k-', color='grey', zorder=0")
	plt.Plot([]float64{0.0, l3[0]}, []float64{0.0, l3[1]}, "'k-', color='grey', zorder=0")

	// reference
	plt.Plot([]float64{0.0, l4[0]}, []float64{0.0, l4[1]}, "'--', color='grey', zorder=-1")
	if full {
		plt.Plot([]float64{0.0, -l4[0]}, []float64{0.0, l4[1]}, "'--', color='grey', zorder=-1")
		plt.Plot([]float64{0.0, 0.0}, []float64{0.0, -l1[1]}, "'--', color='grey', zorder=-1")
	}
	if ref {
		plt.Plot([]float64{0.0, lo[0]}, []float64{0.0, lo[1]}, "'--', color='grey', zorder=-1")
		if full {
			plt.Plot([]float64{0.0, -lo[0]}, []float64{0.0, lo[1]}, "'--', color='grey', zorder=-1")
			plt.Plot([]float64{-cr * r, cr * r}, []float64{0.0, 0.0}, "'--', color='grey', zorder=-1")
		}
	}

	// text
	if withtext {
		plt.Text(l1[0], l1[1], "$-\\sigma_1,\\\\theta=+30^\\circ$", io.Sf("ha='center', fontsize=%g", fsz))
		plt.Text(l2[0], l2[1], "$-\\sigma_3$", io.Sf("ha='right',  fontsize=%g", fsz))
		plt.Text(l3[0], l3[1], "$-\\sigma_2$", io.Sf("ha='left',   fontsize=%g", fsz))
		plt.Text(lo[0], lo[1], "$\\\\theta=0^\\circ$", io.Sf("ha='center', fontsize=%g", fsz))
		plt.Text(l4[0], l4[1], "$\\\\theta=-30^\\circ$", io.Sf("ha='center', fontsize=%g", fsz))
	}
	plt.Equal()
}

// PlotRefOct plots reference failure criterions in octahedral plane:
//  Drucker-Prager and Mohr-Circles
func PlotRefOct(φ, σc float64, withExtCircle bool) {
	if φ > 0 {
		sφ := math.Sin(φ * math.Pi / 180.0)
		μmax := 2.0 * SQ2 * sφ / (3.0 - sφ)
		μmin := 2.0 * SQ2 * sφ / (3.0 + sφ)
		Rmax := μmax * σc
		Rmin := μmin * σc
		d30 := math.Pi / 6.0
		xa := -Rmin * math.Cos(d30)
		ya := Rmin * math.Sin(d30)
		xb := -Rmax * math.Cos(d30)
		yb := -Rmax * math.Sin(d30)
		plt.Plot([]float64{xb, xa, 0, -xa, -xb, 0, xb}, []float64{yb, ya, Rmax, ya, yb, -Rmin, yb}, "'grey', ls='-'")
		if withExtCircle {
			plt.Circle(0, 0, Rmin, "ec='grey'")
		}
		plt.Circle(0, 0, Rmax, "ec='grey'")
	}
}

type Cb_F_t func(A []float64, args ...interface{}) (fval float64, err error)
type Cb_G_t func(dfdA, A []float64, args ...interface{}) (fval float64, err error)

// PlotOct plots a function cross-section and gradients projections on octahedral plane
func PlotOct(filename string, σcCte, rmin, rmax float64, nr, nα int, φ float64, F Cb_F_t, G Cb_G_t,
	notpolarc, simplec, only0, grads, showpts, first, last bool, ferr float64, args ...interface{}) {
	A := make([]float64, 6)
	dfdA := make([]float64, 6)
	dα := 2.0 * math.Pi / float64(nα)
	dr := (rmax - rmin) / float64(nr-1)
	αmin := -math.Pi / 6.0
	if notpolarc {
		rmin = -rmax
		dr = (rmax - rmin) / float64(nr-1)
		nα = nr
		dα = dr
		αmin = rmin
	}
	x, y, f := la.MatAlloc(nα, nr), la.MatAlloc(nα, nr), la.MatAlloc(nα, nr)
	//var gx, gy, L [][]float64
	var gx, gy [][]float64
	if grads {
		gx, gy = la.MatAlloc(nα, nr), la.MatAlloc(nα, nr)
		//L      = O2Lmat()
	}
	var σa, σb []float64
	if showpts {
		σa, σb = make([]float64, nα*nr), make([]float64, nα*nr)
	}
	k := 0
	var err error
	for i := 0; i < nα; i++ {
		for j := 0; j < nr; j++ {
			α := αmin + float64(i)*dα
			r := rmin + float64(j)*dr
			if notpolarc {
				x[i][j] = α // σa
				y[i][j] = r // σb
			} else {
				x[i][j] = r * math.Cos(α) // σa
				y[i][j] = r * math.Sin(α) // σb
			}
			A[0], A[1], A[2] = O2L(x[i][j], y[i][j], σcCte)
			if showpts {
				σa[k] = x[i][j]
				σb[k] = y[i][j]
			}
			f[i][j], err = F(A, args...)
			if err != nil {
				if ferr > 0 {
					f[i][j] = ferr
				} else {
					chk.Panic(_octahedral_err1, "func", err)
				}
			}
			if grads {
				_, err = G(dfdA, A, args...)
				if err != nil {
					chk.Panic(_octahedral_err1, "grads", err)
				}
				/*
				   gx[i][j] = -o.dfdλ[0]*L[0][0] - o.dfdλ[1]*L[1][0] - o.dfdλ[2]*L[2][0]
				   gy[i][j] = -o.dfdλ[0]*L[0][1] - o.dfdλ[1]*L[1][1] - o.dfdλ[2]*L[2][1]
				*/
			}
			k += 1
		}
	}
	if first {
		plt.Reset()
		plt.SetForPng(1, 500, 125)
		PlotRosette(1.1*rmax, true, true, true, 7)
		PlotRefOct(φ, σcCte, false)
		if showpts {
			plt.Plot(σa, σb, "'k.', ms=5")
		}
	}
	if simplec {
		if !only0 {
			plt.ContourSimple(x, y, f, "")
		}
		plt.ContourSimple(x, y, f, "levels=[0], colors=['blue'], linewidths=[2]")
	} else {
		plt.Contour(x, y, f, "fsz=8")
	}
	if grads {
		plt.Quiver(x, y, gx, gy, "")
	}
	if last {
		plt.AxisOff()
		plt.Equal()
		plt.SaveD("results", filename)
	}
}

// error messages
var (
	_octahedral_err1 = "octahedral.go: OctPlotYF: %s evaluation failed:\n%v"
)
