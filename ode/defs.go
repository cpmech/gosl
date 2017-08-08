// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import "github.com/cpmech/gosl/la"

// Func defines the main function d{y}/dx = {f}(x, {y})
//
//   Here, the "main" function receives the stepsize h as well, i.e.
//
//     d{y}/dx := {f}(h=dx, x, {y})
//
//   Input:
//     h -- current stepsize = dx
//     x -- current x
//     y -- current {y}
//   Output:
//     f -- {f}(h, x, {y})
//
type Func func(f la.Vector, h, x float64, y la.Vector) error

// JacF defines the Jacobian matrix of Func
//
//   Here, the Jacobian function receives the stepsize h as well, i.e.
//
//   d{f}/d{y} := [J](h=dx, x, {y})
//
//   Input:
//     h -- current stepsize = dx
//     x -- current x
//     y -- current {y}
//   Output:
//     dfdy -- Jacobian matrix d{f}/d{y} := [J](h=dx, x, {y})
//
type JacF func(dfdy *la.Triplet, h, x float64, y la.Vector) error

// OutF defines a "callback" function to be called during the output of results
//   Input:
//     first -- whether this is the first output or not
//     h     -- stepsize = dx
//     x     -- scalar variable
//     y     -- vector variable
//   Output:
//     error -- this function can return an error to force stopping the simulation
type OutF func(first bool, h, x float64, y la.Vector) error

// YanaF defines a function to be used when computing analytical solutions
type YanaF func(x float64) float64

// ProbPlotF defines a function to plot problems' results
//type ProbPlotF func(dirout, fnk, label string, xf float64, yana yanaType, sol *Solver)
