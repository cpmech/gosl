# Gosl. opt. Solvers for optimisation problems

More information is available in **[the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxopt.html).**

This package provides routines to solve optimisation problems. Currently, linear programming
problems can be solved with the interior-point method. In the future, more solution techniques will
be implemented directly in Go.

## Interior-point method for linear problems

```
LinIpm solves:

        min cᵀx   s.t.   Aᵀx = b, x ≥ 0
         x

or the dual problem:

        max bᵀλ   s.t.   Aᵀλ + s = c, s ≥ 0
         λ
```

Linear problems can be solved with the `LinIpm` structure. First, the problem definitions are
initialised with the `Init` command and by giving the matrix of constraint coefficients (A), the
right-hand side vector (b) of the constraints, and the vector defining the minimisation problem (c).

The matrix `A` is given as compressed-column sparse for efficiency purposes.



### Example 1

Simple linear problem:

```
linear programming problem:

  min cᵀx   s.t.   Aᵀx = b, x ≥ 0
   x

specific problem:

       min      -4*x0 - 5*x1
  {x0,x1,x2,x3}

   s.t.  2*x0 +   x1 ≤ 3
           x0 + 2*x1 ≤ 3
         x0,x1 ≥ 0

standard form:

   2*x0 +   x1 + x2     = 3
     x0 + 2*x1     + x3 = 3
   x0,x1,x2,x3 ≥ 0

as matrix:
                 / x0 \
  [-4  -5  0  0] | x1 | = cᵀ x
                 | x2 |
                 \ x3 /

   _            _   / x0 \
  |  2  1  1  0  |  | x1 | = Aᵀ x
  |_ 1  2  0  1 _|  | x2 |
                    \ x3 /

```

Code:
```go
// coefficients vector
c := []float64{-4, -5, 0, 0}

// constraints as a sparse matrix
var T la.Triplet
T.Init(2, 4, 6) // 2 by 4 matrix, with 6 non-zeros
T.Put(0, 0, 2.0)
T.Put(0, 1, 1.0)
T.Put(0, 2, 1.0)
T.Put(1, 0, 1.0)
T.Put(1, 1, 2.0)
T.Put(1, 3, 1.0)
Am := T.ToMatrix(nil) // compressed-column matrix

// right-hand side
b := []float64{3, 3}

// print arrays
A := Am.ToDense()
la.PrintMat("\nA", A, "%6g", false)
la.PrintVec("\nb", b, "%6g", false)
la.PrintVec("\nc", c, "%6g", false)
io.Pf("\n")

// solve LP
var ipm opt.LinIpm
defer ipm.Free()
ipm.Init(Am, b, c, nil)
err := ipm.Solve(true)
if err != nil {
    io.Pf("%v", err)
    return
}

// print solution
io.Pf("\n")
io.Pf("x = %v\n", ipm.X)
io.Pf("λ = %v\n", ipm.L)
io.Pf("s = %v\n", ipm.S)

// check solution
x := ipm.X[:2]
bchk := make([]float64, 2)
la.MatVecMul(bchk, 1, A, x)
io.Pf("b(check) = %v\n", bchk)

// plotting
plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
f := func(x []float64) float64 { return c[0]*x[0] + c[1]*x[1] }
g := func(x []float64, i int) float64 { return A[i][0]*x[0] + A[i][1]*x[1] - b[i] }
np := 41
argsF := &plt.A{CmapIdx: 0}
argsG := &plt.A{Levels: []float64{0}, Colors: []string{"yellow"}, Lw: 2, Fsz: 10}
vmin, vmax := []float64{-2.0, -2.0}, []float64{2.0, 2.0}
opt.PlotTwoVarsContour(x, np, nil, true, vmin, vmax, argsF, argsG, f,
    func(x []float64) float64 { return g(x, 0) },
    func(x []float64) float64 { return g(x, 1) },
)
plt.Equal()
plt.HideAllBorders()
plt.Gll("$x$", "$y$", &plt.A{LegOut: true})
plt.Save("/tmp/gosl", "opt_ipm01")
```

Output:
```
A =
     2     1     1     0
     1     2     0     1

b =      3     3

c =     -4    -5     0     0

 it            f(x)           error
  0 -9.99000000e+00  1.71974522e-01
  1 -8.65656141e+00  3.63052829e-02
  2 -8.99639576e+00  3.78555516e-04
  3 -8.99996396e+00  3.78424585e-06
  4 -8.99999964e+00  3.78423235e-08
  5 -9.00000000e+00  3.78423337e-10

x = [0.9999999990004347 1.000000000078799 1.9203318816792844e-09 8.419670861842801e-10]
λ = [-1.0000000003319234 -1.9999999997280653]
s = [7.256799795925211e-10 1.218218347079067e-10 1.0000000006656913 2.000000000061833]
b(check) = [2.9999999980796686 2.9999999991580326]
```

Source code: <a href="../examples/opt_ipm01.go">../examples/opt_ipm01.go</a>

