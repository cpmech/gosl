// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

func problem(index int) (name string, xTrial, xReference la.Vector, funcF fun.Vv, funcJsparse fun.Tv, funcJdense fun.Mv) {
	switch index {
	case 0:
		name = "simple 2 equations"
		xTrial = []float64{0.5, 0.5}
		xReference = []float64{1.0, 0.0}
		funcF = func(fx, x la.Vector) {
			fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
			fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
			return
		}
		funcJsparse = func(dfdx *la.Triplet, x la.Vector) {
			dfdx.Start()
			dfdx.Put(0, 0, 3.0*x[0]*x[0])
			dfdx.Put(0, 1, 1.0)
			dfdx.Put(1, 0, -1.0)
			dfdx.Put(1, 1, 3.0*x[1]*x[1])
			return
		}
		funcJdense = func(dfdx *la.Matrix, x la.Vector) {
			dfdx.Set(0, 0, 3.0*x[0]*x[0])
			dfdx.Set(0, 1, 1.0)
			dfdx.Set(1, 0, -1.0)
			dfdx.Set(1, 1, 3.0*x[1]*x[1])
		}
		return

	case 1:
		name = "with exponential function"
		xTrial = []float64{5.0, 5.0}
		xReference = []float64{0.5671, 0.5671}
		funcF = func(fx, x la.Vector) {
			fx[0] = 2.0*x[0] - x[1] - math.Exp(-x[0])
			fx[1] = -x[0] + 2.0*x[1] - math.Exp(-x[1])
			return
		}
		funcJsparse = func(dfdx *la.Triplet, x la.Vector) {
			dfdx.Start()
			dfdx.Put(0, 0, 2.0+math.Exp(-x[0]))
			dfdx.Put(0, 1, -1.0)
			dfdx.Put(1, 0, -1.0)
			dfdx.Put(1, 1, 2.0+math.Exp(-x[1]))
			return
		}
		funcJdense = func(dfdx *la.Matrix, x la.Vector) {
			dfdx.Set(0, 0, 2.0+math.Exp(-x[0]))
			dfdx.Set(0, 1, -1.0)
			dfdx.Set(1, 0, -1.0)
			dfdx.Set(1, 1, 2.0+math.Exp(-x[1]))
		}

	case 2:
		name = "trigonometric and exponential"
		xTrial = []float64{0.4, 3.0}
		xReference = []float64{-0.2605992900257, 0.6225308965998}
		funcF = func(fx, x la.Vector) {
			pi := math.Pi
			e := math.E
			fx[0] = 0.5*sin(x[0]*x[1]) - 0.25*x[1]/pi - 0.5*x[0]
			fx[1] = (1.0-0.25/pi)*(math.Exp(2.0*x[0])-e) + e*x[1]/pi - 2.0*e*x[0]
			return
		}
		funcJsparse = func(dfdx *la.Triplet, x la.Vector) {
			pi := math.Pi
			e := math.E
			dfdx.Start()
			dfdx.Put(0, 0, 0.5*x[1]*cos(x[0]*x[1])-0.5)
			dfdx.Put(0, 1, 0.5*x[0]*cos(x[0]*x[1])-0.25/pi)
			dfdx.Put(1, 0, (2.0-0.5/pi)*math.Exp(2.0*x[0])-2.0*e)
			dfdx.Put(1, 1, e/pi)
			return
		}
		funcJdense = func(dfdx *la.Matrix, x la.Vector) {
			pi := math.Pi
			e := math.E
			dfdx.Set(0, 0, 0.5*x[1]*cos(x[0]*x[1])-0.5)
			dfdx.Set(0, 1, 0.5*x[0]*cos(x[0]*x[1])-0.25/pi)
			dfdx.Set(1, 0, (2.0-0.5/pi)*math.Exp(2.0*x[0])-2.0*e)
			dfdx.Set(1, 1, e/pi)
			return
		}

	default:
		chk.Panic("there is no problem index = %d", index)
	}
	return
}

func checkProblem(tst *testing.T, sol *NlSolver, x, xRef la.Vector, funcF fun.Vv, tolx, tolf float64, checkJacobian bool) {
	// check solution
	fx := make([]float64, len(x))
	funcF(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, xRef)
	io.Pf("f(x) = %v\n", fx)
	chk.Array(tst, "x  ", tolx, x, xRef)
	chk.Array(tst, "f(x) = 0? ", tolf, fx, nil)

	// check Jacobian
	if checkJacobian {
		io.Pforan("\nchecking Jacobian @ %v\n", x)
		sol.CheckJ(x, 1e-5, chk.Verbose)
	}
}

