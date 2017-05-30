# Gosl. num. Fundamental numerical methods

More information is available in **[the documentation of this package](https://godoc.org/github.com/cpmech/gosl/num).**

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


### Using Brent's method:

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


### Using Newton's method:

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

There are two kinds of algorithms: (1) basic methods for discrete data; and (2) using refinment for
integrating general functions.

### Examples. Basic methods for discrete data

Source code: <a href="t_quadDisc_test.go">t_quadDisc_test.go</a>

### Examples. Methods with refinement

Source code: <a href="t_quadElem_test.go">t_quadElem_test.go</a>



## Numerical differentiation

Check first and second derivative of `y(x) = sin(x)`

See source code: <a href="../examples/num_deriv01.go">num_deriv01.go</a>



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
