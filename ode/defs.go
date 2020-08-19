// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "gosl/la"

// Func defines the main function d{y}/dx = {f}(x, {y})
//
//   Here, the "main" function receives the stepsize h as well, i.e.
//
//     d{y}/dx := {f}(h=dx, x, {y})
//
//   INPUT:
//     h -- current stepsize = dx
//     x -- current x
//     y -- current {y}
//
//   OUTPUT:
//     f -- {f}(h, x, {y})
//
type Func func(f la.Vector, h, x float64, y la.Vector)

// JacF defines the Jacobian matrix of Func
//
//   Here, the Jacobian function receives the stepsize h as well, i.e.
//
//   d{f}/d{y} := [J](h=dx, x, {y})
//
//   INPUT:
//     h -- current stepsize = dx
//     x -- current x
//     y -- current {y}
//
//   OUTPUT:
//     dfdy -- Jacobian matrix d{f}/d{y} := [J](h=dx, x, {y})
//
type JacF func(dfdy *la.Triplet, h, x float64, y la.Vector)

// StepOutF defines a callback function to be called when a step is accepted
//
//   INPUT:
//     istep -- index of step (0 is the very first output whereas 1 is the first accepted step)
//     h     -- stepsize = dx
//     x     -- scalar variable
//     y     -- vector variable
//
//   OUTPUT:
//     stop -- stop simulation (nicely)
//
type StepOutF func(istep int, h, x float64, y la.Vector) (stop bool)

// DenseOutF defines a function to produce a dense output (i.e. many equally spaced points,
// regardless of the actual stepsize)
//
//   INPUT:
//     istep -- index of step (0 is the very first output whereas 1 is the first accepted step)
//     h     -- best (current) h
//     x     -- current (just updated) x
//     y     -- current (just updated) y
//     xout  -- selected x to produce an output
//     yout  -- y values computed @ xout
//
//   OUTPUT:
//     stop -- stop simulation (nicely)
//
type DenseOutF func(istep int, h, x float64, y la.Vector, xout float64, yout la.Vector) (stop bool)

// YanaF defines a function to be used when computing analytical solutions
type YanaF func(res []float64, x float64)
