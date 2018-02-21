// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/big"

	"github.com/cpmech/gosl/chk"
)

// flags and table of values
var (
	factorialTable22  []float64   // will be filled during the first call to Factorial22
	factorialTableBig []big.Float // will be filled during the first call to FactorialBig
)

// Factorial22 implements the factorial function; i.e. computes n! up to 22!  According to [1],
// factorials up to 22! have exact double precision representations (52 bits of mantissa, not
// counting powers of two that are absorbed into the exponent)
//   References
//   [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
//        Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func Factorial22(n int) float64 {
	if n < 0 || n > 22 {
		chk.Panic("cannot use Factorial22 with n = %d", n)
	}
	if len(factorialTable22) == 0 {
		factorialTable22 = make([]float64, 23)
		factorialTable22[0] = 1.0
		for i := 1; i <= 22; i++ {
			factorialTable22[i] = factorialTable22[i-1] * float64(i)
		}
	}
	return factorialTable22[n]
}

// Factorial100 returns the factorial n! up to 100! using the math/big package
func Factorial100(n int) big.Float {
	if n < 0 || n > 100 {
		chk.Panic("cannot use Factorial100 with n = %d", n)
	}
	if len(factorialTableBig) == 0 {

		// get precision for 100!
		n100 := new(big.Float)
		n100.SetPrec(big.MaxPrec)
		_, _, err := n100.Parse("93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000", 10)
		if err != nil {
			chk.Panic("cannot parse 100!\n%v\n", err)
		}
		maxprec := n100.MinPrec()

		// compute numbers
		factorialTableBig = make([]big.Float, 101)
		factorialTableBig[0].SetFloat64(1.0)
		t := new(big.Float)
		for i := 1; i <= 100; i++ {
			t.SetInt64(int64(i))
			factorialTableBig[i].SetPrec(maxprec)
			factorialTableBig[i].Mul(&factorialTableBig[i-1], t)
		}
	}
	return factorialTableBig[n]
}

// Beta computes the beta function by calling the Lgamma function
func Beta(a, b float64) float64 {
	la, sgnla := math.Lgamma(a)
	lb, sgnlb := math.Lgamma(b)
	lc, sgnlc := math.Lgamma(a + b)
	return float64(sgnla*sgnlb*sgnlc) * math.Exp(la+lb-lc)
}

// Binomial comptues the binomial coefficient (n k)^T
func Binomial(n, k int) float64 {
	if n < 0 || k < 0 || k > n {
		chk.Panic("Binomial function requires that k <= n (both positive). Incorrect values: n=%v, k=%v", n, k)
	}
	if k == 0 || k == n {
		return 1
	}
	if k == 1 || k == n-1 {
		return float64(n)
	}
	// use fast table lookup. See [1] page 258
	if n <= 22 {
		// the floor function cleans up roundoff error for smaller values of n and k.
		return math.Floor(0.5 + Factorial22(n)/(Factorial22(k)*Factorial22(n-k)))
	}
	// use beta function
	if k > n-k {
		k = n - k // take advantage of symmetry
	}
	res := float64(k) * Beta(float64(k), float64(n-k+1))
	if res == 0 {
		chk.Panic("Binomial function failed with n=%v, k=%v", n, k)
	}
	return math.Floor(0.5 + 1.0/res)
}

// UintBinomial implements the Binomial coefficient using uint64. Panic happens on overflow
// Also, this function uses a loop so it may not be very efficient for large k
// The code below comes from https://en.wikipedia.org/wiki/Binomial_coefficient
// [cannot find a reference to cite]
func UintBinomial(n, k uint64) uint64 {
	if k > n {
		chk.Panic("UintBinomial function requires that k <= n. Incorrect values: n=%v, k=%v", n, k)
	}
	if k == 0 || k == n {
		return 1
	}
	if k == 1 || k == n-1 {
		return n
	}
	if k > n-k {
		k = n - k // take advantage of symmetry
	}
	var c uint64 = 1
	var i uint64
	for i = 1; i <= k; i, n = i+1, n-1 {
		if c/i > math.MaxUint64/n {
			chk.Panic("Overflow in UintBinomial: %v > %v", c/i, math.MaxUint64/n)
		}
		c = c/i*n + c%i*n/i // split c*n/i into (c/i*i + c%i)*n/i
	}
	return c
}

// Rbinomial computes the binomial coefficient with real (non-negative) arguments by calling the Gamma function
func Rbinomial(x, y float64) float64 {
	if x < 0 || y < 0 {
		chk.Panic("Rbinomial requires x and y to be non-negative, at this moment")
	}
	a := math.Gamma(x + 1.0)
	b := math.Gamma(y + 1.0)
	c := math.Gamma(x - y + 1.0)
	return a / (b * c)
}

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

// Heav computes the Heaviside step function (== derivative of Ramp(x))
//
//             │ 0    if x < 0
//   Heav(x) = ┤ 1/2  if x = 0
//             │ 1    if x > 0
//
func Heav(x float64) float64 {
	if x < 0.0 {
		return 0.0
	}
	if x > 0.0 {
		return 1.0
	}
	return 0.5
}

// Sign implements the sign function
//
//             │ -1   if x < 0
//   Sign(x) = ┤  0   if x = 0
//             │  1   if x > 0
//
func Sign(x float64) float64 {
	if x < 0.0 {
		return -1.0
	}
	if x > 0.0 {
		return 1.0
	}
	return 0.0
}

