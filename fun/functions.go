// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import "math"

// SuqCos implements the superquadric auxiliary function that uses cos(x)
func SuqCos(angle, expon float64) float64 {
	return Sign(math.Cos(angle)) * math.Pow(math.Abs(math.Cos(angle)), expon)
}

// SuqSin implements the superquadric auxiliary function that uses sin(x)
func SuqSin(angle, expon float64) float64 {
	return Sign(math.Sin(angle)) * math.Pow(math.Abs(math.Sin(angle)), expon)
}

// Atan2p implements a positive version of atan2, in such a way that: 0 ≤ α ≤ 2π
func Atan2p(y, x float64) (αrad float64) {
	αrad = math.Atan2(y, x)
	if αrad < 0.0 {
		αrad += 2.0 * math.Pi
	}
	return
}

// Atan2pDeg implements a positive version of atan2, in such a way that: 0 ≤ α ≤ 360
func Atan2pDeg(y, x float64) (αdeg float64) {
	αdeg = math.Atan2(y, x) * 180.0 / math.Pi
	if αdeg < 0.0 {
		αdeg += 360.0
	}
	return
}

// Ramp function => MacAulay brackets
func Ramp(x float64) float64 {
	if x < 0.0 {
		return 0.0
	}
	return x
}

// Heaviside step function (== derivative of Ramp(x))
func Heav(x float64) float64 {
	if x < 0.0 {
		return 0.0
	}
	if x > 0.0 {
		return 1.0
	}
	return 0.5
}

// Sign function
func Sign(x float64) float64 {
	if x < 0.0 {
		return -1.0
	}
	if x > 0.0 {
		return 1.0
	}
	return 0.0
}

// Sramp implements a smooth ramp function. Ramp
func Sramp(x, β float64) float64 {
	if -β*x > 500.0 {
		return 0.0
	}
	return x + math.Log(1.0+math.Exp(-β*x))/β
}

// SrampD1 returns the first derivative of Sramp
func SrampD1(x, β float64) float64 {
	if -β*x > 500.0 {
		return 0.0
	}
	return 1.0 / (1.0 + math.Exp(-β*x))
}

// SrampD2 returns the second derivative of Sramp
func SrampD2(x, β float64) float64 {
	if β*x > 500.0 {
		return 0.0
	}
	return β * math.Exp(β*x) / math.Pow(math.Exp(β*x)+1.0, 2.0)
}

// Sabs implements a smooth abs function: Sabs(x) = x*x / (sign(x)*x + eps)
func Sabs(x, eps float64) float64 {
	s := 0.0
	if x > 0.0 {
		s = 1.0
	} else if x < 0.0 {
		s = -1.0
	}
	return x * x / (s*x + eps)
}

// SabsD1 returns the first derivative of Sabs
func SabsD1(x, eps float64) float64 {
	s := 0.0
	if x > 0.0 {
		s = 1.0
	} else if x < 0.0 {
		s = -1.0
	}
	d := s*x + eps
	y := x * x / d
	return (2.0*x - s*y) / d
}

// SabsD2 returns the first derivative of Sabs
func SabsD2(x, eps float64) float64 {
	s := 0.0
	if x > 0.0 {
		s = 1.0
	} else if x < 0.0 {
		s = -1.0
	}
	d := s*x + eps
	y := x * x / d
	dydt := (2.0*x - s*y) / d
	return 2.0 * (1.0 - s*dydt) / d
}