func solveProblem(tst *testing.T, index int, xTrialOrNil, xRefOrNil la.Vector,
	denseMatrix bool, numJacobian bool, lineSearch bool,
	tolx, tolf float64) {

	// problem data
	name, xTrial, xReference, funcF, funcJsparse, funcJdense := problem(index)

	// title
	anaOrNum := "ana"
	if numJacobian {
		anaOrNum = "num"
	}
	sparseOrDense := "sparse"
	if denseMatrix {
		sparseOrDense = "dense"
	}
	lineSearchOrNot := ""
	if lineSearch {
		lineSearchOrNot = "/line"
	}
	chk.PrintTitle(io.Sf("NlSolver %d (%s/%s%s): %s", index, anaOrNum, sparseOrDense, lineSearchOrNot, name))

	// init
	sol := NewNlSolver(len(xTrial), funcF)
	defer sol.Free()

	// config
	sol.config.Verbose = chk.Verbose
	sol.config.LineSearch = lineSearch
	if !numJacobian {
		if denseMatrix {
			sol.SetJacobianFunction(nil, funcJdense)
		} else {
			sol.SetJacobianFunction(funcJsparse, nil)
		}
	}

	// alternative xTrial and xReference
	x := xTrial
	xRef := xReference
	if xTrialOrNil != nil {
		x = xTrialOrNil
	}
	if xRefOrNil != nil {
		xRef = xRefOrNil
	}

	// solve non linear system
	sol.Solve(x)
	checkProblem(tst, sol, x, xRef, funcF, tolx, tolf, true)
}

// ------ problem 0

func TestNlSolver0AnaSparseNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 0, nil, nil, false, false, false, 1e-15, 1e-15)
}

func TestNlSolver0AnaSparseLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 0, nil, nil, false, false, true, 1e-15, 1e-15)
}

func TestNlSolver0AnaDenseNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 0, nil, nil, true, false, false, 1e-15, 1e-15)
}

func TestNlSolver0AnaDenseLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 0, nil, nil, true, false, true, 1e-15, 1e-15)
}

func TestNlSolver0NumNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 0, nil, nil, false, true, false, 1e-13, 1e-13)
}

func TestNlSolver0NumLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 0, nil, nil, false, true, true, 1e-13, 1e-13)
}

// ------ problem 1

func TestNlSolver1AnaSparseNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 1, nil, nil, false, false, false, 1e-4, 1e-14)
}

func TestNlSolver1AnaSparseLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 1, nil, nil, false, false, true, 1e-4, 1e-14)
}

func TestNlSolver1AnaDenseNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 1, nil, nil, true, false, false, 1e-4, 1e-14)
}

func TestNlSolver1AnaDenseLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 1, nil, nil, true, false, true, 1e-4, 1e-14)
}

func TestNlSolver1NumNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 1, nil, nil, false, true, false, 1e-4, 1e-14)
}

func TestNlSolver1NumLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 1, nil, nil, false, true, true, 1e-4, 1e-14)
}

// ------ problem 2

func TestNlSolver2AnaSparseNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 2, nil, nil, false, false, false, 1e-13, 1e-11)
	solveProblem(tst, 2, []float64{0.7, 4.0}, []float64{0.5000000377836, 3.1415927055406}, false, false, false, 1e-7, 1e-14)
	solveProblem(tst, 2, []float64{1.0, 4.0}, []float64{1.65458271876435, -15.819188232171314}, false, false, false, 1e-15, 1e-14)
}

/* this problem does not work with line-search
func TestNlSolver2AnaSparseLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 2, nil, nil, false, false, true, 1e-13, 1e-11)
}
*/

func TestNlSolver2AnaDenseNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 2, nil, nil, true, false, false, 1e-13, 1e-11)
	solveProblem(tst, 2, []float64{0.7, 4.0}, []float64{0.5000000377836, 3.1415927055406}, true, false, false, 1e-7, 1e-14)
	solveProblem(tst, 2, []float64{1.0, 4.0}, []float64{1.65458271876435, -15.819188232171314}, true, false, false, 1e-14, 1e-14)
}

func TestNlSolver2NumNoLine(tst *testing.T) {
	// verbose()
	solveProblem(tst, 2, nil, nil, false, true, false, 1e-13, 1e-11)
	solveProblem(tst, 2, []float64{0.7, 4.0}, []float64{0.5000000377836, 3.1415927055406}, false, true, false, 1e-7, 1e-14)
	solveProblem(tst, 2, []float64{1.0, 4.0}, []float64{1.65458271876435, -15.819188232171314}, false, true, false, 1e-15, 1e-14)
}