// Boxcar implements the boxcar function
//
//   Boxcar(x;a,b) = Heav(x-a) - Heav(x-b)
//
//                   │ 0    if x < a or  x > b
//   Boxcar(x;a,b) = ┤ 1/2  if x = a or  x = b
//                   │ 1    if x > a and x < b
//
//   Note: a ≤ x ≤ b; i.e. b ≥ a (not checked)
//
func Boxcar(x, a, b float64) float64 {
	if x < a || x > b {
		return 0
	}
	if x > a && x < b {
		return 1
	}
	return 0.5
}

// Rect implements the rectangular function
//
//   Rect(x) = Boxcar(x;-0.5,0.5)
//
//             │ 0    if |x| > 1/2
//   Rect(x) = ┤ 1/2  if |x| = 1/2
//             │ 1    if |x| < 1/2
//
func Rect(x float64) float64 {
	if x < -0.5 || x > +0.5 {
		return 0
	}
	if x > -0.5 && x < +0.5 {
		return 1
	}
	return 0.5
}

// Hat implements the hat function
//
//      --———--   o (xc,y0+h)
//         |     / \
//         h    /   \    m = h/l
//         |   /m    \
//   y0 ——————o       o—————————
//
//            |<  2l >|
//
func Hat(x, xc, y0, h, l float64) float64 {
	if x <= xc-l || x >= xc+l {
		return y0
	}
	if x <= xc {
		return y0 + (h/l)*(x-xc+l)
	}
	return y0 + h - (h/l)*(x-xc)
}

// HatD1 returns the first derivative of the hat function
// NOTE: the discontinuity is ignored ⇒ D1(xc-l)=D1(xc+l)=D1(xc)=0
func HatD1(x, xc, y0, h, l float64) float64 {
	if x <= xc-l || x >= xc+l || x == xc {
		return 0
	}
	if x < xc {
		return h / l
	}
	return -h / l
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

// ExpPix uses Euler's formula to compute exp(+i⋅x) = cos(x) + i⋅sin(x)
func ExpPix(x float64) complex128 {
	return complex(math.Cos(x), math.Sin(x))
}

// ExpMix uses Euler's formula to compute exp(-i⋅x) = cos(x) - i⋅sin(x)
func ExpMix(x float64) complex128 {
	return complex(math.Cos(x), -math.Sin(x))
}

// Sinc computes the sine cardinal (sinc) function
//
//   Sinc(x) = |     1      if x = 0
//             | sin(x)/x   otherwise
//
func Sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(x) / x
}

// NegOnePowN computes (-1)ⁿ
func NegOnePowN(n int) float64 {
	if (n & 1) == 0 { // even
		return 1
	}
	return -1
}

// ImagPowN computes iⁿ = (√-1)ⁿ
//
//   i¹ = i      i²  = -1      i³  = -i      i⁴  = 1
//   i⁵ = i      i⁶  = -1      i⁷  = -i      i⁸  = 1
//   i⁹ = i      i¹⁰ = -1      i¹¹ = -i      i¹² = 1
//
func ImagPowN(n int) complex128 {
	if n == 0 {
		return 1
	}
	switch n % 4 {
	case 1:
		return 1i
	case 2:
		return -1
	case 3:
		return -1i
	}
	return 1
}

// ImagXpowN computes (x⋅i)ⁿ
//
//   (x⋅i)¹ = x¹⋅i      (x⋅i)²  = -x²       (x⋅i)³  = -x³ ⋅i      (x⋅i)⁴  = x⁴
//   (x⋅i)⁵ = x⁵⋅i      (x⋅i)⁶  = -x⁶       (x⋅i)⁷  = -x⁷ ⋅i      (x⋅i)⁸  = x⁸
//   (x⋅i)⁹ = x⁹⋅i      (x⋅i)¹⁰ = -x¹⁰      (x⋅i)¹¹ = -x¹¹⋅i      (x⋅i)¹² = x¹²
//
func ImagXpowN(x float64, n int) complex128 {
	if n == 0 {
		return 1
	}
	xn := math.Pow(x, float64(n))
	switch n % 4 {
	case 1:
		return complex(0, xn)
	case 2:
		return complex(-xn, 0)
	case 3:
		return complex(0, -xn)
	}
	return complex(xn, 0)
}

// PowP computes real raised to positive integer xⁿ
func PowP(x float64, n uint32) (r float64) {
	if n == 0 {
		return 1.0
	}
	if n == 1 {
		return x
	}
	if n == 2 {
		return x * x
	}
	if n == 3 {
		return x * x * x
	}
	if n == 4 {
		r = x * x
		return r * r
	}
	if n == 5 {
		r = x * x
		return r * r * x
	}
	if n == 6 {
		r = x * x * x
		return r * r
	}
	if n == 7 {
		r = x * x * x
		return r * r * x
	}
	if n == 8 {
		r = x * x * x * x
		return r * r
	}
	if n == 9 {
		r = x * x * x
		return r * r * r
	}
	if n == 10 {
		r = x * x * x
		return r * r * r * x
	}
	r = 1.0
	for i := n; i > 0; i >>= 1 {
		if i&1 == 1 {
			r *= x
		}
		x *= x
	}
	return
}
