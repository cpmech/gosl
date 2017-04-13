# Gosl. num. Fundamental numerical methods

More information is available in [the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxnum.html).

This package implements basic numerical methods such as for root finding, numerical quadrature,
numerical differentiation, and solution of simple nonlinear problems.



## Root finding

The problem of finding the _roots_ of any function **y(x)** is solved when the *x* values that lead
to **y(x) = 0** are found.

At the moment, two methods are available in the `num` package to solve these types of problems. The
first one is the Newton's method and the second one is the Brent's method; see e.g. [1]. The key
difference between these two methods is that the first requires the derivatives of **y** with
respect to **x**, the Jacobian.

The Newton's method is implemented by the `NlSolver` structure whereas Brent's method is in `Brent`.
Nonetheless, on Newton's method can solve nonlinear system of equations.



### Examples

Find the root of
```
    y(x) = x³ - 0.165 x² + 3.993e-4
```
within [0, 0.11]. We have to make sure that the root is bounded otherwise Brent's method doesn't
work.


Using Brent's method:

```go
// y(x) function
yx := func(x float64) (res float64, err error) {
    res = math.Pow(x, 3.0) - 0.165*math.Pow(x, 2.0) + 3.993e-4
    return
}

// range: be sure to enclose root
xa, xb := 0.0, 0.11

// initialise solver
var o num.Brent
o.Init(yx)

// solve
xo, err := o.Solve(xa, xb, false)

// output
yo, _ := yx(xo)
io.Pf("\n")
io.Pf("x      = %v\n", xo)
io.Pf("f(x)   = %v\n", yo)
io.Pf("nfeval = %v\n", o.NFeval)
io.Pf("niter. = %v\n", o.It)
```

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


Using Newton's method:

```go
// Function: y(x) = fx[0] with x = xvec[0]
fcn := func(fx, xvec []float64) (err error) {
    x := xvec[0]
    fx[0] = math.Pow(x, 3.0) - 0.165*math.Pow(x, 2.0) + 3.993e-4
    return
}

// Jacobian: dfdx(x) function
Jfcn := func(dfdx [][]float64, x []float64) (err error) {
    dfdx[0][0] = 3.0*x[0]*x[0] - 2.0*0.165*x[0]
    return
}

// trial solution
xguess := 0.03

// initialise solver
neq := 1      // number of equations
useDn := true // use dense Jacobian
numJ := false // numerical Jacobian
var o num.NlSolver
o.Init(neq, fcn, nil, Jfcn, useDn, numJ, nil)

// solve
x := []float64{xguess}
err := o.Solve(x, false)
if err != nil {
    chk.Panic("NlSolver filed: %v\n", err)
}

// Output
fx := []float64{123}
fcn(fx, x)
xo, yo := x[0], fx[0]
io.Pf("\n")
io.Pf("x      = %v\n", xo)
io.Pf("f(x)   = %v\n", yo)
io.Pf("nfeval = %v\n", o.NFeval)
io.Pf("niter. = %v\n", o.It)
```

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



## Numerical quadrature

`num` package implements a number of methods to compute numerical integrals; i.e. quadrature.

There are two kinds of functions/structures: (1) elementary ones; and (2) more advanced ones with
refinement.

The elementary methods are:
1. `Trapz` trapezoidal method for given discrete points
2. `TrapzF` trapezoidal method with given function
3. `TrapzRange` trapezoidal method with given function and x-range
4. `Trapz2D` trapezoidal method in 2D
5. `Simpson` Simpson's method
6. `Simps2D` Simpson's method in 2D

The structures using refinement are:
1. `Quadrature` is the interface
2. `Trap` implements the trapezoidal rule with refinement
3. `Simp` implements the Simpson's rule with refinement

### Examples. Elementary methods

Source code: <a href="t_integ_test.go">t_integ_test.go</a>

Discrete:
```go
x := []float64{4, 6, 8}
y := []float64{1, 2, 3}
A := Trapz(x, y) // should be 8.0
```

Function and range:
```go
y := func(x float64) float64 {
    return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
}
n := 11
x := utl.LinSpace(0, 1, n)
A := TrapzF(x, y)
A_ := TrapzRange(0, 1, n, y)
Acor := 1.08306090851465 // correct value
chk.Scalar(tst, "A", 1e-15, A, Acor)
chk.Scalar(tst, "A_", 1e-15, A_, Acor)
```

Volume:
```go
// Γ(1/4, 1)
gamma_1div4_1 := 0.2462555291934987088744974330686081384629028737277219

x := utl.LinSpace(0, 1, 11)
y := utl.LinSpace(0, 1, 11)
m, n := len(x), len(y)
f := la.MatAlloc(m, n)
for i := 0; i < m; i++ {
    for j := 0; j < n; j++ {
        f[i][j] = 8.0 * math.Exp(-math.Pow(x[i], 2)-math.Pow(y[j], 4))
    }
}
dx, dy := x[1]-x[0], y[1]-y[0]
Vt := Trapz2D(dx, dy, f)
Vs := Simps2D(dx, dy, f)
Vc := math.Sqrt(math.Pi) * math.Erf(1) * (math.Gamma(1.0/4.0) - gamma_1div4_1)
chk.Scalar(tst, "Vt", 0.0114830435645548, Vt, Vc)
chk.Scalar(tst, "Vs", 1e-4, Vs, Vc)
```

