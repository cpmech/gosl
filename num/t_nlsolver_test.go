// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package num

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fdm"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

func Test_nls01(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	//utl.UseColors = false
	chk.PrintTitle("nls01. 2 eqs system")

	ffcn := func(fx, x []float64) error {
		fx[0] = math.Pow(x[0], 3.0) + x[1] - 1.0
		fx[1] = -x[0] + math.Pow(x[1], 3.0) + 1.0
		return nil
	}
	Jfcn := func(dfdx *la.Triplet, x []float64) error {
		dfdx.Start()
		dfdx.Put(0, 0, 3.0*x[0]*x[0])
		dfdx.Put(0, 1, 1.0)
		dfdx.Put(1, 0, -1.0)
		dfdx.Put(1, 1, 3.0*x[1]*x[1])
		return nil
	}

	x := []float64{0.5, 0.5}
	atol, rtol, ftol := 1e-10, 1e-10, 10*EPS
	fx := make([]float64, len(x))
	neq := len(x)

	prms := map[string]float64{
		"atol":    atol,
		"rtol":    rtol,
		"ftol":    ftol,
		"lSearch": 1.0,
	}

	io.PfYel("\n-------------------- Analytical Jacobian -------------------\n")

	// init
	var nls_ana NlSolver
	nls_ana.Init(neq, ffcn, Jfcn, nil, false, false, prms)
	defer nls_ana.Clean()

	// solve
	err := nls_ana.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{1.0, 0.0})
	io.Pf("f(x) = %v\n", fx)
	chk.Vector(tst, "f(x) = 0? ", 1e-16, fx, []float64{})

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	_, err = nls_ana.CheckJ(x, 1e-5, false, true)
	if err != nil {
		chk.Panic(err.Error())
	}

	io.PfYel("\n\n-------------------- Numerical Jacobian --------------------\n")
	xx := []float64{0.5, 0.5}

	// init
	var nls_num NlSolver
	nls_num.Init(neq, ffcn, nil, nil, false, true, prms)
	defer nls_num.Clean()

	// solve
	err = nls_num.Solve(xx, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	ffcn(fx, xx)
	io.Pf("xx    = %v  expected = %v\n", xx, []float64{1.0, 0.0})
	io.Pf("f(xx) = %v\n", fx)
	chk.Vector(tst, "f(x) = 0? ", 1e-16, fx, []float64{})
	chk.Vector(tst, "x == xx", 1e-15, x, xx)

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	_, err = nls_ana.CheckJ(x, 1e-5, false, true)
	if err != nil {
		chk.Panic(err.Error())
	}
}

func Test_nls02(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	//utl.UseColors = false
	chk.PrintTitle("nls02. 2 eqs system with exp function")

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

	io.PfYel("\n-------------------- Analytical Jacobian -------------------\n")

	// init
	var nls_ana NlSolver
	nls_ana.Init(neq, ffcn, Jfcn, nil, false, false, prms)
	defer nls_ana.Clean()

	// solve
	err := nls_ana.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5671, 0.5671})
	io.Pf("f(x) = %v\n", fx)
	chk.Vector(tst, "f(x) = 0? ", 1e-15, fx, []float64{})

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	_, err = nls_ana.CheckJ(x, 1e-5, false, true)
	if err != nil {
		chk.Panic(err.Error())
	}

	io.PfYel("\n\n-------------------- Numerical Jacobian --------------------\n")
	xx := []float64{5.0, 5.0}

	// init
	var nls_num NlSolver
	nls_num.Init(neq, ffcn, nil, nil, false, true, prms)
	defer nls_num.Clean()

	// solve
	err = nls_num.Solve(xx, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	ffcn(fx, x)
	io.Pf("xx    = %v  expected = %v\n", x, []float64{0.5671, 0.5671})
	io.Pf("f(xx) = %v\n", fx)
	chk.Vector(tst, "f(x) = 0? ", 1e-15, fx, []float64{})
	chk.Vector(tst, "x == xx", 1e-15, x, xx)

	// check Jacobian
	io.Pforan("\nchecking Jacobian @ %v\n", x)
	_, err = nls_ana.CheckJ(x, 1e-5, false, true)
	if err != nil {
		chk.Panic(err.Error())
	}
}

