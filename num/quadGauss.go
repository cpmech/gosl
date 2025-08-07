// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/utl"
)

// QuadGaussL10 approximates the integral of the function f(x) between a and b, by ten-point
// Gauss-Legendre integration. The function is evaluated exactly ten times at interior points
// in the range of integration. See page 180 of [1].
//
//	Reference:
//	[1] Press WH, Teukolsky SA, Vetterling WT, Flannery BP (2007) Numerical Recipes: The Art of
//	    Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func QuadGaussL10(a, b float64, f fun.Ss) (res float64) {

	// constants
	x := []float64{0.1488743389816312, 0.4333953941292472, 0.6794095682990244, 0.8650633666889845, 0.9739065285171717}
	w := []float64{0.2955242247147529, 0.2692667193099963, 0.2190863625159821, 0.1494513491505806, 0.0666713443086881}

	// auxiliary variables
	xm := 0.5 * (b + a)
	xr := 0.5 * (b - a)
	s := 0.0 // will be twice the average value of the function, since the ten weights (five numbers above each used twice) sum to 2.

	// execute sum
	var dx, fp, fm float64
	for j := 0; j < 5; j++ {
		dx = xr * x[j]
		fp = f(xm + dx)
		fm = f(xm - dx)
		s += w[j] * (fp + fm)
	}
	res = s * xr // scale the answer to the range of integration.
	return
}

// GaussLegendreXW computes positions (xi) and weights (wi) to perform Gauss-Legendre integrations
//
//	Input:
//	  x1 -- lower limit of integration
//	  x2 -- upper limit of integration
//	  n  -- number of points for quadrature formula
//	Reference:
//	[1] Press WH, Teukolsky SA, Vetterling WT, Flannery BP (2007) Numerical Recipes: The Art of
//	    Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func GaussLegendreXW(x1, x2 float64, n int) (x, w []float64) {
	x = make([]float64, n)
	w = make([]float64, n)
	EPS := MACHEPS // relative precision.
	var z1, z, xm, xl, pp, p3, p2, p1 float64
	m := (n + 1) / 2 // The roots are symmetric in the interval, so we only have to find half of them.
	xm = 0.5 * (x2 + x1)
	xl = 0.5 * (x2 - x1)
	for i := 0; i < m; i++ { // Loop over the desired roots.
		z = math.Cos(3.141592654 * (float64(i) + 0.75) / (float64(n) + 0.5))
		// starting with this approximation to the ith root, we enter the main loop of refinement by Newton's method.
		it, MAXIT := 0, 10
		for it < MAXIT {
			p1 = 1.0
			p2 = 0.0
			for j := 0; j < n; j++ { // Loop up the recurrence relation to get the Legendre polynomial evaluated at z.
				p3 = p2
				p2 = p1
				p1 = ((2.0*float64(j)+1.0)*z*p2 - float64(j)*p3) / (float64(j) + 1.0)
			}
			// p1 is now the desired Legendre polynomial. We next compute pp, its derivative, by a standard relation involving also p2, the polynomial of one lower order.
			pp = float64(n) * (z*p1 - p2) / (z*z - 1.0)
			z1 = z
			z = z1 - p1/pp // Newton's method.
			if math.Abs(z-z1) < EPS {
				break
			}
			it++
		}
		if it == MAXIT {
			chk.Panic("Newton's method did not converge after %d iterations", it)
		}
		x[i] = xm - xl*z // Scale the root to the desired interval, and put in its symmetric counterpart.
		x[n-1-i] = xm + xl*z
		w[i] = 2.0 * xl / ((1.0 - z*z) * pp * pp) //  Compute the weight and its symmetric counterpart.
		w[n-1-i] = w[i]
	}
	return
}

