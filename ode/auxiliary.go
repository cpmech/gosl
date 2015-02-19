// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
)

// callbacks
type Cb_ycorr func(y []float64, x float64, args ...interface{}) // y(x) correct

// plot results corresponding to one run
func Plot(dirout, fnkey, method string, bres *bytes.Buffer, ycps []int, ndim int, ycfcn Cb_ycorr, xa, xb float64,
	withdx, show bool, extra string) {

	// save file with results
	os.MkdirAll(dirout, 0777)
	if bres != nil {
		io.WriteFileD(dirout, fnkey+".res", bres)
	}

	// new python script
	var b bytes.Buffer
	fmt.Fprintf(&b, "from gosl import *\n")
	fmt.Fprintf(&b, "d = Read('%s/%s.res')\n", dirout, fnkey)

	// closed-form solution
	var xc []float64
	if ycfcn != nil {
		np := 101
		dx := (xb - xa) / float64(np-1)
		fmt.Fprintf(&b, "yc = array([\n")
		xc = make([]float64, np)
		yc := make([]float64, ndim)
		for i := 0; i < np; i++ {
			xc[i] = xa + dx*float64(i)
			ycfcn(yc, xc[i])
			fmt.Fprintf(&b, "[")
			for j := 0; j < ndim; j++ {
				if j == ndim-1 {
					fmt.Fprintf(&b, "%g", yc[j])
				} else {
					fmt.Fprintf(&b, "%g,", yc[j])
				}
			}
			if i == np-1 {
				fmt.Fprintf(&b, "]")
			} else {
				fmt.Fprintf(&b, "],\n")
			}
		}
		fmt.Fprintf(&b, "])\n")
		fmt.Fprintf(&b, "xc = array([")
		for i := 0; i < np; i++ {
			if i == np-1 {
				fmt.Fprintf(&b, "%g", xc[i])
			} else {
				fmt.Fprintf(&b, "%g,", xc[i])
			}
		}
		fmt.Fprintf(&b, "])\n")
	}

	// number of subplots
	nplt := len(ycps)
	if withdx {
		nplt += 1
	}

	// plot
	for i, cp := range ycps {
		fmt.Fprintf(&b, "subplot(%d,1,%d)\n", nplt, i+1)
		if ycfcn != nil {
			fmt.Fprintf(&b, "plot(xc, yc[:,%d], 'y-', lw=6, clip_on=0, label='solution')\n", cp)
		}
		fmt.Fprintf(&b, "plot(d['x'], d['y%d'], 'b-', marker='.', lw=1, clip_on=0, label='%s')\n", cp, method)
		fmt.Fprintf(&b, "Gll('x', 'y%d')\n", cp)
	}
	if withdx {
		fmt.Fprintf(&b, "subplot(%d,1,%d)\n", nplt, len(ycps)+1)
		fmt.Fprintf(&b, "plot(d['x'], d['dx'], 'b-', marker='.', lw=1, clip_on=0, label='%s')\n", method)
		fmt.Fprintf(&b, "gca().set_yscale('log')\n")
		fmt.Fprintf(&b, "Gll('x', 'step size')\n")
	}
	fmt.Fprintf(&b, extra)
	fmt.Fprintf(&b, "show()\n")

	// write file
	fn := fmt.Sprintf("%s/%s.py", dirout, fnkey)
	io.WriteFile(fn, &b)

	// run script
	if show {
		_, err := exec.Command("python", fn).Output()
		if err != nil {
			chk.Panic("failed when calling python %s\n%v", fn, err)
		}
	}
}

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
	var o ODE
	o.Init(method, ndim, fcn, jac, M, out, true)
	o.PredCtrl = true

	// for python script
	var b0, b1, b2, b3 bytes.Buffer
	fmt.Fprintf(&b0, "from gosl import *\n")
	fmt.Fprintf(&b0, "tols = array([")
	fmt.Fprintf(&b1, "errs = array([")
	fmt.Fprintf(&b2, "nfev = array([")

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
			fmt.Fprintf(&b0, "%g", tol)
			fmt.Fprintf(&b1, "%g", re.value)
			fmt.Fprintf(&b2, "%d", o.nfeval)
		} else {
			fmt.Fprintf(&b0, "%g,", tol)
			fmt.Fprintf(&b1, "%g,", re.value)
			fmt.Fprintf(&b2, "%d,", o.nfeval)
		}
	}
	fmt.Fprintf(&b0, "])\n")
	fmt.Fprintf(&b1, "])\n")
	fmt.Fprintf(&b2, "], dtype=float)\n")

	// python script
	fmt.Fprintf(&b3, "X, Y = -log10(errs), log10(nfev)\n")
	fmt.Fprintf(&b3, "plot(X, Y, clip_on=0)\n")
	if len(orders) > 0 {
		fmt.Fprintf(&b3, "dX = X[-1] - X[0]\n")
	}
	for _, ord := range orders {
		fmt.Fprintf(&b3, "plot([X[0], X[0]+dX], [Y[0], Y[0] + dX/float(%g)], 'k--', clip_on=0)\n", ord)
	}
	fmt.Fprintf(&b3, "Gll('correctness = -log10(error)', 'work = log10(nfev)', leg=0)\n")
	fmt.Fprintf(&b3, "show()\n")

	// write file
	fnpath := fmt.Sprintf("%s/%s_wc.py", dirout, fnkey)
	io.WriteFileD(dirout, fnkey+"_wc.py", &b0, &b1, &b2, &b3)

	// run script
	if show {
		_, err := exec.Command("python", fnpath).Output()
		if err != nil {
			chk.Panic("failed when calling python %s\n%v", fnpath, err)
		}
	}
}