<div id="container">
<p><img src="../examples/figs/opt_ipm01.png" width="500"></p>
</div>



### Example 2

Another linear problem:

```
linear programming problem:

  min cᵀx   s.t.   Aᵀx = b, x ≥ 0
   x

specific problem:

  min   2*x0 +   x1
  s.t.   -x0 +   x1 ≤ 1
          x0 +   x1 ≥ 2   →  -x0 - x1 ≤ -2
          x0 - 2*x1 ≤ 4
        x1 ≥ 0

standard (step 1) add slack
  s.t.   -x0 +   x1 + x2           = 1
         -x0 -   x1      + x3      = -2
          x0 - 2*x1           + x4 = 4

standard (step 2)
   replace x0 := x0_ - x5
   because it's unbounded

   min  2*x0_ +   x1                - 2*x5
   s.t.  -x0_ +   x1 + x2           +   x5 = 1
         -x0_ -   x1      + x3      +   x5 = -2
          x0_ - 2*x1           + x4 -   x5 = 4
        x0_,x1,x2,x3,x4,x5 ≥ 0
```

Code
```go
// coefficients vector
c := []float64{2, 1, 0, 0, 0, -2}

// constraints as a sparse matrix
var T la.Triplet
T.Init(3, 6, 12) // 3 by 6 matrix, with 12 non-zeros
T.Put(0, 0, -1)
T.Put(0, 1, 1)
T.Put(0, 2, 1)
T.Put(0, 5, 1)
T.Put(1, 0, -1)
T.Put(1, 1, -1)
T.Put(1, 3, 1)
T.Put(1, 5, 1)
T.Put(2, 0, 1)
T.Put(2, 1, -2)
T.Put(2, 4, 1)
T.Put(2, 5, -1)
Am := T.ToMatrix(nil) // compressed-column matrix

// right-hand side
b := []float64{1, -2, 4}

// print arrays
A := Am.ToDense()
la.PrintMat("\nA", A, "%6g", false)
la.PrintVec("\nb", b, "%6g", false)
la.PrintVec("\nc", c, "%6g", false)
io.Pf("\n")

// solve LP
var ipm opt.LinIpm
defer ipm.Free()
ipm.Init(Am, b, c, nil)
err := ipm.Solve(true)
if err != nil {
    io.Pf("%v", err)
    return
}

// print solution
io.Pf("\n")
io.Pf("x = %v\n", ipm.X)
io.Pf("λ = %v\n", ipm.L)
io.Pf("s = %v\n", ipm.S)

// check solution
x := ipm.X[:2]
x[0] -= ipm.X[5]
io.Pf("x = %v\n", x)
bchk := make([]float64, 3)
la.MatVecMul(bchk, 1, A, x)
io.Pf("b(check) = %v\n", bchk)

// plotting
plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
f := func(x []float64) float64 { return c[0]*x[0] + c[1]*x[1] }
g := func(x []float64, i int) float64 { return A[i][0]*x[0] + A[i][1]*x[1] - b[i] }
np := 41
argsF := &plt.A{CmapIdx: 1}
argsG := &plt.A{Levels: []float64{0}, Colors: []string{"yellow"}, Lw: 2, Fsz: 10}
vmin, vmax := []float64{-2.0, -2.0}, []float64{2.0, 2.0}
opt.PlotTwoVarsContour(x, np, nil, true, vmin, vmax, argsF, argsG, f,
    func(x []float64) float64 { return g(x, 0) },
    func(x []float64) float64 { return g(x, 1) },
)
plt.Equal()
plt.HideAllBorders()
plt.Gll("$x$", "$y$", &plt.A{LegOut: true})
plt.Save("/tmp/gosl", "opt_ipm02")
```

Output:
```
A =
    -1     1     1     0     0     1
    -1    -1     0     1     0     1
     1    -2     0     0     1    -1

b =      1    -2     4

c =      2     1     0     0     0    -2

 it            f(x)           error
  0  4.82195674e+00  4.72141263e-01
  1  3.66300276e+00  4.21080232e-01
  2  2.67385434e+00  3.69702809e-02
  3  2.50182089e+00  5.20560741e-04
  4  2.50001821e+00  5.20845343e-06
  5  2.50000018e+00  5.20848030e-08
  6  2.50000000e+00  5.20848253e-10

x = [3.4562270104040635 1.4999999986259591 2.971515056222099e-09 2.2343305942772616e-10 6.499999995654445 2.9562270088065894]
λ = [-0.49999999992944033 -1.4999999999550573 4.316183389793475e-12]
s = [2.489351257414523e-10 1.207644468431149e-10 0.5000000000671894 1.5000000000928062 1.3343281402633547e-10 2.6562805299653977e-11]
x = [0.5000000015974742 1.4999999986259591]
b(check) = [0.999999997028485 -2.0000000002234333 -2.499999995654444]
```

Source code: <a href="../examples/opt_ipm02.go">../examples/opt_ipm02.go</a>

<div id="container">
<p><img src="../examples/figs/opt_ipm02.png" width="500"></p>
</div>
