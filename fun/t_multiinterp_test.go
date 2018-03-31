package fun

import (
	"testing"

	"github.com/cpmech/gosl/chk"
)

func TestMultiInterp01(t *testing.T) {

	chk.PrintTitle("MultiInterp01. bilinear interpolation")

	// checking f(x,y) = x^2 + 2y^2
	f := []float64{
		0.00, 0.25, 1.00, 4.00,
		2.00, 2.25, 3.00, 6.00,
		8.00, 8.25, 9.00, 12.00,
	}

	xx := []float64{0.0, 0.5, 1.0, 2.0}
	yy := []float64{0.0, 1.0, 2.0}

	o := NewBiLinear(f, xx, yy)

	for i, x := range xx {
		for j, y := range yy {
			chk.Float64(t, "P(x,y)", 1e-17, o.P(x, y), f[i+j*len(xx)])
		}
	}

	fref := []float64{1.125, 1.625, 3.2, 8.4}
	xref := []float64{0.25, 0.75, 1.2, 1.2}
	yref := []float64{0.5, 0.5, 0.8, 1.8}

	for i := 0; i < len(fref); i++ {
		chk.Float64(t, "P(xref,yref)", 1e-17, o.P(xref[i], yref[i]), fref[i])
	}

}