### Examples. Methods with refinement

Source code: <a href="t_quadrature_test.go">t_quadrature_test.go</a>

```go
// problem
y := func(x float64) float64 {
    return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
}
var err error
Acor := 1.08268158558

// trapezoidal rule
var T Quadrature
T = new(Trap)
T.Init(y, 0, 1, 1e-11)
A, err := T.Integrate()
if err != nil {
    io.Pforan(err.Error())
}
io.Pforan("A  = %v\n", A)
chk.Scalar(tst, "A", 1e-11, A, Acor)

// Simpson's rule
var S Quadrature
S = new(Simp)
S.Init(y, 0, 1, 1e-11)
A, err = S.Integrate()
if err != nil {
    io.Pforan(err.Error())
}
io.Pforan("A  = %v\n", A)
chk.Scalar(tst, "A", 1e-11, A, Acor)
```


## Numerical differentiation

Source code: <a href="../examples/num_deriv01">num_deriv01.go</a>

Check first and second derivative of `y(x) = sin(x)`

```go
// define function and derivative function
y_fcn := func(x float64) float64 { return math.Sin(x) }
dydx_fcn := func(x float64) float64 { return math.Cos(x) }
d2ydx2_fcn := func(x float64) float64 { return -math.Sin(x) }

// run test for 11 points
X := utl.LinSpace(0, 2*math.Pi, 11)
io.Pf("          %8s %23s %23s %23s\n", "x", "analytical", "numerical", "error")
for _, x := range X {

    // analytical derivatives
    dydx_ana := dydx_fcn(x)
    d2ydx2_ana := d2ydx2_fcn(x)

    // numerical derivative: dydx
    dydx_num, _ := num.DerivCentral(func(t float64, args ...interface{}) float64 {
        return y_fcn(t)
    }, x, 1e-3)

    // numerical derivative d2ydx2
    d2ydx2_num, _ := num.DerivCentral(func(t float64, args ...interface{}) float64 {
        return dydx_fcn(t)
    }, x, 1e-3)

    // check
    chk.PrintAnaNum(io.Sf("dy/dx   @ %.6f", x), 1e-10, dydx_ana, dydx_num, true)
    chk.PrintAnaNum(io.Sf("d²y/dx² @ %.6f", x), 1e-10, d2ydx2_ana, d2ydx2_num, true)
}

// generate 101 points for plotting
X = utl.LinSpace(0, 2*math.Pi, 101)
Y := make([]float64, len(X))
for i, x := range X {
    Y[i] = y_fcn(x)
}

// plot
plt.SetForPng(0.75, 300, 150)
plt.Plot(X, Y, "'b.-', clip_on=0, markevery=10, label='y(x)=sin(x)'")
plt.Gll("x", "y", "")
plt.SaveD("/tmp/gosl", "num_deriv01.png")
```


## Nonlinear problems

### Examples

Source code: <a href="t_nlsolver_test.go">t_nlsolver_test.go</a>

Find `x0` and `x1` such that `f0` and `f1` are zero, with:
```
f0(x0,x1) = 2.0*x0 - x1 - exp(-x0)
f1(x0,x1) = -x0 + 2.0*x1 - exp(-x1)
```

### Using analytical (sparse) Jacobian matrix

```go
ffcn := func(fx, x []float64) error {
    fx[0] = 2.0*x[0] - x[1] - math.Exp(-x[0])
    fx[1] = -x[0] + 2.0*x[1] - math.Exp(-x[1])
    return nil
}

Jfcn := func(dfdx *la.Triplet, x []float64) error {
    dfdx.Start()
    dfdx.Put(0, 0, 2.0+math.Exp(-x[0]))
    dfdx.Put(0, 1, -1.0)
    dfdx.Put(1, 0, -1.0)
    dfdx.Put(1, 1, 2.0+math.Exp(-x[1]))
    return nil
}

x := []float64{5.0, 5.0}
atol, rtol, ftol := 1e-10, 1e-10, 10*EPS
fx := make([]float64, len(x))
neq := len(x)

prms := map[string]float64{
    "atol":    atol,
    "rtol":    rtol,
    "ftol":    ftol,
    "lSearch": 1.0,
}

var o num.NlSolver
o.Init(neq, ffcn, Jfcn, nil, false, false, prms)
defer o.Free()

err := o.Solve(x, false)
if err != nil {
    chk.Panic(err.Error())
}
```

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

Just change the call to Init
```go
o.Init(neq, ffcn, nil, nil, false, true, prms)
```

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


## References

[1] G.Forsythe, M.Malcolm, C.Moler, Computer methods for mathematical
    computations. M., Mir, 1980, p.180 of the Russian edition
