# Gosl. ode. Ordinary differential equations

[![GoDoc](https://pkg.go.dev/github.com/cpmech/gosl/ode?status.svg)](https://pkg.go.dev/github.com/cpmech/gosl/ode)

More information is available in **[the documentation of this package](https://pkg.go.dev/github.com/cpmech/gosl/ode).**

Package `ode` implements solution techniques to ordinary differential equations, such as the
Runge-Kutta method. Methods that can handle stiff problems are also available.

## Examples

### Robertson's Equation

From Hairer-Wanner VII-p3 Eq.(1.4) [2].

<div id="container">
<p><img src="../examples/figs/rober.png" width="400"></p>
Solution of Robertson's equation
</div>

## Output of Tests

### Convergence of explicit Runge-Kutta methods

Source code: <a href="t_erk_test.go">t_erk_test.go</a>

<div id="container">
<p><img src="../examples/figs/t_erk04.png" width="500"></p>
</div>

<div id="container">
<p><img src="../examples/figs/t_erk05.png" width="500"></p>
</div>

## References

[1] Hairer E, Norset SP, Wanner G. Solving O. Solving Ordinary Differential Equations I. Nonstiff
Problems, Springer. 1987

[2] Hairer E, Wanner G. Solving Ordinary Differential Equations II. Stiff and Differential-Algebraic
Problems, Second Revision Edition. Springer. 1996