// GaussJacobiXW computes positions (xi) and weights (wi) to perform Gauss-Jacobi integrations.
// The largest abscissa is returned in x[0], the smallest in x[n-1].
// The interval of integration is x ϵ [-1, 1]
//
//	Input:
//	  alp -- coefficient of the Jacobi polynomial
//	  bet -- coefficient of the Jacobi polynomial
//	  n   -- number of points for quadrature formula
//	Reference:
//	[1] Press WH, Teukolsky SA, Vetterling WT, Flannery BP (2007) Numerical Recipes: The Art of
//	    Scientific Computing. Third Edition. Cambridge University Press. 1235p.
func GaussJacobiXW(alf, bet float64, n int) (x, w []float64) {
	x = make([]float64, n)
	w = make([]float64, n)
	EPS := MACHEPS // relative precision.
	N := float64(n)
	var alfbet, an, bn, r1, r2, r3 float64
	var a, b, c, p1, p2, p3, pp, temp, z, z1 float64
	var l1, l2, l3, l4 float64
	for i := 0; i < n; i++ { // Loop over the desired roots. Initial guess for the largest root.
		if i == 0 {
			an = alf / N
			bn = bet / N
			r1 = (1.0 + alf) * (2.78/(4.0+N*N) + 0.768*an/N)
			r2 = 1.0 + 1.48*an + 0.96*bn + 0.452*an*an + 0.83*an*bn
			z = 1.0 - r1/r2
		} else if i == 1 { // Initial guess for the second largest root.
			r1 = (4.1 + alf) / ((1.0 + alf) * (1.0 + 0.156*alf))
			r2 = 1.0 + 0.06*(N-8.0)*(1.0+0.12*alf)/N
			r3 = 1.0 + 0.012*bet*(1.0+0.25*math.Abs(alf))/N
			z -= (1.0 - z) * r1 * r2 * r3
		} else if i == 2 { // Initial guess for the third largest root.
			r1 = (1.67 + 0.28*alf) / (1.0 + 0.37*alf)
			r2 = 1.0 + 0.22*(N-8.0)/N
			r3 = 1.0 + 8.0*bet/((6.28+bet)*N*N)
			z -= (x[0] - z) * r1 * r2 * r3
		} else if i == n-2 { // Initial guess for the second smallest root.
			r1 = (1.0 + 0.235*bet) / (0.766 + 0.119*bet)
			r2 = 1.0 / (1.0 + 0.639*(N-4.0)/(1.0+0.71*(N-4.0)))
			r3 = 1.0 / (1.0 + 20.0*alf/((7.5+alf)*N*N))
			z += (z - x[n-4]) * r1 * r2 * r3
		} else if i == n-1 { // Initial guess for the smallest root.
			r1 = (1.0 + 0.37*bet) / (1.67 + 0.28*bet)
			r2 = 1.0 / (1.0 + 0.22*(N-8.0)/N)
			r3 = 1.0 / (1.0 + 8.0*alf/((6.28+alf)*N*N))
			z += (z - x[n-3]) * r1 * r2 * r3
		} else { // Initial guess for the other roots.
			z = 3.0*x[i-1] - 3.0*x[i-2] + x[i-3]
		}
		alfbet = alf + bet
		it, MAXIT := 0, 10
		for it = 0; it < MAXIT; it++ { // Refinement by Newton’s method.
			temp = 2.0 + alfbet // Start the recurrence with P0 and P1 to avoid a division by zero when α+β=0 or -1
			p1 = (alf - bet + temp*z) / 2.0
			p2 = 1.0
			for j := 2; j <= n; j++ { // Loop up the recurrence relation to get the Jacobi polynomial evaluated at z.
				J := float64(j)
				p3 = p2
				p2 = p1
				temp = 2*J + alfbet
				a = 2 * J * (J + alfbet) * (temp - 2.0)
				b = (temp - 1.0) * (alf*alf - bet*bet + temp*(temp-2.0)*z)
				c = 2.0 * (J - 1 + alf) * (J - 1 + bet) * temp
				p1 = (b*p2 - c*p3) / a
			}
			pp = (N*(alf-bet-temp*z)*p1 + 2.0*(N+alf)*(N+bet)*p2) / (temp * (1.0 - z*z))
			// p1 is now the desired Jacobi polynomial. We next compute pp, its derivative, by
			// a standard relation involving also p2, the polynomial of one lower order.
			z1 = z
			z = z1 - p1/pp // Newton's formula.
			if math.Abs(z-z1) <= EPS {
				break
			}
		}
		if it == MAXIT {
			chk.Panic("Newton's method did not converge after %d iterations", it)
		}
		x[i] = z // Store the root and the weight.
		l1, _ = math.Lgamma(alf + N)
		l2, _ = math.Lgamma(bet + N)
		l3, _ = math.Lgamma(N + 1.0)
		l4, _ = math.Lgamma(N + alfbet + 1.0)
		w[i] = math.Exp(l1+l2-l3-l4) * temp * math.Pow(2.0, alfbet) / (pp * p2)
	}
	// sort positions
	utl.Qsort2(x, w)
	return
}
