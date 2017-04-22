# Gosl. ode. Ordinary differential equations

More information is available in **[the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxode.html).**

Package `ode` implements solution techniques to ordinary differential equations. Initially,
only algorithms based on the Runge-Kutta method are impelmented; i.e. the multisteps and others are
not available.

Also, some focus is given to methods that are able to handle stiff problems.

In summary, this package is _largelly_ based on Hairer and coworkers' book [1, 2].

The _golden_ method here is the Radau5 which is basically a rewrite from Hairer's Fortran code.
Nonetheless, capabilities to run in parallel have been added to it.

Roughly, the files are:
1. erk.go: Explicit Runge-Kutta
2. fweuler.go: Forward-Euler
3. bweuler.go: Backward-Euler
4. radau5.go: Radau5 !!!
5. ode.go: the _main_ file

Tests files are prefixed with `t_`

MPI tests (i.e. parallel) are suffixed with `_main` and can be run with the `xrunmpitests.bash`
script.

## Examples

### 1 Robertson's Equation

From Hairer-Wanner VII-p3 Eq.(1.4) [2].


```go
// problem definition
fcn := func(f []float64, dx, x float64, y []float64, args ...interface{}) error {
    f[0] = -0.04*y[0] + 1.0e4*y[1]*y[2]
    f[1] = 0.04*y[0] - 1.0e4*y[1]*y[2] - 3.0e7*y[1]*y[1]
    f[2] = 3.0e7 * y[1] * y[1]
    return nil
}
jac := func(dfdy *la.Triplet, dx, x float64, y []float64, args ...interface{}) error {
    if dfdy.Max() == 0 {
        dfdy.Init(3, 3, 9)
    }
    dfdy.Start()
    dfdy.Put(0, 0, -0.04)
    dfdy.Put(0, 1, 1.0e4*y[2])
    dfdy.Put(0, 2, 1.0e4*y[1])
    dfdy.Put(1, 0, 0.04)
    dfdy.Put(1, 1, -1.0e4*y[2]-6.0e7*y[1])
    dfdy.Put(1, 2, -1.0e4*y[1])
    dfdy.Put(2, 0, 0.0)
    dfdy.Put(2, 1, 6.0e7*y[1])
    dfdy.Put(2, 2, 0.0)
    return nil
}

// data
silent := false
fixstp := false
method := "Radau5"
xa, xb := 0.0, 0.3
ya := []float64{1.0, 0.0, 0.0}
ndim := len(ya)

// allocate ODE object
var o Solver
o.Init(method, ndim, fcn, jac, nil, SimpleOutput, silent)

// tolerances and initial step size
rtol := 1e-2
atol := rtol * 1e-6
o.SetTol(atol, rtol)
o.IniH = 1.0e-6

// solve problem
y := make([]float64, ndim)
copy(y, ya)
o.Solve(y, xa, xb, xb-xa, fixstp, nil)
```

Output:
```
number of F evaluations   =    87
number of J evaluations   =     8
total number of steps     =    17
number of accepted steps  =    15
number of rejected steps  =     1
number of decompositions  =    15
number of lin solutions   =    24
max number of iterations  =     2
```

<div id="container">
<p><img src="../examples/figs/rober.png" width="400"></p>
Solution of Robertson's equation
</div>


## References

[1] Hairer E, Norset SP, Wanner G. Solving O. Solving Ordinary Differential Equations I. Nonstiff
Problems, Springer. 1987

[2] Hairer E, Wanner G. Solving Ordinary Differential Equations II. Stiff and Differential-Algebraic
Problems, Second Revision Edition. Springer. 1996