func Test_nls03(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("nls03. 2 eqs system with trig functions")

	e := math.E
	ffcn := func(fx, x []float64) error {
		fx[0] = 0.5*sin(x[0]*x[1]) - 0.25*x[1]/pi - 0.5*x[0]
		fx[1] = (1.0-0.25/pi)*(math.Exp(2.0*x[0])-e) + e*x[1]/pi - 2.0*e*x[0]
		return nil
	}
	Jfcn := func(dfdx *la.Triplet, x []float64) error {
		dfdx.Start()
		dfdx.Put(0, 0, 0.5*x[1]*cos(x[0]*x[1])-0.5)
		dfdx.Put(0, 1, 0.5*x[0]*cos(x[0]*x[1])-0.25/pi)
		dfdx.Put(1, 0, (2.0-0.5/pi)*math.Exp(2.0*x[0])-2.0*e)
		dfdx.Put(1, 1, e/pi)
		return nil
	}
	JfcnD := func(dfdx [][]float64, x []float64) error {
		dfdx[0][0] = 0.5*x[1]*cos(x[0]*x[1]) - 0.5
		dfdx[0][1] = 0.5*x[0]*cos(x[0]*x[1]) - 0.25/pi
		dfdx[1][0] = (2.0-0.5/pi)*math.Exp(2.0*x[0]) - 2.0*e
		dfdx[1][1] = e / pi
		return nil
	}

	x := []float64{0.4, 3.0}
	fx := make([]float64, len(x))
	atol := 1e-6
	rtol := 1e-3
	ftol := 10 * EPS
	neq := len(x)

	prms := map[string]float64{
		"atol":    atol,
		"rtol":    rtol,
		"ftol":    ftol,
		"lSearch": 0.0, // does not work with line search
	}

	// init
	var nls_sps NlSolver // sparse
	var nls_den NlSolver // dense
	nls_sps.Init(neq, ffcn, Jfcn, nil, false, false, prms)
	nls_den.Init(neq, ffcn, nil, JfcnD, true, false, prms)
	defer nls_sps.Clean()
	defer nls_den.Clean()

	io.PfMag("\n/////////////////////// sparse //////////////////////////////////////////\n")

	x = []float64{0.4, 3.0}
	io.PfYel("\n--- sparse ------------- with x = %v --------------\n", x)
	err := nls_sps.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{-0.2605992900257, 0.6225308965998})
	io.Pf("f(x) = %v\n", fx)
	chk.Vector(tst, "x", 1e-13, x, []float64{-0.2605992900257, 0.6225308965998})
	chk.Vector(tst, "f(x) = 0? ", 1e-11, fx, nil)

	x = []float64{0.7, 4.0}
	io.PfYel("\n--- sparse ------------- with x = %v --------------\n", x)
	//rtol = 1e-2
	err = nls_sps.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5000000377836, 3.1415927055406})
	io.Pf("f(x) = %v\n", fx)
	chk.Vector(tst, "x  ", 1e-7, x, []float64{0.5000000377836, 3.1415927055406})
	chk.Vector(tst, "f(x) = 0? ", 1e-7, fx, nil)

	x = []float64{1.0, 4.0}
	io.PfYel("\n--- sparse ------------- with x = %v ---------------\n", x)
	//lSearch, chkConv := false, true  // this combination fails due to divergence
	//lSearch, chkConv := false, false // this combination works but results are different
	//lSearch, chkConv := true, true   // this combination works but results are wrong => fails
	nls_sps.ChkConv = false
	err = nls_sps.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5, pi})
	io.Pf("f(x) = %v << converges to a different solution\n", fx)
	chk.Vector(tst, "f(x) = 0? ", 1e-8, fx, nil)

	io.PfMag("\n/////////////////////// dense //////////////////////////////////////////\n")

	x = []float64{0.4, 3.0}
	io.PfYel("\n--- dense ------------- with x = %v --------------\n", x)
	err = nls_den.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{-0.2605992900257, 0.6225308965998})
	io.Pf("f(x) = %v\n", fx)
	chk.Vector(tst, "x", 1e-13, x, []float64{-0.2605992900257, 0.6225308965998})
	chk.Vector(tst, "f(x) = 0? ", 1e-11, fx, nil)

	x = []float64{0.7, 4.0}
	io.PfYel("\n--- dense ------------- with x = %v --------------\n", x)
	//rtol = 1e-2
	err = nls_den.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5000000377836, 3.1415927055406})
	io.Pf("f(x) = %v\n", fx)
	chk.Vector(tst, "x  ", 1e-7, x, []float64{0.5000000377836, 3.1415927055406})
	chk.Vector(tst, "f(x) = 0? ", 1e-7, fx, nil)

	x = []float64{1.0, 4.0}
	io.PfYel("\n--- dense ------------- with x = %v ---------------\n", x)
	nls_den.ChkConv = false
	err = nls_den.Solve(x, false)
	if err != nil {
		chk.Panic(err.Error())
	}
	ffcn(fx, x)
	io.Pf("x    = %v  expected = %v\n", x, []float64{0.5, pi})
	io.Pf("f(x) = %v << converges to a different solution\n", fx)
	chk.Vector(tst, "f(x) = 0? ", 1e-8, fx, nil)
}

