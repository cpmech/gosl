// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"bytes"
	"math"
	"os/exec"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// work/correctness analysis
func WcAnalysis(dirout, fnkey, method string, fcn Cb_fcn, jac Cb_jac, M *la.Triplet, ycfcn Cb_ycorr, ya []float64, xa, xb float64,
	orders []float64, show bool) {

	// structure holding error data
	type RmsErr struct {
		value float64
		count int
	}

	// output function => calculate error indicator
	ndim := len(ya)
	yc := make([]float64, ndim)
	out := func(first bool, dx, x float64, y []float64, args ...interface{}) error {
		ycfcn(yc, x)
		for i := 0; i < ndim; i++ {
			args[0].(*RmsErr).value += math.Pow(math.Abs(y[i]-yc[i])/(1.0+y[i]), 2.0)
		}
		args[0].(*RmsErr).count += 1
		return nil
	}

	// initialise ode
	var o Solver
	o.Init(method, ndim, fcn, jac, M, out, true)
	o.PredCtrl = true

	// for python script
	var b0, b1, b2, b3 bytes.Buffer
	io.Ff(&b0, "from gosl import *\n")
	io.Ff(&b0, "tols = array([")
	io.Ff(&b1, "errs = array([")
	io.Ff(&b2, "nfev = array([")

	// run for a number of tolerances
	nt := 13
	tols := make([]float64, nt)
	for i := 1; i < nt+1; i++ {
		tols[i-1] = math.Pow(10.0, -float64(i))
	}
	io.Pf("tols = %v\n", tols)
	y := make([]float64, ndim)
	for i, tol := range tols {

		// set tolerances
		o.Atol, o.Rtol = tol, tol
		//o.Atol, o.Rtol = 1.0e-9, tol
		o.NmaxSS = 10000

		// run
		copy(y, ya)
		var re RmsErr
		o.Solve(y, xa, xb, xb-xa, false, &re)
		re.value = math.Sqrt(re.value / float64(re.count))
		io.Pf("tol = %e  =>  err = %e  =>  feval = %d\n", tol, re.value, o.nfeval)

		// python script
		if i == len(tols)-1 {
			io.Ff(&b0, "%g", tol)
			io.Ff(&b1, "%g", re.value)
			io.Ff(&b2, "%d", o.nfeval)
		} else {
			io.Ff(&b0, "%g,", tol)
			io.Ff(&b1, "%g,", re.value)
			io.Ff(&b2, "%d,", o.nfeval)
		}
	}
	io.Ff(&b0, "])\n")
	io.Ff(&b1, "])\n")
	io.Ff(&b2, "], dtype=float)\n")

	// python script
	io.Ff(&b3, "X, Y = -log10(errs), log10(nfev)\n")
	io.Ff(&b3, "plot(X, Y, clip_on=0)\n")
	if len(orders) > 0 {
		io.Ff(&b3, "dX = X[-1] - X[0]\n")
	}
	for _, ord := range orders {
		io.Ff(&b3, "plot([X[0], X[0]+dX], [Y[0], Y[0] + dX/float(%g)], 'k--', clip_on=0)\n", ord)
	}
	io.Ff(&b3, "Gll('correctness = -log10(error)', 'work = log10(nfev)', leg=0)\n")
	io.Ff(&b3, "show()\n")

	// write file
	fnpath := io.Sf("%s/%s_wc.py", dirout, fnkey)
	io.WriteFileD(dirout, fnkey+"_wc.py", &b0, &b1, &b2, &b3)

	// run script
	if show {
		_, err := exec.Command("python", fnpath).Output()
		if err != nil {
			chk.Panic("failed when calling python %s\n%v", fnpath, err)
		}
	}
}
