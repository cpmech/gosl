// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"

	"github.com/cpmech/gosl/io"
)

// TexPiRadFmt formats number in radians/pi notation
func TexPiRadFmt(x float64) string {
	n := int((x / (math.Pi / 6.0)) + math.Pi/12.0)
	switch n {
	case 0:
		return "$0$"
	case 1:
		return "$\\\\frac{\\\\pi}{6}$"
	case 2:
		return "$\\\\frac{\\\\pi}{3}$"
	case 3:
		return "$\\\\frac{\\\\pi}{2}$"
	case 4:
		return "$2\\\\frac{\\\\pi}{3}$"
	case 6:
		return "$\\\\pi$"
	case 8:
		return "$4\\\\frac{\\\\pi}{3}$"
	case 9:
		return "$3\\\\frac{\\\\pi}{2}$"
	case 10:
		return "$5\\\\frac{\\\\pi}{3}$"
	case 12:
		return "$2\\\\pi$"
	}
	return io.Sf("$%d\\\\frac{\\\\pi}{6}$", n)
}