func Test_nls04(tst *testing.T) {

	prevTs := verbose()
	defer func() {
		verbose() = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//verbose()
	chk.PrintTitle("nls04. finite differences problem")

	// grid
	var g fdm.Grid2D
	g.Init(1.0, 1.0, 6, 6)

	// equations numbering
	var e fdm.Equations
	peq := utl.IntUnique(g.L, g.R, g.B, g.T)
	e.Init(g.N, peq)

	// K11 and K12
	var K11, K12 la.Triplet
	fdm.InitK11andK12(&K11, &K12, &e)

	// assembly
	F1 := make([]float64, e.N1)
	fdm.Assemble(&K11, &K12, F1, nil, &g, &e)

	// prescribed values
	U2 := make([]float64, e.N2)
	for _, eq := range g.L {
		U2[e.FR2[eq]] = 50.0
	}
	for _, eq := range g.R {
		U2[e.FR2[eq]] = 0.0
	}
	for _, eq := range g.B {
		U2[e.FR2[eq]] = 0.0
	}
	for _, eq := range g.T {
		U2[e.FR2[eq]] = 50.0
	}

	// functions
	k11 := K11.ToMatrix(nil)
	k12 := K12.ToMatrix(nil)
	ffcn := func(fU1, U1 []float64) error { // K11*U1 + K12*U2 - F1
		la.VecCopy(fU1, -1, F1)            // fU1 := (-F1)
		la.SpMatVecMulAdd(fU1, 1, k11, U1) // fU1 += K11*U1
		la.SpMatVecMulAdd(fU1, 1, k12, U2) // fU1 += K12*U2
		return nil
	}
	Jfcn := func(dfU1dU1 *la.Triplet, U1 []float64) error {
		fdm.Assemble(dfU1dU1, &K12, F1, nil, &g, &e)
		return nil
	}
	JfcnD := func(dfU1dU1 [][]float64, U1 []float64) error {
		la.MatCopy(dfU1dU1, 1, K11.ToMatrix(nil).ToDense())
		return nil
	}

	prms := map[string]float64{
		"atol":    1e-8,
		"rtol":    1e-8,
		"ftol":    1e-12,
		"lSearch": 0.0,
	}

	// init
	var nls_sps NlSolver // sparse analytical
	var nls_num NlSolver // sparse numerical
	var nls_den NlSolver // dense analytical
	nls_sps.Init(e.N1, ffcn, Jfcn, nil, false, false, prms)
	nls_num.Init(e.N1, ffcn, nil, nil, false, true, prms)
	nls_den.Init(e.N1, ffcn, nil, JfcnD, true, false, prms)
	defer nls_sps.Clean()
	defer nls_num.Clean()
	defer nls_den.Clean()

	// results
	U1sps := make([]float64, e.N1)
	U1num := make([]float64, e.N1)
	U1den := make([]float64, e.N1)
	Usps := make([]float64, e.N)
	Unum := make([]float64, e.N)
	Uden := make([]float64, e.N)

	// solution
	Uc := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 50.0, 25.0, 325.0 / 22.0, 100.0 / 11.0, 50.0 / 11.0,
		0.0, 50.0, 775.0 / 22.0, 25.0, 375.0 / 22.0, 100.0 / 11.0, 0.0, 50.0, 450.0 / 11.0, 725.0 / 22.0,
		25.0, 325.0 / 22.0, 0.0, 50.0, 500.0 / 11.0, 450.0 / 11.0, 775.0 / 22.0, 25.0, 0.0, 50.0, 50.0,
		50.0, 50.0, 50.0, 50.0,
	}

	io.PfYel("\n---- sparse -------- Analytical Jacobian -------------------\n")

	// solve
	err := nls_sps.Solve(U1sps, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	fdm.JoinVecs(Usps, U1sps, U2, &e)
	chk.Vector(tst, "Usps", 1e-14, Usps, Uc)

	// plot
	if false {
		g.Contour("results", "fig_t_heat_square", nil, Usps, 11, false)
	}

	io.PfYel("\n---- dense -------- Analytical Jacobian -------------------\n")

	// solve
	err = nls_den.Solve(U1den, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	fdm.JoinVecs(Uden, U1den, U2, &e)
	chk.Vector(tst, "Uden", 1e-14, Uden, Uc)

	io.PfYel("\n---- sparse -------- Numerical Jacobian -------------------\n")

	// solve
	err = nls_num.Solve(U1num, false)
	if err != nil {
		chk.Panic(err.Error())
	}

	// check
	fdm.JoinVecs(Unum, U1num, U2, &e)
	chk.Vector(tst, "Unum", 1e-14, Unum, Uc)
}
