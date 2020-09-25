# Gosl. fun. Special functions, DFT, FFT, Bessel, elliptical integrals, orthogonal polynomials, interpolators

This package implements _special_ functions such as orthogonal polynomials and elliptical functions
of first, second and third kind.

Routines to interpolate and/or assist on spectral methods are also available; e.g. FourierInterp,
ChebyInterp.

## API

**go doc**

```
package fun // import "gosl/fun"

Package fun (functions) implements special functions such as elliptical,
orthogonal polynomials, Bessel, discrete Fourier transform, polynomial
interpolators, and more.

FUNCTIONS

func Atan2p(y, x float64) (αrad float64)
    Atan2p implements a positive version of atan2, in such a way that: 0 ≤ α ≤
    2π

func Atan2pDeg(y, x float64) (αdeg float64)
    Atan2pDeg implements a positive version of atan2, in such a way that: 0 ≤ α
    ≤ 360

func Beta(a, b float64) float64
    Beta computes the beta function by calling the Lgamma function

func Binomial(n, k int) float64
    Binomial comptues the binomial coefficient (n k)^T

func Boxcar(x, a, b float64) float64
    Boxcar implements the boxcar function

        Boxcar(x;a,b) = Heav(x-a) - Heav(x-b)

                        │ 0    if x < a or  x > b
        Boxcar(x;a,b) = ┤ 1/2  if x = a or  x = b
                        │ 1    if x > a and x < b

        Note: a ≤ x ≤ b; i.e. b ≥ a (not checked)

func CarlsonRc(x, y float64) float64
    CarlsonRc computes Carlson’s degenerate elliptic integral according to [1]
    Computes Rc(x,y) where x must be nonnegative and y must be nonzero. If y <
    0, the Cauchy principal value is returned.

func CarlsonRd(x, y, z float64) float64
    CarlsonRd computes Carlson’s elliptic integral of the second kind according
    to [1] Computes Rf(x,y,z) where x,y must be non-negative and at most one can
    be zero. z must be positive.

        References:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func CarlsonRf(x, y, z float64) float64
    CarlsonRf computes Carlson's elliptic integral of the first kind according
    to [1]. See also [2] Computes Rf(x,y,z) where x,y,z must be non-negative and
    at most one can be zero.

        References:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.
        [2] Carlson BC (1977) Elliptic Integrals of the First Kind, SIAM Journal on Mathematical
            Analysis, vol. 8, pp. 231-242.

func CarlsonRj(x, y, z, p float64) float64
    CarlsonRj computes Carlson’s elliptic integral of the third kind according
    to [1] Computes Rj(x,y,z,p) where x,y,z must be nonnegative, and at most one
    can be zero. p must be nonzero. If p < 0, the Cauchy principal value is
    returned.

func ChebyshevT(n int, x float64) float64
    ChebyshevT directly computes the Chebyshev polynomial of first kind Tn(x)
    using the trigonometric functions.

                │ (-1)ⁿ cosh[n⋅acosh(-x)]   if x < -1
        Tn(x) = ┤       cosh[n⋅acosh( x)]   if x > 1
                │       cos [n⋅acos ( x)]   if |x| ≤ 1

func ChebyshevTdiff1(n int, x float64) float64
    ChebyshevTdiff1 computes the first derivative of the Chebyshev function
    Tn(x)

        dTn
        ————
         dx

func ChebyshevTdiff2(n int, x float64) float64
    ChebyshevTdiff2 computes the second derivative of the Chebyshev function
    Tn(x)

        d²Tn
        —————
         dx²

func ChebyshevXgauss(N int) (X []float64)
    ChebyshevXgauss computes Chebyshev-Gauss roots considering symmetry

                    /  (2i+1)⋅π  \
        X[i] = -cos | —————————— |       i = 0 ... N
                    \   2N + 2   /

func ChebyshevXlob(N int) (X []float64)
    ChebyshevXlob computes Chebyshev-Gauss-Lobatto points using the sin function
    and considering symmetry

                            /  π⋅(N-2i)  \
                X[i] = -sin | —————————— |       i = 0 ... N
                            \    2⋅N     /
          or
                            /  i⋅π  \
                X[i] = -cos | ————— |       i = 0 ... N
                            \   N   /

        Reference:
        [1] Baltensperger R and Trummer MR (2003) Spectral differencing with a twist, SIAM J. Sci.
            Comput. 24(5):1465-1487

func Dft1d(data []complex128, inverse bool)
    Dft1d computes the discrete Fourier transform (DFT) in 1D. It replaces data
    by its discrete Fourier transform, if inverse==false or replaces data by its
    inverse discrete Fourier transform, if inverse==true

        Computes:
                           N-1         -i 2 π j k / N                 __
          forward:  X[k] =  Σ  x[j] ⋅ e                     with i = √-1
                           j=0

                           N-1         +i 2 π j k / N
          inverse:  Y[k] =  Σ  y[j] ⋅ e                     thus x[k] = Y[k] / N
                           j=0

        NOTE: (1) the inverse operation does not divide by N
              (2) ideally, N=len(data) is an integer power of 2.
              (3) using FFTW: http://fftw.org/fftw3_doc/What-FFTW-Really-Computes.html

func Elliptic1(φ, k float64) float64
    Elliptic1 computes Legendre elliptic integral of the first kind F(φ,k),
    evaluated using Carlson’s function Rf [1]. The argument ranges are 0 ≤ φ ≤
    π/2 and 0 ≤ k·sin(φ) ≤ 1

        Computes:
                           φ
                          ⌠          dt
              F(φ, k)  =  │  ___________________
                          │     _______________
                          ⌡   \╱ 1 - k² sin²(t)
                         0
        where:
                 0 ≤ φ ≤ π/2
                 0 ≤ k·sin(φ) ≤ 1

        References:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func Elliptic2(φ, k float64) float64
    Elliptic2 computes Legendre elliptic integral of the second kind E(φ,k),
    evaluated using Carlson's functions Rf and Rd [1]. The argument ranges are 0
    ≤ φ ≤ π/2 and 0 ≤ k⋅sin(φ) ≤ 1

        Computes:
                           φ
                          ⌠     _______________
              E(φ, k)  =  │   \╱ 1 - k² sin²(t)  dt
                          ⌡
                         0
        where:
                 0 ≤ φ ≤ π/2
                 0 ≤ k·sin(φ) ≤ 1

        References:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func Elliptic3(n, φ, k float64) float64
    Elliptic3 computes Legendre elliptic integral of the third kind Π(n,φ,k),
    evaluated using Carlson's functions Rf and Rj. NOTE that the sign convention
    on n corresponds to that of Abramowitz and Stegun [2] and not to [1]. The
    argument ranges are 0 ≤ φ ≤ π/2 and 0 ≤ k⋅sin(φ) ≤ 1

        Computes:
                              φ
                             ⌠                  dt
              Π(n, φ, k)  =  │  ___________________________________
                             │                     _______________
                             ⌡   (1 - n sin²(t)) \╱ 1 - k² sin²(t)
                            0
        where:
                 0 ≤ φ ≤ π/2
                 0 ≤ k·sin(φ) ≤ 1

        References:
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
            Scientific Computing. Third Edition. Cambridge University Press. 1235p.
        [2] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
            and Mathematical Tables. U.S. Department of Commerce, NIST

func ExpMix(x float64) complex128
    ExpMix uses Euler's formula to compute exp(-i⋅x) = cos(x) - i⋅sin(x)

func ExpPix(x float64) complex128
    ExpPix uses Euler's formula to compute exp(+i⋅x) = cos(x) + i⋅sin(x)

func Factorial100(n int) big.Float
    Factorial100 returns the factorial n! up to 100! using the math/big package

func Factorial22(n int) float64
    Factorial22 implements the factorial function; i.e. computes n! up to 22!
    According to [1], factorials up to 22! have exact double precision
    representations (52 bits of mantissa, not counting powers of two that are
    absorbed into the exponent)

        References
        [1] Press WH, Teukolsky SA, Vetterling WT, Fnannery BP (2007) Numerical Recipes: The Art of
             Scientific Computing. Third Edition. Cambridge University Press. 1235p.

func Hat(x, xc, y0, h, l float64) float64
    Hat implements the hat function

           --———--   o (xc,y0+h)
              |     / \
              h    /   \    m = h/l
              |   /m    \
        y0 ——————o       o—————————

                 |<  2l >|

func HatD1(x, xc, y0, h, l float64) float64
    HatD1 returns the first derivative of the hat function NOTE: the
    discontinuity is ignored ⇒ D1(xc-l)=D1(xc+l)=D1(xc)=0

func Heav(x float64) float64
    Heav computes the Heaviside step function (== derivative of Ramp(x))

                  │ 0    if x < 0
        Heav(x) = ┤ 1/2  if x = 0
                  │ 1    if x > 0

func ImagPowN(n int) complex128
    ImagPowN computes iⁿ = (√-1)ⁿ

        i¹ = i      i²  = -1      i³  = -i      i⁴  = 1
        i⁵ = i      i⁶  = -1      i⁷  = -i      i⁸  = 1
        i⁹ = i      i¹⁰ = -1      i¹¹ = -i      i¹² = 1

func ImagXpowN(x float64, n int) complex128
    ImagXpowN computes (x⋅i)ⁿ

        (x⋅i)¹ = x¹⋅i      (x⋅i)²  = -x²       (x⋅i)³  = -x³ ⋅i      (x⋅i)⁴  = x⁴
        (x⋅i)⁵ = x⁵⋅i      (x⋅i)⁶  = -x⁶       (x⋅i)⁷  = -x⁷ ⋅i      (x⋅i)⁸  = x⁸
        (x⋅i)⁹ = x⁹⋅i      (x⋅i)¹⁰ = -x¹⁰      (x⋅i)¹¹ = -x¹¹⋅i      (x⋅i)¹² = x¹²

func Logistic(z float64) float64
    Logistic implements the sigmoid/logistic function

func LogisticD1(z float64) float64
    LogisticD1 implements the first derivative of the sigmoid/logistic function

func ModBesselI0(x float64) (ans float64)
    ModBesselI0 returns the modified Bessel function I0(x) for any real x.

func ModBesselI1(x float64) (ans float64)
    ModBesselI1 returns the modified Bessel function I1(x) for any real x.

func ModBesselIn(n int, x float64) (ans float64)
    ModBesselIn returns the modified Bessel function In(x) for any real x and n
    ≥ 0

func ModBesselK0(x float64) float64
    ModBesselK0 returns the modified Bessel function K0(x) for positive real x.

        Special cases:
          K0(x=0) = +Inf
          K0(x<0) = NaN

func ModBesselK1(x float64) float64
    ModBesselK1 returns the modified Bessel function K1(x) for positive real x.

        Special cases:
          K0(x=0) = +Inf
          K0(x<0) = NaN

func ModBesselKn(n int, x float64) float64
    ModBesselKn returns the modified Bessel function Kn(x) for positive x and n
    ≥ 0

func NegOnePowN(n int) float64
    NegOnePowN computes (-1)ⁿ

func Pow2(x float64) float64
    Pow2 computes x²

func Pow3(x float64) float64
    Pow3 computes x³

func PowP(x float64, n uint32) (r float64)
    PowP computes real raised to positive integer xⁿ

func Ramp(x float64) float64
    Ramp function => MacAulay brackets

func Rbinomial(x, y float64) float64
    Rbinomial computes the binomial coefficient with real (non-negative)
    arguments by calling the Gamma function

func Rect(x float64) float64
    Rect implements the rectangular function

        Rect(x) = Boxcar(x;-0.5,0.5)

                  │ 0    if |x| > 1/2
        Rect(x) = ┤ 1/2  if |x| = 1/2
                  │ 1    if |x| < 1/2

func Sabs(x, eps float64) float64
    Sabs implements a smooth abs function: Sabs(x) = x*x / (sign(x)*x + eps)

func SabsD1(x, eps float64) float64
    SabsD1 returns the first derivative of Sabs

func SabsD2(x, eps float64) float64
    SabsD2 returns the first derivative of Sabs

func Sign(x float64) float64
    Sign implements the sign function

                  │ -1   if x < 0
        Sign(x) = ┤  0   if x = 0
                  │  1   if x > 0

func Sinc(x float64) float64
    Sinc computes the sine cardinal (sinc) function

        Sinc(x) = |     1      if x = 0
                  | sin(x)/x   otherwise

func Sramp(x, β float64) float64
    Sramp implements a smooth ramp function. Ramp

func SrampD1(x, β float64) float64
    SrampD1 returns the first derivative of Sramp

func SrampD2(x, β float64) float64
    SrampD2 returns the second derivative of Sramp

func SuqCos(angle, expon float64) float64
    SuqCos implements the superquadric auxiliary function that uses cos(x)

func SuqSin(angle, expon float64) float64
    SuqSin implements the superquadric auxiliary function that uses sin(x)

func UintBinomial(n, k uint64) uint64
    UintBinomial implements the Binomial coefficient using uint64. Panic happens
    on overflow Also, this function uses a loop so it may not be very efficient
    for large k The code below comes from
    https://en.wikipedia.org/wiki/Binomial_coefficient [cannot find a reference
    to cite]


TYPES

type Axis struct {

	// configuration data
	DisableHunt bool // do not use hunt code at all

	// Has unexported fields.
}
    Axis implements a type to hold an arbitrarily spaced discrete data

func NewAxis(data []float64, interpType InterpType) (o *Axis)
    NewAxis builds a new Axis type from a data slice for an InterpType

func (o *Axis) Get(i int) float64
    Get returns the value at data[i]

type BiLinear struct {
	// Has unexported fields.
}
    BiLinear implements a two dimensional interpolant

func NewBiLinear(f, xx, yy []float64) (o *BiLinear)
    NewBiLinear builds a two dimensional bi-linear interpolant Input:

        xx -- function sample points abscissas
        yy -- function sample points ordinates
        f  -- function values
        	f(i,j) is stored at f[len(xx)*j + i]

    Ref:

         f(x,y) = x^2 + 2y^2		 	  2 h  	i  	j
         xx = [0.00,0.50,1.00]    		|
         yy = [0.00,1.00,2.00] 	 	  1 d  	e  	f
         f  = [a:0.00,b:0.25,c:1.00,		|
        		  d:2.00,e:2,25,f:3.00,		a___b___c
        		  h:8.00,i:8.25,j:9.00]	   0   0.5 1.0

func (o *BiLinear) P(x, y float64) float64
    P is the interpolation polynomial

func (o *BiLinear) Reset(f, xx, yy []float64)
    Reset (Re)Set the axis and matrix for the interpolant

func (o *BiLinear) SetDisableHunt(disable bool)
    SetDisableHunt disables the hunt function for both axis

type ChebyInterp struct {

	// input
	N     int  // degree of polynomial
	Gauss bool // use roots (Gauss) or points (Lobatto)?

	// options
	Bary  bool // [default=true] use barycentric formulae in for ℓ_i and I{f} [default=true]
	Nst   bool // [default=true] use the "negative sum trick" to compute the diagonal components according to:
	Trig  bool // [default=false] use trigonometric identities (to reduce round-off errors)
	Flip  bool // [default=false] compute lower-diagonal part from upper diagonal part with D_{N-j,N-l} = -D_{j,l}
	StdD2 bool // [default=false] compute D2 using standard formula, otherwise use D1

	// constants
	EstimationN int // N to use when estimating CoefP [default=128]

	// derived
	X     []float64 // points. NOTE: mirrowed version of Chebyshev X; i.e. from +1 to -1
	Wb    []float64 // weights for Gaussian quadrature
	Gamma []float64 // denominator of coefficients equation ~ ‖p[i]‖²
	Lam   []float64 // λ_i barycentric weights (also w_i in some papers)

	// computed by auxiliary methods
	CoefI  []float64 // coefficients of interpolant
	CoefP  []float64 // coefficients of projection (estimated)
	CoefIs []float64 // coefficients of interpolation using Lagrange cardinal functions

	// computed
	C  *la.Matrix // physical to transform space conversion matrix
	Ci *la.Matrix // transform to physical space conversion matrix
	D1 *la.Matrix // (dℓj/dx)(xi)
	D2 *la.Matrix // (d²ℓj/dx²)(xi)
}
    ChebyInterp defines a structure for efficient computations with Chebyshev
    polynomials such as projecttion or interpolation

        Some equations are based on [1,2,3]

        References:
          [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
              Single Domains. Springer. 563p
          [2] Webb M, Trefethen LN, Gonnet P (2012) Stability of barycentric interpolation formulas for
              extrapolation, SIAM J. Sci. Comput. Vol. 34, No. 6, pp. A3009-A3015
          [3] Baltensperger R, Trummer M (2003) Spectral Differencing with a twist,
              SIAM J. Sci. Comput., Vol. 24, No. 5, pp. 1465-1487

func NewChebyInterp(N int, gaussChebyshev bool) (o *ChebyInterp)
    NewChebyInterp returns a new ChebyInterp structure

        gaussChebyshev == true:

                          /  (2⋅j+1)⋅π  \
                X_j = cos | ——————————— |       j = 0 ... N
                          \   2⋅N + 2   /

        gaussChebyshev == false: (Gauss-Lobatto-Chebyshev)

                          /  j⋅π  \
                X_j = cos | ————— |       j = 0 ... N
                          \   N   /

        NOTE: X here is the mirrowed version of Chebyshev X; i.e. from +1 to -1

func (o *ChebyInterp) CalcCoefI(f Ss)
    CalcCoefI computes the coefficients of the interpolant by (slow) direct
    formula

                   1    N
        CoefI_k = ——— ⋅ Σ  f(x_j) ⋅ T_k(x_j) ⋅ wb_j
                  γ_k  j=0

        Thus (for Gauss-Lobatto):

                    2       N    1                   /  k⋅j⋅π  \
        CoefI_k = —————— ⋅  Σ  —————— ⋅ f(x_j) ⋅ cos | ——————— |
                  N⋅cb_k   j=0  cb_j                 \    N    /

        where:
                cb_k = 2 if j=0,N   or   1 if j=1...N-1

        NOTE: the results will be stored in o.CoefI

func (o *ChebyInterp) CalcCoefIs(f Ss)
    CalcCoefIs computes the coefficients for interpolation with Lagrange
    cardinal functions ℓ_l(x)

func (o *ChebyInterp) CalcCoefP(f Ss)
    CalcCoefP computes the coefficients of the projection (slow) using
    o.EstimationN + 1 points

                    ∫ f(x)⋅T_k(x)⋅w(x) dx      (f, T_k)_w
        CoefP_k = ————————————————————————— = ————————————
                   ∫ T_k(x)⋅T_k(x)⋅w(x) dx      ‖ T_k ‖²

        NOTE: the results will be stored in o.CoefP

func (o *ChebyInterp) CalcConvMats()
    CalcConvMats computes conversion matrices C and Ci

                      N
        trans(u)_k =  Σ  C_{kj} ⋅ u(x_j)      e.g. trans(u)_k == coefI_k
                     j=0

               N
        u_j =  Σ   C⁻¹_{jk} ⋅ trans(u)_k
              k=0

                       2             /  j⋅k⋅π  \
        C_{kj} = ————————————— ⋅ cos | ——————— |
                  cb_k⋅cb_j⋅N        \    N    /

                       /  j⋅k⋅π  \
        C⁻¹_{jk} = cos | ——————— |
                       \    N    /

func (o *ChebyInterp) CalcD1()
    CalcD1 computes the differentiation matrix D1 of the function L_i

         d I{f}(x)     N            d ℓ_l(x)
        ——————————— =  Σ   f(x_l) ⋅ ————————
             dx       l=0              dx

         d I{f}(x)  |         N
        ——————————— |      =  Σ   D1_jl ⋅ f(x_l)
             dx     |x=x_j   l=0

        where:

                 dℓ_l  |
         D1_jl = ————— |
                  dx   |x=x_j

        Equations (2.4.31) and (2.4.33), page 89 of [1]

        If Nst==true (negative-sum-trick):

                                          N
                               D_{jj} = - Σ  D_{jl}
                                         l=0
                                         l≠j

        NOTE: (1) the signs are swapped (compared to [1]) because X are reversed here (from -1 to +1)
              (2) this method is only available for Gauss-Lobatto points

func (o *ChebyInterp) CalcD2()
    CalcD2 calculates the second derivative

                  d²ℓ_l  |
          D2_jl = —————— |
                   dx²   |x=x_j

        NOTE: this function will call CalcD1() because the D1 values required to compute D2,
              unless StdD2=true where the "standard" formula (Eq. 2.4.32) is used instead => less accurate

func (o *ChebyInterp) CalcErrorD1(dfdxAna Ss) (maxDiff float64)
    CalcErrorD1 computes the maximum error due to differentiation (@ X[i]) using
    the D1 matrix

        NOTE: CoefIs and D1 matrix must be computed previously

func (o *ChebyInterp) CalcErrorD2(d2fdx2Ana Ss) (maxDiff float64)
    CalcErrorD2 computes the maximum error due to differentiation (@ X[i]) using
    the D2 matrix

        NOTE: CoefIs and D2 matrix must be computed previously

func (o *ChebyInterp) EstimateMaxErr(f Ss, projection bool) (maxerr, xloc float64)
    EstimateMaxErr estimates the maximum error using 10000 stations along [-1,1]
    This function also returns the location (xloc) of the estimated max error

        maxerr = max(|f - I{f}|)  or  maxerr = max(|f - P{f}|)

        NOTE: CoefI or CoefP must be computed first

func (o *ChebyInterp) HierarchicalT(i int, x float64) float64
    HierarchicalT computes Tn(x) using hierarchical definition (but NOT
    recursive)

        NOTE: this function is not as efficient as ChebyshevT and should be used for testing only

func (o *ChebyInterp) I(x float64) (res float64)
    I computes the interpolation

                   N
        I{f}(x) =  Σ  CoefI_k ⋅ T_k(x)
                  k=0

        Thus:

                  N
        f(x_j) =  Σ   CoefI_k ⋅ T_k(x_j)
                 k=0

        NOTE: CoefI coefficients must be computed first

func (o *ChebyInterp) Il(x float64) (res float64)
    Il computes the interpolation using the Lagrange cardinal functions ℓ_i(x)

                   N                                         N
        I{f}(x) =  Σ   f(x_i) ⋅ ℓ_i(x)    or      I{f}(x) =  Σ  CoefIs_i ⋅ ℓ_i(x)
                  l=0                                       i=0

        NOTE: (1) CoefIs == f(x_i) coefficients must be computed (or set) first
              (2) ℓ is symbolised by ℓ in [1]

func (o *ChebyInterp) L(i int, x float64) float64
    L evaluates the Lagrange cardinal function ℓ_i(x) of degree N with
    Gauss-Lobatto points

                   N
        I{f}(x) =  Σ   f(x_i) ⋅ ℓ_i(x)
                  i=0

        Equation (2.4.30), page 88 of [1]

        NOTE: must not use with Gauss (roots) points

func (o *ChebyInterp) P(x float64) (res float64)
    P computes the (approximated) projection

                   ∞
        S{f}(x) =  Σ  CoefP_k ⋅ T_k(x)   (series representation)
                  k=0

        Thus:

                   N
        P{f}(x) =  Σ  CoefP_k ⋅ T_k(x)   (truncated series)
                  k=0

        NOTE: CoefP coefficients must be computed first

type DataInterp struct {

	// configuration data
	DisableHunt bool // do not use hunt code at all

	// output data
	Dy float64 // error estimate

	// Has unexported fields.
}
    DataInterp implements numeric interpolators to be used with discrete data

func NewDataInterp(Type string, p int, xx, yy []float64) (o *DataInterp)
    NewDataInterp creates new interpolator for data point sets xx and yy (with
    same lengths)

        Type -- type of interpolator
           "lin"  : linear
           "poly" : polynomial

        p  -- order of interpolator
        xx -- x-data
        yy -- y-data

func (o *DataInterp) P(x float64) float64
    P computes P(x); i.e. performs the interpolation

func (o *DataInterp) Reset(xx, yy []float64)
    Reset re-assigns xx and yy data sets

type FourierInterp struct {

	// main
	N int        // number of terms. must be power of 2; i.e. N = 2ⁿ
	X la.Vector  // point coordinates == 2⋅π.j/N
	K la.Vector  // k values computed from j such that j = 0...N-1 ⇒ k = -N/2...N/2-1
	A la.VectorC // coefficients for interpolation. from FFT
	S la.VectorC // smothing coefficients

	// computed (U may be set externally)
	U      la.Vector  // values of f(x) at grid points (nodes) X[j]
	Du     la.Vector  // p-order derivative of u
	Du1    la.Vector  // 1st derivative of f(x) at grid points (nodes) X[j]
	Du2    la.Vector  // 2nd derivative of f(x) at grid points (nodes) X[j]
	DuHat  la.VectorC // spectral coefficient corresponding to p-derivative
	Du1Hat la.VectorC // spectral coefficient corresponding to 1st derivative
	Du2Hat la.VectorC // spectral coefficient corresponding to 1st derivative

	// Has unexported fields.
}
    FourierInterp performs interpolation using truncated Fourier series

                    N/2 - 1
                     ————          +i k X[j]
          f(x[j]) =  \     A[k] ⋅ e                   with    X[j] = 2 π j / N
                     /
                     ————
                    k = -N/2                 Eq (2.1.27) of [1]    x ϵ [0, 2π]

          where:

                      N - 1
                  1   ————             -i k X[j]
          A[k] = ———  \     f(x[j]) ⋅ e              with    X[j] = 2 π j / N
                  N   /
                      ————
                     j = 0                                  Eq (2.1.25) of [1]

        NOTE: (1) f=u in [1] and A[k] is the tilde(u[k]) of [1]
              (2) FFTW says "so you should initialize your input data after creating the plan."
                  Therefore, the plan can be created and reused several times.
                  [http://www.fftw.org/fftw3_doc/Planner-Flags.html]
                  Also: "The plan can be reused as many times as needed. In typical high-performance
                  applications, many transforms of the same size are computed"
                  [http://www.fftw.org/fftw3_doc/Introduction.html]

        Create a new object with NewFourierInterp(...) AND deallocate memory with Free()

        Reference:
          [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
              Single Domains. Springer. 563p

func NewFourierInterp(N int, smoothing string) (o *FourierInterp)
    NewFourierInterp allocates a new FourierInterp object

        N -- number of terms. must be even; ideally power of 2, e.g. N = 2ⁿ

        smoothing -- type of smoothing: use SmoNoneKind for no smoothing
          "" or "none" : no smoothing
          "lanc"       : Lanczos (sinc)
          "rcos"       : Raised Cosine
          "ces"        : Cesaro

        NOTE: remember to call Free in the end to release memory allocatedy by FFTW; e.g.
              defer o.Free()

func (o *FourierInterp) CalcA()
    CalcA calculates the coefficients A of the interpolation using (fwd) FFT

                       N - 1
                   1   ————             -i k X[j]
           A[k] = ———  \     f(x[j]) ⋅ e              with    X[j] = 2 π j / N
                   N   /
                       ————
                      j = 0                                  Eq (2.1.25) of [1]

        NOTE: remember to set U (or call CalcU) first

func (o *FourierInterp) CalcAwithAliasRemoval(f Ss)
    CalcAwithAliasRemoval calculates the coefficients A by using the 3/2-rule to
    remove alias error via the padding method

        NOTE: with the 3/2-rule, the intepolatory property is not exact; i.e. I(xi)≈f(xi) only

func (o *FourierInterp) CalcD(p int)
    CalcD calculates the p-derivative of the interpolated function @ grid points
    using the FFT (with smoothing or not)

                        p      |
                       d(I{f}) |
               dfdx =  ——————— |             len(res) must be equal to N
                           p   |
                         dx    |x=x[j]

         INPUT:
            p -- derivative order

         OUTPUT:
            Du and DuHat will contain the results

        NOTE: remember to call CalcA first

func (o *FourierInterp) CalcD1()
    CalcD1 calculates the 1st derivative using function CalcD and internal
    arrays. See function CalcD for further details

        OUTPUT: the results will be stored in Du1 and Du1Hat

func (o *FourierInterp) CalcD2()
    CalcD2 calculates the 2nd derivative using function CalcD and internal
    arrays. See function CalcD for further details

        OUTPUT: the results will be stored in Du2 and Du2Hat

func (o *FourierInterp) CalcJ(k float64) int
    CalcJ computes j-index from k-index where j corresponds to the FFT index

        k ϵ [-N/2, N/2-1]
        j ϵ [0, N-1]

        Example with N = 8:

             k=0 ⇒ j=0      k=-4 ⇒ j=4
             k=1 ⇒ j=1      k=-3 ⇒ j=5      j = { N + k  if  k < 0
             k=2 ⇒ j=2      k=-2 ⇒ j=6          {     k  otherwise
             k=3 ⇒ j=3      k=-1 ⇒ j=7

func (o *FourierInterp) CalcK(j int) float64
    CalcK computes k-index from j-index where j corresponds to the FFT index

        FFT returns the A coefficients as:

           {A[0], A[1], ..., A[N/2-1], A[-N/2], A[-N/2+1], ... A[-1]}

        k ϵ [-N/2, N/2-1]
        j ϵ [0, N-1]

        Example with N = 8:

             j=0 ⇒ k=0      j=4 ⇒ k=-4
             j=1 ⇒ k=1      j=5 ⇒ k=-3
             j=2 ⇒ k=2      j=6 ⇒ k=-2
             j=3 ⇒ k=3      j=7 ⇒ k=-1

func (o *FourierInterp) CalcU(f Ss)
    CalcU calculates f(x) at grid points (to be used later with CalcA and/or
    CalcD)

func (o *FourierInterp) Free()
    Free releases resources allocated for FFTW

func (o *FourierInterp) I(x float64) float64
    I computes the interpolation (with smoothing or not)

                     N/2 - 1
                       ————          +i k x
           I {f}(x) =  \     A[k] ⋅ e                 x ϵ [0, 2π]
            N          /
                       ————
                      k = -N/2                 Eq (2.1.28) of [1]

        NOTE: remember to call CalcA first

func (o *FourierInterp) Idiff(p int, x float64) float64
    Idiff performs the differentiation of the interpolation; i.e. computes the
    p-derivative of the interpolation (with smoothing or not)

                         p       N/2 - 1
              p         d(I{f})    ————       p           +i k x
        res: DI{f}(x) = ——————— =  \     (i⋅k)  ⋅ A[k] ⋅ e
              N             p      /
                          dx       ————
                                  k = -N/2                   x ϵ [0, 2π]

        NOTE: remember to call CalcA first

type GeneralOrthoPoly struct {

	// input
	Kind string // type of orthogonal polynomial
	N    int    // (max) degree of polynomial. Lower order can be quickly obtained after this polynomial with max(N) is generated

	// Has unexported fields.
}
    GeneralOrthoPoly implements general orthogonal polynomials. It uses a
    general format and is NOT very efficient for large degrees. For efficiency,
    use the OrthoPoly structure instead.

        Reference:
        [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
            and Mathematical Tables. U.S. Department of Commerce, NIST

        NOTE: this structure should be not be used for high-performance computing;
              it's probably useful for verifications or learning purposes only

func NewGeneralOrthoPoly(kind string, N int, alpha, beta float64) (o *GeneralOrthoPoly)
    NewGeneralOrthoPoly creates a new orthogonal polynomial

        kind -- is the type of orthognal polynomial:
          "J" or "jac"    : Jacobi
          "L" or "leg"    : Legendre
          "H" or "her"    : Hermite
          "T" or "cheby1" : Chebyshev first kind
          "U" or "cheby2" : Chebyshev second kind

        N -- is the (max) degree of the polynomial.
             Lower order can later be quickly obtained after this
             polynomial with max(N) is created

        alpha -- Jacobi only: α coefficient

        beta -- Jacobi only: β coefficient

        NOTE: all coefficients for the 0...N polynomials will be generated

        NOTE: this structure should be not be used for high-performance computing;
              it's probably useful for verifications or learning purposes only

func (o *GeneralOrthoPoly) F(x float64) (res float64)
    F computes P(n,x) with n=N (max)

        Since GeneralOrthoPoly is a general form, the summations are directly implement; i.e. no
        advantages are taken w.r.t the structure of the polynomial. Thus, these functions are not
        highly efficient for large degrees N

func (o *GeneralOrthoPoly) P(n int, x float64) (res float64)
    P computes P(n,x) where n must be ≤ N

        Since GeneralOrthoPoly is a general form, the summations are directly implement; i.e. no
        advantages are taken w.r.t the structure of the polynomial. Thus, these functions are not
        highly efficient for large degrees N

type InterpCubic struct {
	A, B, C, D float64 // coefficients of polynomial
	TolDen     float64 // tolerance to avoid zero denominator
}
    InterpCubic computes a cubic polynomial to perform interpolation either
    using 4 points or 3 points and a known derivative

func NewInterpCubic() (o *InterpCubic)
    NewInterpCubic returns a new object

func (o *InterpCubic) Critical() (xmin, xmax, xifl float64, hasMin, hasMax, hasIfl bool)
    Critical returns the critical points

        xmin -- x @ min and y(xmin)
        xmax -- x @ max and y(xmax)
        xifl -- x @ inflection point and y(ifl)
        hasMin, hasMax, hasIfl -- flags telling what is available

func (o *InterpCubic) F(x float64) float64
    F computes y = f(x) curve

func (o *InterpCubic) Fit3pointsD(x0, y0, x1, y1, x2, y2, x3, d3 float64) (err error)
    Fit3pointsD fits polynomial to 3 points and known derivative

        (x0, y0) -- first point
        (x1, y1) -- second point
        (x2, y2) -- third point
        (x3, d3) -- derivative @ x3

func (o *InterpCubic) Fit4points(x0, y0, x1, y1, x2, y2, x3, y3 float64) (err error)
    Fit4points fits polynomial to 3 points

        (x0, y0) -- first point
        (x1, y1) -- second point
        (x2, y2) -- third point
        (x3, y3) -- fourth point

func (o *InterpCubic) G(x float64) float64
    G computes y' = df/x|(x) curve

type InterpQuad struct {
	A, B, C float64 // coefficients of polynomial
	TolDen  float64 // tolerance to avoid zero denominator
}
    InterpQuad computes a quadratic polynomial to perform interpolation either
    using 3 points or 2 points and a known derivative

func NewInterpQuad() (o *InterpQuad)
    NewInterpQuad returns a new object

func (o *InterpQuad) F(x float64) float64
    F computes y = f(x) curve

func (o *InterpQuad) Fit2pointsD(x0, y0, x1, y1, x2, d2 float64) (err error)
    Fit2pointsD fits polynomial to 2 points and known derivative

        (x0, y0) -- first point
        (x1, y1) -- second point
        (x2, d2) -- derivative @ x2

func (o *InterpQuad) Fit3points(x0, y0, x1, y1, x2, y2 float64) (err error)
    Fit3points fits polynomial to 3 points

        (x0, y0) -- first point
        (x1, y1) -- second point
        (x2, y2) -- third point

func (o *InterpQuad) G(x float64) float64
    G computes y' = df/x|(x) curve

func (o *InterpQuad) Optimum() (xopt, fopt float64)
    Optimum returns the minimum or maximum point; i.e. the point with zero
    derivative

        xopt -- x @ optimum
        fopt -- f(xopt) = y @ optimum

type InterpType int
    InterpType specifies the type of interpolant

const (

	// BiLinearType defines the bi-linear type
	BiLinearType InterpType = 2

	// BiCubicType defines the bi-cubic type
	BiCubicType InterpType = 3
)
type LagIntSet []*LagrangeInterp
    LagIntSet is groups interpolators together; e.g. 2D, 3D

func NewLagIntSet(ndim int, degrees []int, gridTypes []string) (lis LagIntSet)
    NewLagIntSet returns a set of LagrangeInterp

type LagrangeInterp struct {

	// general
	N int       // degree: N = len(X)-1
	X la.Vector // grid points: len(X) = P+1; generated in [-1, 1]
	U la.Vector // function evaluated @ nodes: f(x_i)

	// barycentric
	Bary   bool      // [default=true] use barycentric weights
	UseEta bool      // [default=true] use ηk when computing D1
	Eta    la.Vector // sum of log of differences: ηk = Σ ln(|xk-xl|) (k≠l)
	Lam    la.Vector // normalised barycentric weights λk = pow(-1, k+N) ⋅ ηk / (2ⁿ⁻¹/n)

	// computed
	D1 *la.Matrix // (dℓj/dx)(xi)
	D2 *la.Matrix // (d²ℓj/dx²)(xi)
}
    LagrangeInterp implements Lagrange interpolators associated with a grid X

        An interpolant I^X_N{f} (associated with a grid X; of degree N; with N+1 points)
        is expressed in the Lagrange form as follows:

                          N
              X          ————             X
             I {f}(x) =  \     f(x[i]) ⋅ ℓ (x)
              N          /                i
                         ————
                         i = 0

        where ℓ^X_i(x) is the i-th Lagrange cardinal polynomial associated with grid X and given by:

                      N
              N      ━━━━    x  -  X[j]
             ℓ (x) = ┃  ┃  —————————————           0 ≤ i ≤ N
              i      ┃  ┃   X[i] - X[j]
                    j = 0
                    j ≠ i

        or, barycentric form:

                          N   λ[i] ⋅ f[i]
                          Σ   ———————————
              X          i=0   x - x[i]
             I {f}(x) = ——————————————————
              N            N     λ[i]
                           Σ   ————————
                          i=0  x - x[i]

        with:

                          λ[i]
                        ————————
              N         x - x[i]
             ℓ (x) = ———————————————
              i        N     λ[k]
                       Σ   ————————
                      k=0  x - x[k]

        The barycentric weights λk are normalised and computed from ηk as follows:

           ηk = Σ ln(|xk-xl|) (k≠l)

                 a ⋅ b             k+N
           λk =  —————     a = (-1)        b = exp(m)    m = -ηk
                  lf0

           lf0 = 2ⁿ⁻¹/n

         or, if N > 700:

                 / a ⋅ b \   /  b  \   /  b  \
           λk =  | ————— | ⋅ | ——— | ⋅ | ——— |      b = exp(m/3)
                 \  lf0  /   \ lf1 /   \ lf2 /

           lf0⋅lf1⋅lf2 = 2ⁿ⁻¹/n

        References:
          [1] Canuto C, Hussaini MY, Quarteroni A, Zang TA (2006) Spectral Methods: Fundamentals in
              Single Domains. Springer. 563p
          [2] Berrut JP, Trefethen LN (2004) Barycentric Lagrange Interpolation,
              SIAM Review Vol. 46, No. 3, pp. 501-517
          [3] Costa B, Don WS (2000) On the computation of high order pseudospectral derivatives,
              Applied Numerical Mathematics, 33:151-159.

func NewLagrangeInterp(N int, gridType string) (o *LagrangeInterp)
    NewLagrangeInterp allocates a new LagrangeInterp

        N -- degree

        gridType -- type of grid:
           "uni" : uniform 1D grid kind
           "cg"  : Chebyshev-Gauss 1D grid kind
           "cgl" : Chebyshev-Gauss-Lobatto 1D grid kind

        NOTE: the grid will be generated in [-1, 1]

func (o *LagrangeInterp) CalcD1()
    CalcD1 computes the differentiation matrix D1 of the function L_i

         d I{f}(x)  |         N
        ——————————— |      =  Σ   D1_kj ⋅ f(x_j)
             dx     |x=x_k   j=0

        see [2]

func (o *LagrangeInterp) CalcD2()
    CalcD2 calculates the second derivative

                  d²ℓ_l  |
          D2_jl = —————— |
                   dx²   |x=x_j

        NOTE: this function will call CalcD1() because the D1 values required to compute D2

func (o *LagrangeInterp) CalcErrorD1(dfdxAna Ss) (maxDiff float64)
    CalcErrorD1 computes the maximum error due to differentiation (@ X[i]) using
    the D1 matrix

        NOTE: U and D1 matrix must be computed previously

func (o *LagrangeInterp) CalcErrorD2(d2fdx2Ana Ss) (maxDiff float64)
    CalcErrorD2 computes the maximum error due to differentiation (@ X[i]) using
    the D2 matrix

        NOTE: U and D2 matrix must be computed previously

func (o *LagrangeInterp) CalcU(f Ss)
    CalcU computes f(x_i); i.e. function f(x) @ all nodes

func (o *LagrangeInterp) EstimateLebesgue() (ΛN float64)
    EstimateLebesgue estimates the Lebesgue constant by using 10000 stations
    along [-1,1]

func (o *LagrangeInterp) EstimateMaxErr(nStations int, f Ss) (maxerr, xloc float64)
    EstimateMaxErr estimates the maximum error using 10000 stations along [-1,1]
    This function also returns the location (xloc) of the estimated max error

        Computes:
                  maxerr = max(|f(x) - I{f}(x)|)

        e.g. nStations := 10000 (≥2) will generate several points along [-1,1]

func (o *LagrangeInterp) I(x float64) (res float64)
    I computes the interpolation I^X_N{f}(x) @ x

                          N
              X          ————          X
             I {f}(x) =  \     U[i] ⋅ ℓ (x)       with   U[i] = f(x[i])
              N          /             i
                         ————
                         i = 0

        or (barycentric):

                         N   λ[i] ⋅ f[i]
                         Σ   ———————————
             X          i=0   x - x[i]
            I {f}(x) = ——————————————————
             N            N     λ[i]
                          Σ   ————————
                         i=0  x - x[i]

        NOTE: U[i] = f(x[i]) must be calculated with o.CalcU or set first

func (o *LagrangeInterp) L(i int, x float64) (lix float64)
    L computes the i-th Lagrange cardinal polynomial ℓ^X_i(x) associated with
    grid X

                      N
              X      ━━━━    x  -  X[j]
             ℓ (x) = ┃  ┃  —————————————           0 ≤ i ≤ N
              i      ┃  ┃   X[i] - X[j]
                    j = 0
                    j ≠ i

        or (barycentric):

                         λ[i]
                       ————————
             X         x - x[i]
            ℓ (x) = ———————————————
             i        N     λ[k]
                      Σ   ————————
                     k=0  x - x[k]

        Input:
           i -- index of X[i] point
           x -- where to evaluate the polynomial
        Output:
           lix -- ℓ^X_i(x)

func (o *LagrangeInterp) Om(x float64) (ω float64)
    Om computes the generating (nodal) polynomial associated with grid X. The
    nodal polynomial is the unique polynomial of degree N+1 and leading
    coefficient whose zeros are the N+1 nodes of X.

                 N
         X      ━━━━
        ω (x) = ┃  ┃ (x - X[i])
        N+1     ┃  ┃
               i = 0

type Mm func(f, m *la.Matrix)
    Mm defines a matrix function f(m) of a matrix argument m (matrix matrix))

        Input:
          m -- input matrix
        Output:
          M -- output matrix

type Mss func(m *la.Matrix, a, b float64)
    Mss defines a matrix function f(a,b) of two scalar arguments (matrix scalar
    scalar)

        Input:
          a -- first input scalar
          b -- second input scalar
        Output:
          m -- output matrix

type Mv func(f *la.Matrix, v la.Vector)
    Mv defines a matrix function f(v) of a vector argument v (matrix vector))

        Input:
          v -- input vector
        Output:
          f -- output matrix

type Sinusoid struct {

	// input
	Period     float64 // T: period; e.g. [s]
	MeanValue  float64 // A0: mean value; e.g. [m]
	Amplitude  float64 // C1: amplitude; e.g. [m]
	PhaseShift float64 // θ: phase shift; e.g. [rad]

	// derived
	Frequency   float64 // f: frequency; e.g. [Hz] or [1 cycle/s]
	AngularFreq float64 // ω0 = 2⋅π⋅f: angular frequency; e.g. [rad⋅s⁻¹]
	TimeShift   float64 // ts = θ / ω0: time shift; e.g. [s]

	// derived: coefficients
	A []float64 // A0, A1, A2, ... (if series mode)
	B []float64 // B0, B1, B2, ... (if series mode)
}
    Sinusoid implements the sinusoid equation:

        y(t) = A0 + C1⋅cos(ω0⋅t + θ)             [essential-form]

        y(t) = A0 + A1⋅cos(ω0⋅t) + B1⋅sin(ω0⋅t)  [basis-form]

        A1 =  C1⋅cos(θ)
        B1 = -C1⋅sin(θ)
        θ  = arctan(-B1 / A1)   if A1<0, θ += π
        C1 = sqrt(A1² + B1²)

func NewSinusoidBasis(T, A0, A1, B1 float64) (o *Sinusoid)
    NewSinusoidBasis creates a new Sinusoid object with the "basis" parameters
    set

        T  -- period; e.g. [s]
        A0 -- mean value; e.g. [m]
        A1 -- coefficient of the cos term
        B1 -- coefficient of the sin term

func NewSinusoidEssential(T, A0, C1, θ float64) (o *Sinusoid)
    NewSinusoidEssential creates a new Sinusoid object with the "essential"
    parameters set

        T  -- period; e.g. [s]
        A0 -- mean value; e.g. [m]
        C1 -- amplitude; e.g. [m]
        θ  -- phase shift; e.g. [rad]

func (o *Sinusoid) ApproxSquareFourier(N int)
    ApproxSquareFourier approximates sinusoid using Fourier series with N terms

func (o *Sinusoid) TestPeriodicity(tmin, tmax float64, npts int) bool
    TestPeriodicity tests that f(t) = f(T + t)

func (o *Sinusoid) Ybasis(t float64) (res float64)
    Ybasis computes y(t) = A0 + A1⋅cos(ω0⋅t) + B1⋅sin(ω0⋅t) [basis-form]

func (o *Sinusoid) Yessen(t float64) float64
    Yessen computes y(t) = A0 + C1⋅cos(ω0⋅t + θ [essential-form]

type Ss func(s float64) float64
    Ss defines a scalar function f(s) of a scalar argument s (scalar scalar)

        Input:
          s -- input scalar
        Returns:
          scalar

type Sss func(s1, s2 float64) float64
    Sss defines a scalar function f(r,s) of two scalar arguments (scalar scalar
    scalar)

        Input:
          s1, s2 -- input scalar
        Returns:
          scalar

type Sv func(v la.Vector) float64
    Sv defines a scalar functioin f(v) of a vector argument v (scalar vector)

        Input:
          v -- input vector
        Returns:
          scalar

type Svs func(v la.Vector, s float64) float64
    Svs defines a scalar function f(v,s) of a vector and a scalar

        Input:
          s -- the scalar
          v -- the vector
        Returns:
          scalar

type Tt func(f, t *la.Triplet)
    Tt defines a triplet (matrix) function f(t) of a triplet (matrix) argument t
    (triplet triplet)

        Input:
          t -- input triplet
        Output:
          f -- output triplet

type Tv func(f *la.Triplet, v la.Vector)
    Tv defines a triplet (matrix) function f(v) of a vector argument v (triplet
    vector)

        Input:
          v -- input vector
        Output:
          f -- output triplet

type Vs func(f la.Vector, s float64)
    Vs defines a vector function f(s) of a scalar argument s (vector scalar)

        Input:
          s -- input scalar
        Output:
          f -- output vector

type Vss func(f la.Vector, a, b float64)
    Vss defines a vector function f(a,b) of two scalar arguments (vector scalar
    scalar)

        Input:
          a -- first input scalar
          b -- second input scalar
        Output:
          f -- output vector

type Vv func(f, v la.Vector)
    Vv defines a vector function f(v) of a vector argument v (vector vector)

        Input:
          v -- input vector
        Output:
          f -- output vector

type Vvss func(u, v la.Vector, a, b float64)
    Vvss defines two vector functions u(a,b) and v(a,b) of two scalar arguments
    (vector vector scalar scalar)

        Input:
          a -- first input scalar
          b -- second input scalar
        Output:
          u -- first output vector
          v -- second output vector

type Vvvss func(u, v, w la.Vector, a, b float64)
    Vvvss defines three vector functions u(a,b), v(a,b) and w(a,b) of two scalar
    arguments (vector vector vector scalar scalar)

        Input:
          a -- first input scalar
          b -- second input scalar
        Output:
          u -- first output vector
          v -- second output vector
          w -- second output vector

```
