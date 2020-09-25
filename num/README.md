# Gosl. num. Fundamental numerical methods

This package implements essential numerical methods such as for root finding, numerical quadrature,
numerical differentiation, and solution of simple nonlinear problems.

While the supackage [num/qpck](https://github.com/cpmech/gosl/tree/master/num/qpck) provides
advanced quadrature schemes (by wrapping [Quadpack](http://www.netlib.org/quadpack/)), this package
implements few (simpler) methods to compute numerical integrals. Here, there are two kinds of
algorithms: (1) basic methods for discrete data; and (2) using refinment for integrating general
functions.

## Example: Using Brent's method:

Find the root of

```
    y(x) = x³ - 0.165 x² + 3.993e-4
```

within [0, 0.11]. We have to make sure that the root is bounded otherwise Brent's method doesn't
work.

Source: <a href="../examples/num_brent01.go">../examples/num_brent01.go</a>

Output of Brent's solution:

```
  it                      x                   f(x)                    err
                                                                  1.0e-14
   0  1.100000000000000e-01 -2.662000000000001e-04  5.500000000000000e-02
   1  6.600000000000000e-02 -3.194400000000011e-05  3.300000000000000e-02
   2  6.044444444444443e-02  1.730305075445823e-05  2.777777777777785e-03
   3  6.239640011030302e-02 -1.676981032316081e-07  9.759778329292944e-04
   4  6.237766369176578e-02 -7.323468182796403e-10  9.666096236606754e-04
   5  6.237758151338346e-02  3.262039076357137e-15  4.108919116063703e-08
   6  6.237758151374950e-02  0.000000000000000e+00  4.108900814037142e-08

x      = 0.0623775815137495
f(x)   = 0
nfeval = 8
niter. = 6
```

<div id="container">
<p><img src="../examples/figs/num_brent01.png" width="400"></p>
Simple root finding problem solved by Brent's method.
</div>

## Example: Using Newton's method:

Same problem as before.

Source: <a href="../examples/num_newton01.go">../examples/num_newton01.go</a>

Output of NlSolver:

```
  it                    Ldx                 fx_max
                  (1.0e-04)              (1.0e-09)
   0  0.000000000000000e+00  2.778000000000000e-04
   1  3.745954692556634e+06  5.421253067129628e-05
   2  6.176571448942142e+05  1.391803634400563e-06
   2  1.515117884960284e+04  5.314115983194589e-10
. . . converged with fx_max. nit=2, nFeval=4, nJeval=3

x      = 0.062377521883073835
f(x)   = 5.314115983194589e-10
nfeval = 4
niter. = 2
```

<div id="container">
<p><img src="../examples/figs/num_newton01.png" width="400"></p>
Simple root finding problem solved by Newton's method.
</div>

## Example: Quadrature with discrete data

Source code: <a href="t_quadDisc_test.go">t_quadDisc_test.go</a>

## Example: Quadrature with methods using refinement

Source code: <a href="t_quadElem_test.go">t_quadElem_test.go</a>

## Example: numerical differentiation

Check first and second derivative of `y(x) = sin(x)`

See source code: <a href="../examples/num_deriv01.go">num_deriv01.go</a>

## Examples: nonlinear problems

Source code: <a href="t_nlsolver_test.go">t_nlsolver_test.go</a>

Find `x0` and `x1` such that `f0` and `f1` are zero, with:

```
f0(x0,x1) = 2.0*x0 - x1 - exp(-x0)
f1(x0,x1) = -x0 + 2.0*x1 - exp(-x1)
```

### Using analytical (sparse) Jacobian matrix

Output:

```
  it                    Ldx                 fx_max
                  (1.0e-05)              (1.0e-15)
   0  0.000000000000000e+00  4.993262053000914e+00
   1  8.266404824090484e+09  9.204814001140181e-01
   2  7.824673760247719e+08  9.107574803964624e-02
   3  9.482829646746747e+07  9.541527544986161e-04
   4  1.014523823737919e+06  1.051153347697564e-07
   5  1.117908116077260e+02  1.221245327087672e-15
   5  1.298802024636321e-06  1.110223024625157e-16
. . . converged with fx_max. nit=5, nFeval=12, nJeval=6
x    = [0.5671432904097838 0.5671432904097838]  expected = [0.5671 0.5671]
f(x) = [-1.1102230246251565e-16 -1.1102230246251565e-16]
```

### Using numerical Jacobian matrix

Output:

```
  it                    Ldx                 fx_max
                  (1.0e-05)              (1.0e-15)
   0  0.000000000000000e+00  4.993262053000914e+00
   1  8.266404846831476e+09  9.204814268661379e-01
   2  7.824673975180904e+08  9.107574923794759e-02
   3  9.482829862677714e+07  9.541518971115659e-04
   4  1.014522914054942e+06  1.051134990159852e-07
   5  1.117888589104628e+02  1.554312234475219e-15
   5  1.653020764599351e-06  1.110223024625157e-16
. . . converged with fx_max. nit=5, nFeval=24, nJeval=6
xx    = [0.5671432904097838 0.5671432904097838]  expected = [0.5671 0.5671]
f(xx) = [-1.1102230246251565e-16 -1.1102230246251565e-16]
```

### Using analytical dense Jacobian matrix

Just replace `Jfcn` with

```
JfcnD := func(dfdx [][]float64, x []float64) error {
    dfdx[0][0] = 2.0+math.Exp(-x[0])
    dfdx[0][1] = -1.0
    dfdx[1][0] = -1.0
    dfdx[1][1] = 2.0+math.Exp(-x[1])
    return nil
}
```

## API

**go doc**

```
package num // import "gosl/num"

Package num implements fundamental numerical methods such as numerical
derivative and quadrature, root finding solvers (Brent's and Newton's
methods), among others.

VARIABLES

var (
	MACHEPS = math.Nextafter(1, 2) - 1.0 // smallest number satisfying 1 + EPS > 1
)
    constants


FUNCTIONS

func CompareJac(tst *testing.T, ffcn fun.Vv, Jfcn fun.Tv, x []float64, tol float64)
    CompareJac compares Jacobian matrix (e.g. for testing)

func CompareJacDense(tst *testing.T, ffcn fun.Vv, Jfcn fun.Mv, x []float64, tol float64)
    CompareJacDense compares Jacobian matrix (e.g. for testing) in dense format

func CompareJacMpi(tst *testing.T, comm *mpi.Communicator, ffcn fun.Vv, Jfcn fun.Tv, x la.Vector, tol float64, distr bool)
    CompareJacMpi compares Jacobian matrix (e.g. for testing)

func DerivBwd4(x, h float64, f fun.Ss) (res float64)
    DerivBwd4 approximates the derivative df/dx using backward differences with
    4 points.

func DerivCen5(x, h float64, f fun.Ss) (res float64)
    DerivCen5 approximates the derivative df/dx using central differences with 5
    points.

func DerivFwd4(x, h float64, f fun.Ss) (res float64)
    DerivFwd4 approximates the derivative df/dx using forward differences with 4
    points.

func EqCubicSolveReal(a, b, c float64) (x1, x2, x3 float64, nx int)
    EqCubicSolveReal solves a cubic equation, ignoring the complex answers.

        The equation is specified by:
         x³ + a x² + b x + c = 0
        Notes:
         1) every cubic equation with real coefficients has at least one solution
            x among the real numbers
         2) from Numerical Recipes 2007, page 228
        Output:
         x[i] -- roots
         nx   -- number of real roots: 1, 2 or 3

func GaussJacobiXW(alf, bet float64, n int) (x, w []float64)
    GaussJacobiXW computes positions (xi) and weights (wi) to perform
    Gauss-Jacobi integrations. The largest abscissa is returned in x[0], the
    smallest in x[n-1]. The interval of integration is x ϵ [-1, 1]

        Input:
          alp -- coefficient of the Jacobi polynomial
          bet -- coefficient of the Jacobi polynomial
          n   -- number of points for quadrature formula
        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func GaussLegendreXW(x1, x2 float64, n int) (x, w []float64)
    GaussLegendreXW computes positions (xi) and weights (wi) to perform
    Gauss-Legendre integrations

        Input:
          x1 -- lower limit of integration
          x2 -- upper limit of integration
          n  -- number of points for quadrature formula
        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func Jacobian(J *la.Triplet, ffcn fun.Vv, x, fx, w []float64)
    Jacobian computes Jacobian (sparse) matrix

            Calculates (with N=n-1):
                df0dx0, df0dx1, df0dx2, ... df0dxN
                df1dx0, df1dx1, df1dx2, ... df1dxN
                     . . . . . . . . . . . . .
                dfNdx0, dfNdx1, dfNdx2, ... dfNdxN
        INPUT:
            ffcn : f(x) function
            x    : station where dfdx has to be calculated
            fx   : f @ x
            w    : workspace with size == n == len(x)
        RETURNS:
            J : dfdx @ x [must be pre-allocated]

func JacobianMpi(comm *mpi.Communicator, J *la.Triplet, ffcn fun.Vv, x, fx, w []float64, distr bool)
    JacobianMpi computes Jacobian (sparse) matrix

            Calculates (with N=n-1):
                df0dx0, df0dx1, df0dx2, ... df0dxN
                df1dx0, df1dx1, df1dx2, ... df1dxN
                     . . . . . . . . . . . . .
                dfNdx0, dfNdx1, dfNdx2, ... dfNdxN
        INPUT:
            ffcn : f(x) function
            x    : station where dfdx has to be calculated
            fx   : f @ x
            w    : workspace with size == n == len(x)
        RETURNS:
            J : dfdx @ x [must be pre-allocated]

func LinFit(x, y []float64) (a, b float64)
    LinFit computes linear fitting parameters. Errors on y-direction only

        y(x) = a + b⋅x

        See page 780 of [1]
        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func LinFitSigma(x, y []float64) (a, b, σa, σb, χ2 float64)
    LinFitSigma computes linear fitting parameters and variances (σa,σb) in the
    estimates of a and b Errors on y-direction only

        y(x) = a + b⋅x

        See page 780 of [1]
        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func LineSearch(x, fx []float64, ffcn fun.Vv, dx, x0, dφdx0 []float64, φ0 float64, maxIt int, dxIsMdx bool) (nFeval int)
    LineSearch finds a new point x along the direction dx, from x0, where the
    function has decreased sufficiently. The new function value is returned in
    fx

        INPUT:
            ffcn    -- f(x) callback
            dx      -- direction vector
            x0      -- initial x
            dφdx0   -- initial dφdx0 = fx * dfdx
            φ0      -- initial φ = 0.5 * dot(fx,fx)
            maxIt   -- max number of iterations
            dxIsMdx -- whether dx is actually -dx ==> IMPORTANT: dx will then be changed dx := -dx

        OUTPUT:
            x      -- updated x (along dx)
            fx     -- updated f(x)
            φ0     -- updated φ = 0.5 * dot(fx,fx)
            dx     -- changed to -dx if dx_is_mdx == true
            nFeval -- number of calls to f(x)

func QuadCs(a, b, ω float64, useSin bool, fid int, f func(x float64) float64) (res float64)
    QuadCs performs automatic integration (quadrature) using the cosine or sine
    weights QUADPACK routine AWOE (Automatic with weight, Oscillatory)

        INPUT:
          a      -- lower limit of integration
          b      -- upper limit of integration
          ω      -- omega
          useSin -- use sin(ω⋅x) instead of cos(ω⋅x)
          fid    -- index of goroutine (to avoid race problems)
          f      -- function defining the integrand

        OUTPUT:          b                                     b
                  res = ∫  f(x) ⋅ cos(ω⋅x) dx     or    res = ∫ f(x) ⋅ sin(ω⋅x) dx
                        a                                     a

func QuadDiscreteSimps2d(Lx, Ly float64, f [][]float64) (V float64)
    QuadDiscreteSimps2d approximates a double integral over the x-y plane with
    the elevation given by data points f[npts][npts]. Thus, the result is an
    estimate of the volume below the f[][] opints and the plane ortogonal to z @
    x=0. The very simple Simpson's method is used here.

        Lx -- total length of plane along x
        Ly -- total length of plane along y
        f  -- elevations f(x,y)

func QuadDiscreteSimpsonRF(a, b float64, n int, f fun.Ss) (res float64)
    QuadDiscreteSimpsonRF approximates the area below the discrete curve defined
    by [xa,xy] range and y function. Computations are carried out with the (very
    simple) Simpson method from xa to xb, with npts points

func QuadDiscreteTrapz2d(Lx, Ly float64, f [][]float64) (V float64)
    QuadDiscreteTrapz2d approximates a double integral over the x-y plane with
    the elevation given by data points f[npts][npts]. Thus, the result is an
    estimate of the volume below the f[][] opints and the plane ortogonal to z @
    x=0. The very simple trapezoidal method is used here.

        Lx -- total length of plane along x
        Ly -- total length of plane along y
        f  -- elevations f(x,y)

func QuadDiscreteTrapzRF(xa, xb float64, npts int, y fun.Ss) (A float64)
    QuadDiscreteTrapzRF approximates the area below the discrete curve defined
    by [xa,xy] range and y function. Computations are carried out with the (very
    simple) trapezoidal rule from xa to xb, with npts points

func QuadDiscreteTrapzXF(x []float64, y fun.Ss) (A float64)
    QuadDiscreteTrapzXF approximates the area below the discrete curve defined
    by x points and y function. Computations are carried out with the (very
    simple) trapezoidal rule.

func QuadDiscreteTrapzXY(x, y []float64) (A float64)
    QuadDiscreteTrapzXY approximates the area below the discrete curve defined
    by x and y points. Computations are carried out with the trapezoidal rule.

func QuadExpIx(a, b, m float64, fid int, f func(x float64) float64) (res complex128)
    QuadExpIx approximates the integral of f(x) ⋅ exp(i⋅m⋅x) with i = √-1

        INPUT:
          a      -- lower limit of integration
          b      -- upper limit of integration
          m      -- coefficient of x
          fid    -- index of goroutine (to avoid race problems)
          f      -- function defining the integrand

        OUTPUT:        b                           b                           b
                res = ∫  f(x) ⋅ exp(i⋅m⋅x) dx   = ∫  f(x) ⋅ cos(m⋅x) dx + i ⋅ ∫  f(x) ⋅ sin(m⋅x) dx
                      a                           a                           a

func QuadGaussL10(a, b float64, f fun.Ss) (res float64)
    QuadGaussL10 approximates the integral of the function f(x) between a and b,
    by ten-point Gauss-Legendre integration. The function is evaluated exactly
    ten times at interior points in the range of integration. See page 180 of
    [1].

        Reference:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func QuadGen(a, b float64, fid int, f func(x float64) float64) (res float64)
    QuadGen performs automatic integration (quadrature) using the
    general-purpose QUADPACK routine AGSE (Automatic, general-purpose,
    end-points singularities).

        INPUT:
          a      -- lower limit of integration
          b      -- upper limit of integration
          fid    -- index of goroutine (to avoid race problems)
          f      -- function defining the integrand

        OUTPUT:          b
                  res = ∫  f(x) dx
                        a

func SecondDerivCen3(x, h float64, f fun.Ss) float64
    SecondDerivCen3 approximates the second derivative d²f/dx² using central
    differences with 3 points

func SecondDerivCen5(x, h float64, f fun.Ss) float64
    SecondDerivCen5 approximates the second derivative d²f/dx² using central
    differences with 5 points

func SecondDerivMixedO2(x, y, h float64, f fun.Sss) float64
    SecondDerivMixedO2 approximates ∂²f/∂x∂y @ x={x,y} using O(h²) formula

func SecondDerivMixedO4v1(x, y, h float64, f fun.Sss) float64
    SecondDerivMixedO4v1 approximates ∂²f/∂x∂y @ x={x,y} using O(h⁴) formula
    from
    http://www.holoborodko.com/pavel/numerical-methods/numerical-derivative/central-differences

func SecondDerivMixedO4v2(x, y, h float64, f fun.Sss) float64
    SecondDerivMixedO4v2 approximates ∂²f/∂x∂y @ x={x,y} using O(h⁴) formula
    from
    http://www.holoborodko.com/pavel/numerical-methods/numerical-derivative/central-differences


TYPES

type Bracket struct {

	// configuration
	MaxIt   int  // max iterations
	Verbose bool // show messages

	// statistics
	NumFeval int // number of calls to Ffcn (function evaluations)
	NumIter  int // number of iterations from last call to Solve

	// Has unexported fields.
}
    Bracket implements routines to bracket roots or optima

        A root of a function is known to be bracketed by a pair of points, a and b,
        when the function has opposite sign at those two points.

        A minimum is known to be bracketed only when there is a triplet of points,
        a < b < c (or c < b < a), such that f(b) is less than both f(a) and f(c)

        REFERENCES:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes:
            The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func NewBracket(ffcn fun.Ss) (o *Bracket)
    NewBracket returns a new bracket-er object

func (o *Bracket) Min(a0, b0 float64) (a, b, c, fa, fb, fc float64)
    Min brackets minimum

        Given a function and given distinct initial points a0 and b0, search in the downhill direction
        (defined by the function as evaluated at the initial points) and return new points a, b, c
        that bracket a minimum of the function.

        Returns also the function values at the three points, fa, fb, and fc

type Brent struct {

	// configuration
	MaxIt   int     // max iterations
	Tol     float64 // tolerance
	Verbose bool    // show messages

	// statistics
	NumFeval int // number of calls to Ffcn (function evaluations)
	NumJeval int // number of calls to Jfcn (Jacobian/derivatives)
	NumIter  int // number of iterations from last call to Solve

	Jfcn fun.Ss // Jfcn(x) = dy/dx [optional / may be nil]

	// Has unexported fields.
}
    Brent implements Brent's method for finding the roots of an equation

func NewBrent(ffcn, Jfcn fun.Ss) (o *Brent)
    NewBrent returns a new Brent structure

        ffcn -- function f(x)
        Jfcn -- derivative df(x)/dx [optinal / may be nil]

func (o *Brent) Min(xa, xb float64) (xAtMin float64)
    Min finds the minimum of f(x) in [xa, xb]

        Based on ZEROIN C math library: http://www.netlib.org/c/
        By: Oleg Keselyov <oleg@ponder.csci.unt.edu, oleg@unt.edu> May 23, 1991

         G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
         computations. M., Mir, 1980, p.202 of the Russian edition

         The function makes use of the "gold section" procedure combined with
         the parabolic interpolation.
         At every step program operates three abscissae - x,v, and w.
         x - the last and the best approximation to the minimum location,
             i.e. f(x) <= f(a) or/and f(x) <= f(b)
             (if the function f has a local minimum in (a,b), then the both
             conditions are fulfilled after one or two steps).
         v,w are previous approximations to the minimum location. They may
         coincide with a, b, or x (although the algorithm tries to make all
         u, v, and w distinct). Points x, v, and w are used to construct
         interpolating parabola whose minimum will be treated as a new
         approximation to the minimum location if the former falls within
         [a,b] and reduces the range enveloping minimum more efficient than
         the gold section procedure.
         When f(x) has a second derivative positive at the minimum location
         (not coinciding with a or b) the procedure converges superlinearly
         at a rate order about 1.324

         The function always obtains a local minimum which coincides with
         the global one only if a function under investigation being
         unimodular. If a function being examined possesses no local minimum
         within the given range, Fminbr returns 'a' (if f(a) < f(b)), otherwise
         it returns the right range boundary value b.

func (o *Brent) MinUseD(xa, xb float64) (xAtMin float64)
    MinUseD finds minimum and uses information about derivatives

        Given a function and deriva funcd that computes a function and also its derivative function df, and
        given a bracketing triplet of abscissas ax, bx, cx [such that bx is between ax and cx, and
        f(bx) is less than both f(ax) and f(cx)], this routine isolates the minimum to a fractional
        precision of about tol using a modification of Brent’s method that uses derivatives. The
        abscissa of the minimum is returned as xAtMin, and the minimum function value is returned
        as min, the returned function value.

        REFERENCES:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes:
            The Art of Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func (o *Brent) Root(xa, xb float64) (res float64)
    Root solves y(x) = 0 for x in [xa, xb] with f(xa) * f(xb) < 0

        Based on ZEROIN C math library: http://www.netlib.org/c/
        By: Oleg Keselyov <oleg@ponder.csci.unt.edu, oleg@unt.edu> May 23, 1991

         G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
         computations. M., Mir, 1980, p.180 of the Russian edition

         The function makes use of the bissection procedure combined with
         the linear or quadric inverse interpolation.
         At every step program operates on three abscissae - a, b, and c.
         b - the last and the best approximation to the root
         a - the last but one approximation
         c - the last but one or even earlier approximation than a that
             1) |f(b)| <= |f(c)|
             2) f(b) and f(c) have opposite signs, i.e. b and c confine
                the root
         At every step Zeroin selects one of the two new approximations, the
         former being obtained by the bissection procedure and the latter
         resulting in the interpolation (if a,b, and c are all different
         the quadric interpolation is utilized, otherwise the linear one).
         If the latter (i.e. obtained by the interpolation) point is
         reasonable (i.e. lies within the current interval [b,c] not being
         too close to the boundaries) it is accepted. The bissection result
         is used in the other case. Therefore, the range of uncertainty is
         ensured to be reduced at least by the factor 1.6

type ElementarySimpson struct {
	// Has unexported fields.
}
    ElementarySimpson structure implements the Simpson's method for quadrature
    with refinement.

func (o *ElementarySimpson) Init(f fun.Ss, a, b, eps float64)
    Init initialises Simp structure

func (o *ElementarySimpson) Integrate() (res float64)
    Integrate performs the numerical integration

func (o *ElementarySimpson) Next() (res float64)
    Next returns the nth stage of refinement of the extended trapezoidal rule.
    On the first call (n=1), R b the routine returns the crudest estimate of a f
    .x/dx. Subsequent calls set n=2,3,... and improve the accuracy by adding 2
    n-2 additional interior points.

type ElementaryTrapz struct {
	// Has unexported fields.
}
    ElementaryTrapz structure is used for the trapezoidal integration rule with
    refinement.

func (o *ElementaryTrapz) Init(f fun.Ss, a, b, eps float64)
    Init initialises Trap structure

func (o *ElementaryTrapz) Integrate() (res float64)
    Integrate performs the numerical integration

func (o *ElementaryTrapz) Next() (res float64)
    Next returns the nth stage of refinement of the extended trapezoidal rule.
    On the first call (n=1), R b the routine returns the crudest estimate of a f
    .x/dx. Subsequent calls set n=2,3,... and improve the accuracy by adding 2
    n-2 additional interior points.

type LineSolver struct {

	// configuration
	UseDeriv bool // use Jacobian function [default = true if Jfcn is provided]

	Jfcn fun.Vv // vector function of vector: {J} = df/d{x} @ {x} [optional / may be nil]

	// Stat
	NumFeval int // number of function evalutions
	NumJeval int // number of Jacobian evaluations

	// Has unexported fields.
}
    LineSolver finds the scalar λ that zeroes or minimizes f(x+λ⋅n)

func NewLineSolver(size int, ffcn fun.Sv, Jfcn fun.Vv) (o *LineSolver)
    NewLineSolver returns a new LineSolver object

        size -- length(x)
        ffcn -- scalar function of vector: y = f({x})
        Jfcn -- vector function of vector: {J} = df/d{x} @ {x} [optional / may be nil]

func (o *LineSolver) G(λ float64) float64
    G implements g(λ) := f({y}(λ)) where {y}(λ) := {x} + λ⋅{n}

func (o *LineSolver) H(λ float64) float64
    H implements h(λ) = dg/dλ = df/d{y} ⋅ d{y}/dλ where {y} == {x} + λ⋅{n}

func (o *LineSolver) Min(x, n la.Vector) (λ float64)
    Min finds the scalar λ that minimizes f(x+λ⋅n)

func (o *LineSolver) MinUpdateX(x, n la.Vector) (λ, fmin float64)
    MinUpdateX finds the scalar λ that minimizes f(x+λ⋅n), updates x and returns
    fmin = f({x})

        Input:
          x -- initial point
          n -- direction
        Output:
          λ -- scale parameter
          x -- x @ minimum
          fmin -- f({x})

func (o *LineSolver) Root(x, n la.Vector) (λ float64)
    Root finds the scalar λ that zeroes f(x+λ⋅n)

func (o *LineSolver) Set(x, n la.Vector)
    Set sets x and n vectors as required by G(λ) and H(λ) functions

type NlSolver struct {

	// stats
	Niter  int // number of iterations from the last call to Solve
	Nfeval int // number of calls to Ffcn (function evaluations)
	Njeval int // number of calls to Jfcn (Jacobian evaluations)
	// Has unexported fields.
}
    NlSolver implements a solver to nonlinear systems of equations

        References:
         [1] G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
             computations. M., Mir, 1980, p.180 of the Russian edition

func NewNlSolver(neq int, F fun.Vv) (o *NlSolver)
    NewNlSolver creates a new NlSolver F is the f(x) function f:vector, x:vector
    Will use numerical Jacobian (with sparse solver) by default

func (o *NlSolver) CheckJ(x []float64, tol float64, verbose bool) (cnd float64)
    CheckJ check Jacobian matrix

        Ouptut: cnd -- condition number (with Frobenius norm)

func (o *NlSolver) Free()
    Free frees memory

func (o *NlSolver) SetJacobianFunction(Jsparse fun.Tv, Jdense fun.Mv)
    SetJacobianFunction sets function to compute the Jacobian (dense or sparse)
    One of sparse [recommended] or dense must be given. If both sparse and dense
    functions are given, the sparse will be used. With Jdense, matrix inversion
    is used (not very efficient. use for small systems)

func (o *NlSolver) Solve(x []float64)
    Solve solves non-linear problem f(x) == 0 x -- trial x "near" the solution;
    otherwise it may not converge

type NlSolverConfig struct {
	// input
	Verbose          bool // show messages
	ConstantJacobian bool // constant Jacobian (Modified Newton's method)
	LineSearch       bool // use line search
	LineSearchMaxIt  int  // line search maximum iterations
	MaxIterations    int  // Newton's method maximum iterations
	EnforceConvRate  bool // check and enforce convergence rate

	// function to be called during each output
	OutCallback func(x []float64) // output callback function

	// configurations for linear solver
	LinSolConfig *la.SparseConfig // configurations for sparse linear solver

	// Has unexported fields.
}
    NlSolverConfig holds the configuration input for NlSolver

func NewNlSolverConfig() (o *NlSolverConfig)
    NewNlSolverConfig creates a new NlSolverConfig Default values:

        CteJac      = false
        LinSearch   = false
        LinSchMaxIt = 20
        MaxIt       = 20
        ChkConv     = false
        Atol        = 1e-8
        Rtol        = 1e-8
        Ftol        = 1e-9

func (o *NlSolverConfig) SetTolerances(atol, rtol, ftol float64)
    SetTolerances sets all tolerances

type QuadElementary interface {
	Init(f fun.Ss, a, b, eps float64) // The constructor takes as inputs f, the function or functor to be integrated between limits a and b, also input.
	Integrate() float64               // Returns the integral for the specified input data
}
    QuadElementary defines the interface for elementary quadrature algorithms
    with refinement.

```
