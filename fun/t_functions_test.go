// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"
	"math/big"
	"math/cmplx"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func Test_functions01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("functions01")

	x := utl.LinSpace(-2, 2, 21)
	ym := make([]float64, len(x))
	yh := make([]float64, len(x))
	ys := make([]float64, len(x))
	yAbs2Ramp := make([]float64, len(x))
	yHea2Ramp := make([]float64, len(x))
	ySig2Heav := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		ym[i] = Ramp(x[i])
		yh[i] = Heav(x[i])
		ys[i] = Sign(x[i])
		yAbs2Ramp[i] = (x[i] + math.Abs(x[i])) / 2.0
		yHea2Ramp[i] = x[i] * yh[i]
		ySig2Heav[i] = (1.0 + ys[i]) / 2.0
	}
	chk.Array(tst, "abs => ramp", 1e-17, ym, yAbs2Ramp)
	chk.Array(tst, "hea => ramp", 1e-17, ym, yHea2Ramp)
	chk.Array(tst, "sig => heav", 1e-17, yh, ySig2Heav)
}

// numderiv employs a 1st order forward difference to approximate the derivative of f(x) w.r.t x @ x
func numderiv(f func(x float64) float64, x float64) float64 {
	eps, cte1 := 1e-16, 1e-5
	delta := math.Sqrt(eps * utl.Max(cte1, math.Abs(x)))
	return (f(x+delta) - f(x)) / delta
}

func Test_functions02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("functions02")

	β := 6.0
	f := func(x float64) float64 { return Sramp(x, β) }
	ff := func(x float64) float64 { return SrampD1(x, β) }

	np := 101
	//x  := utl.LinSpace(-5e5, 5e5, np)
	//x  := utl.LinSpace(-5e2, 5e2, np)
	x := utl.LinSpace(-5e1, 5e1, np)
	y := make([]float64, np)
	g := make([]float64, np)
	h := make([]float64, np)
	tolg, tolh := 1e-6, 1e-5
	withErr := false
	for i := 0; i < np; i++ {
		y[i] = Sramp(x[i], β)
		g[i] = SrampD1(x[i], β)
		h[i] = SrampD2(x[i], β)
		gnum := numderiv(f, x[i])
		hnum := numderiv(ff, x[i])
		errg := math.Abs(g[i] - gnum)
		errh := math.Abs(h[i] - hnum)
		clrg, clrh := "[1;32m", "[1;32m"
		if errg > tolg {
			clrg, withErr = "[1;31m", true
		}
		if errh > tolh {
			clrh, withErr = "[1;31m", true
		}
		io.Pf("errg = %s%23.15e   errh = %s%23.15e[0m\n", clrg, errg, clrh, errh)
	}

	if withErr {
		chk.Panic("errors found")
	}
}

func TestLogistic01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Logistic01")

	f := func(x float64) float64 { return Logistic(x) }

	np := 101
	//x := utl.LinSpace(-5e5, 5e5, np)
	//x := utl.LinSpace(-5e2, 5e2, np)
	x := utl.LinSpace(-5e1, 5e1, np)
	y := make([]float64, np)
	g := make([]float64, np)
	tolg := 1e-6
	withErr := false
	for i := 0; i < np; i++ {
		y[i] = Logistic(x[i])
		g[i] = LogisticD1(x[i])
		gnum := numderiv(f, x[i])
		errg := math.Abs(g[i] - gnum)
		clrg := "[1;32m"
		if errg > tolg {
			clrg, withErr = "[1;31m", true
		}
		io.Pf("errg = %s%23.15e[0m\n", clrg, errg)
	}

	if withErr {
		chk.Panic("errors found")
	}

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.25})
		plt.Subplot(2, 1, 1)
		plt.Plot(x, y, &plt.A{C: plt.C(0, 0), Ls: "-", Lw: 1.5, L: "logistic", NoClip: true})
		plt.Gll("z", "g(z)", nil)
		plt.Subplot(2, 1, 2)
		plt.Plot(x, g, &plt.A{C: plt.C(2, 0), Ls: "-", Lw: 1.5, L: "derivative", NoClip: true})
		plt.Gll("z", "dgdz(z)", nil)
		plt.Save("/tmp/gosl/fun", "logistic01")
	}
}

func Test_functions03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("functions03")

	eps := 1e-2
	f := func(x float64) float64 { return Sabs(x, eps) }
	ff := func(x float64) float64 { return SabsD1(x, eps) }

	np := 101
	//x  := utl.LinSpace(-5e5, 5e5, np)
	//x  := utl.LinSpace(-5e2, 5e2, np)
	x := utl.LinSpace(-5e1, 5e1, np)
	Y := make([]float64, np)
	y := make([]float64, np)
	g := make([]float64, np)
	h := make([]float64, np)
	tolg, tolh := 1e-6, 1e-5
	withErr := false
	for i := 0; i < np; i++ {
		Y[i] = math.Abs(x[i])
		y[i] = Sabs(x[i], eps)
		g[i] = SabsD1(x[i], eps)
		h[i] = SabsD2(x[i], eps)
		gnum := numderiv(f, x[i])
		hnum := numderiv(ff, x[i])
		errg := math.Abs(g[i] - gnum)
		errh := math.Abs(h[i] - hnum)
		clrg, clrh := "[1;32m", "[1;32m"
		if errg > tolg {
			clrg, withErr = "[1;31m", true
		}
		if errh > tolh {
			clrh, withErr = "[1;31m", true
		}
		io.Pf("errg = %s%23.15e   errh = %s%23.15e[0m\n", clrg, errg, clrh, errh)
	}

	if withErr {
		chk.Panic("errors found")
	}

	if false {
		//if true {
		plt.Subplot(3, 1, 1)
		plt.Plot(x, y, &plt.A{C: "k", Ls: "--", L: "abs"})
		plt.Plot(x, y, &plt.A{C: "b", Ls: "-", L: "sabs"})
		plt.Gll("x", "y", nil)

		plt.Subplot(3, 1, 2)
		plt.Plot(x, g, &plt.A{C: "b", Ls: "-", L: "sabs"})
		plt.Gll("x", "dy/dx", nil)

		plt.Subplot(3, 1, 3)
		plt.Plot(x, h, &plt.A{C: "b", Ls: "-", L: "sabs"})
		plt.Gll("x", "d2y/dx2", nil)

		plt.Show()
	}
}

func Test_suq01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("suq01. superquadric functions")

	if chk.Verbose {
		np := 101
		X := utl.LinSpace(0, math.Pi, np)
		Y := make([]float64, np)
		for i := 0; i < np; i++ {
			Y[i] = SuqCos(X[i], 4)
		}
		plt.Plot(X, Y, nil)
		plt.Gll("x", "y", nil)
		plt.Save("/tmp/gosl", "t_suq01")
	}
}

func Test_factorial01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("factorial01. Factorial22")

	n0 := Factorial22(0)
	n1 := Factorial22(1)
	n2 := Factorial22(2)
	n3 := Factorial22(3)
	n10 := Factorial22(10)
	n22 := Factorial22(22)

	chk.Float64(tst, "0!", 1e-15, n0, 1)
	chk.Float64(tst, "1!", 1e-15, n1, 1)
	chk.Float64(tst, "2!", 1e-15, n2, 2)
	chk.Float64(tst, "3!", 1e-15, n3, 6)
	chk.Float64(tst, "10!", 1e-15, n10, 3628800)
	chk.Float64(tst, "22!", 1e-15, n22, 1124000727777607680000)

	// printing max int sizes, out of curiosity
	MaxUint := ^uint(0)
	MaxInt := int(MaxUint >> 1)
	MinInt := -MaxInt - 1
	io.Pl()
	io.Pf("MaxUint = %v  %v\n", MaxUint, uint64(math.MaxUint64))
	io.Pf("MaxInt  = %v  %v\n", MaxInt, math.MaxInt64)
	io.Pf("MinInt  = %v  %v\n", MinInt, math.MinInt64)
}

func Test_factorial02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("factorial02. Factorial100")

	values := []int{22, 23, 50, 100}
	answers := []string{ // from http://www.tsm-resources.com/alists/fact.html
		"1124000727777607680000",
		"25852016738884976640000",
		"30414093201713378043612608166064768844377641568960512000000000000",
		"93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000",
	}
	for idx, value := range values {

		// compute factorial using big.Int and convert to big.Float
		ibig := new(big.Int)
		fbig := new(big.Float)
		ibig.MulRange(1, int64(value))
		fbig.SetPrec(big.MaxPrec)
		fbig.SetInt(ibig)
		txt := fbig.Text('f', 0)
		chk.String(tst, txt, answers[idx])

		// compute factorial using Factorial100
		f := Factorial100(value)
		diff := new(big.Float)
		diff.SetPrec(big.MaxPrec)
		diff.Sub(fbig, &f)
		d, a := diff.Float64()
		chk.String(tst, a.String(), "Exact")
		chk.Float64(tst, "diff", 1e-15, d, 0)
	}
}

func Test_beta01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("beta01. Beta function")

	aValues := []float64{1, 3, 10}
	bValues := []float64{5, 2, 11, -0.5}
	answers := [][]float64{ // values from wxMaxima beta(a,b) function
		{1.0 / 5.0, 1.0 / 2.0, 1.0 / 11, -2.0},
		{1.0 / 105.0, 1.0 / 12.0, 1.0 / 858.0, -16.0 / 3.0},
		{1.0 / 10010.0, 1.0 / 110, 1.0 / 1847560, -131072.0 / 12155.0},
	}
	for i, a := range aValues {
		for j, b := range bValues {
			tol := 1e-15
			if i == 2 && j == 3 {
				tol = 5e-14
			}
			res := Beta(a, b)
			chk.Float64(tst, io.Sf("Beta(%2f,%2f)", a, b), tol, res, answers[i][j])
		}
	}
}

func Test_binomial01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("binomial01. binomial coefficient")

	aValues := []int{10, 22, 50}
	bValues := []int{5, 2, 10}
	answers := [][]float64{ // values from wxMaxima beta(a,b) function
		{252, 45, 1},
		{26334, 231, 646646},
		{2118760, 1225, 10272278170},
	}
	for i, a := range aValues {
		for j, b := range bValues {
			res := Binomial(a, b)
			ures := UintBinomial(uint64(a), uint64(b))
			chk.Float64(tst, io.Sf("Binomial(%2d,%2d)", a, b), 1e-15, res, answers[i][j])
			chk.Float64(tst, "ures", 1e-15, float64(ures), answers[i][j])
		}
	}

	r49 := Binomial(50, 49)     // k = n-1
	r26 := Binomial(50, 26)     // k > n-k
	u49 := UintBinomial(50, 49) // k = n-1
	u26 := UintBinomial(50, 26) // k > n-k
	chk.Float64(tst, "Binomial(50,49)", 1e-15, r49, 50)
	chk.Float64(tst, "Binomial(50,26)", 1e-15, r26, 121548660036300-1) // cannot get 121548660036300
	chk.Float64(tst, "UintBinomial(50,49)", 1e-15, float64(u49), 50)
	chk.Float64(tst, "UintBinomial(50,26)", 1e-15, float64(u26), 121548660036300)
	io.Pforan("r26 = %.1f (should be 121548660036300.0)\n", r26)
	io.Pforan("u26 = %v\n", u26)

	// The following test fails with overflow in UintBinomial and incorrect results in Binomial
	// We need to use math/big for these
	if false {
		n100k50 := Binomial(100, 50)
		u100k50 := UintBinomial(100, 50)
		n100k50maxima := 100891344545564193334812497256.0
		io.Pforan("Binomial(100,50) = %v\n", n100k50)
		io.Pforan("UintBinomial(100,50) = %v\n", u100k50)
		chk.Float64(tst, "Binomial(100,50)", 1e-15, n100k50, n100k50maxima)
	}
}

func Test_binomial02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("binomial02. binomial with real arguments")

	aValues := []float64{10, 22, 50}
	bValues := []float64{5, 2, 10}
	answers := [][]float64{ // values from wxMaxima beta(a,b) function
		{252, 45, 1},
		{26334, 231, 646646},
		{2118760, 1225, 10272278170},
	}
	for i, a := range aValues {
		for j, b := range bValues {
			res := Rbinomial(a, b)
			tol := 1e-15
			if i == 2 {
				switch j {
				case 0:
					tol = 1e-9
					break
				case 1:
					tol = 2.3e-13 // For aarch64
					break
				case 2:
					tol = 1e-5 // for MacBook Air
					break
				default:
					tst.Errorf("Unexpected key value: %d\n", j)
				}
			}
			chk.AnaNum(tst, io.Sf("Rbinomial(%g,%g)", a, b), tol, res, answers[i][j], chk.Verbose)
		}
	}
}

func TestEuler01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Euler01. Euler's formula")

	a := ExpPix(math.Pi)
	A := cmplx.Exp(complex(0, math.Pi))
	io.Pforan("exp(+i⋅π) = %v  (%v)\n", a, A)
	chk.Complex128(tst, "exp(+i⋅π) == -1", 1e-15, a, -1)
	chk.Complex128(tst, "a == A", 1e-17, a, A)

	b := ExpMix(math.Pi)
	B := cmplx.Exp(complex(0, -math.Pi))
	io.Pforan("exp(-i⋅π) = %v  (%v)\n", b, B)
	chk.Complex128(tst, "exp(-i⋅π) == -1", 1e-15, b, -1)
	chk.Complex128(tst, "b == B", 1e-17, b, B)

	c := ExpPix(1)
	C := cmplx.Exp(1i)
	io.Pforan("exp(i) = %v  (%v)\n", c, C)
	chk.Float64(tst, "real(exp(i))", 1e-17, real(c), math.Cos(1))
	chk.Float64(tst, "imag(exp(i))", 1e-17, imag(c), math.Sin(1))
	chk.Complex128(tst, "c == C", 1e-17, c, C)
}

func TestSinc01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Sinc01. sine cardinal function")

	chk.Float64(tst, "sinc(0)", 1e-17, Sinc(0), 1)
	chk.Float64(tst, "sinc(π)", 1e-16, Sinc(math.Pi), 0)
	chk.Float64(tst, "sinc(π/2)", 1e-17, Sinc(math.Pi/2), 2.0/math.Pi)
	chk.Float64(tst, "sinc(3π/2)", 1e-17, Sinc(3*math.Pi/2), -2.0/(3.0*math.Pi))

	if chk.Verbose {
		X := utl.LinSpace(-15, 15, 201)
		Y := utl.GetMapped(X, func(x float64) float64 { return Sinc(x) })
		plt.Reset(true, nil)
		plt.Plot(X, Y, &plt.A{C: "r", NoClip: true})
		plt.Gll("x", "sinc(x)", nil)
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "sinc01")
	}
}

func TestBoxcar01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Boxcar01. boxcar function")

	a, b := 0.5, 1.0

	chk.Float64(tst, "boxcar(a)", 1e-17, Boxcar(a, a, b), 0.5)
	chk.Float64(tst, "boxcar(b)", 1e-17, Boxcar(b, a, b), 0.5)
	chk.Float64(tst, "boxcar((a+b)/2)", 1e-17, Boxcar((a+b)/2, a, b), 1)
	chk.Float64(tst, "boxcar(a-1)", 1e-17, Boxcar(a-1, a, b), 0)
	chk.Float64(tst, "boxcar(b+1)", 1e-17, Boxcar(b+1, a, b), 0)
	for _, x := range []float64{-1, 0, 0.5, 0.7, 1.0, 1.5} {
		chk.Float64(tst, "H(x-a)-H(x-b)", 1e-17, Boxcar(x, a, b), Heav(x-a)-Heav(x-b))
	}

	if chk.Verbose {
		Xa := utl.LinSpace(-1.0, 1.5, 201)
		Xb := utl.LinSpace(-1.0, 1.5, 16)
		Ya := utl.GetMapped(Xa, func(x float64) float64 { return Boxcar(x, a, b) })
		Yb := utl.GetMapped(Xb, func(x float64) float64 { return Boxcar(x, a, b) })
		plt.Reset(true, nil)
		plt.Plot(Xa, Ya, &plt.A{C: "b", NoClip: true})
		plt.Plot(Xb, Yb, &plt.A{C: "r", Ls: "none", M: ".", NoClip: true})
		plt.Gll("x", "boxcar(x)", nil)
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "boxcar01")
	}
}

func TestRect01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Rect01. rectangular function")

	chk.Float64(tst, "rect(-0.5)", 1e-17, Rect(-0.5), 0.5)
	chk.Float64(tst, "rect(+0.5)", 1e-17, Rect(+0.5), 0.5)
	chk.Float64(tst, "rect(0)", 1e-17, Rect(0), 1)
	chk.Float64(tst, "rect(-0.7)", 1e-17, Rect(-0.7), 0)
	chk.Float64(tst, "rect(+0.7)", 1e-17, Rect(+0.7), 0)

	if chk.Verbose {
		Xa := utl.LinSpace(-1.5, 1.5, 201)
		Xb := utl.LinSpace(-1.5, 1.5, 16)
		Ya := utl.GetMapped(Xa, func(x float64) float64 { return Rect(x) })
		Yb := utl.GetMapped(Xb, func(x float64) float64 { return Rect(x) })
		plt.Reset(true, nil)
		plt.Plot(Xa, Ya, &plt.A{C: "b", NoClip: true})
		plt.Plot(Xb, Yb, &plt.A{C: "r", Ls: "none", M: ".", NoClip: true})
		plt.Gll("x", "rect(x)", nil)
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "rect01")
	}
}

func TestHat01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Hat01. 'hat' function")

	xc := 1.0
	y0 := 1.0
	h := 1.0
	l := 2.0

	chk.Float64(tst, "rect(-2.0)", 1e-17, Hat(-2.0, xc, y0, h, l), y0)
	chk.Float64(tst, "rect(-1.0)", 1e-17, Hat(-1.0, xc, y0, h, l), y0)
	chk.Float64(tst, "rect(+0.0)", 1e-17, Hat(+0.0, xc, y0, h, l), y0+h/2)
	chk.Float64(tst, "rect(+1.0)", 1e-17, Hat(+1.0, xc, y0, h, l), y0+h)
	chk.Float64(tst, "rect(+2.0)", 1e-17, Hat(+2.0, xc, y0, h, l), y0+h/2)
	chk.Float64(tst, "rect(+3.0)", 1e-17, Hat(+3.0, xc, y0, h, l), y0)
	chk.Float64(tst, "rect(+4.0)", 1e-17, Hat(+4.0, xc, y0, h, l), y0)

	Xb := utl.LinSpace(-2, 4, 13)
	for _, x := range Xb {
		chk.AnaNum(tst, io.Sf("Hat(%+.2f)", x), 1e-17, Hat(x, xc, y0, h, l), utl.Max(y0, y0+h-math.Abs(x-xc)/2), chk.Verbose)
	}

	Xb = utl.LinSpace(-2, 4, 8)
	for _, x := range Xb {
		chk.AnaNum(tst, io.Sf("Hat(%+.2f)", x), 1e-17, Hat(x, xc, y0, h, l), utl.Max(y0, y0+h-math.Abs(x-xc)/2), chk.Verbose)
		chk.DerivScaSca(tst, "HatD1", 1e-12, HatD1(x, xc, y0, h, l), x, 1e-3, chk.Verbose, func(t float64) float64 {
			return Hat(t, xc, y0, h, l)
		})
	}

	if chk.Verbose {
		Xa := utl.LinSpace(-2, 4, 201)
		Ya := utl.GetMapped(Xa, func(x float64) float64 { return Hat(x, xc, y0, h, l) })
		Yb := utl.GetMapped(Xb, func(x float64) float64 { return Hat(x, xc, y0, h, l) })
		plt.Reset(true, nil)
		plt.Plot(Xa, Ya, &plt.A{C: "b", NoClip: true})
		plt.Plot(Xb, Yb, &plt.A{C: "r", Ls: "none", M: ".", NoClip: true})
		plt.Equal()
		plt.Gll("x", "hat(x)", nil)
		plt.Cross(0, 0, nil)
		plt.HideAllBorders()
		plt.Save("/tmp/gosl/fun", "hat01")
	}
}

func TestNegOnePowN(tst *testing.T) {

	//verbose()
	chk.PrintTitle("NegOnePowN. calc (-1)ⁿ")

	for n := 0; n < 7; n++ {
		res := NegOnePowN(n)
		io.Pf("n=%d ⇒ (-1)ⁿ = %v\n", n, res)
		if n%2 == 0 { // even
			if res != 1 {
				tst.Errorf("pow(-1, %d) must be equal to 1\n", n)
				return
			}
		} else {
			if res != -1 {
				tst.Errorf("pow(+1, %d) must be equal to -1\n", n)
				return
			}
		}
		chk.Float64(tst, io.Sf("pow(-1,%d)", n), 1e-17, res, math.Pow(-1, float64(n)))
	}
}

func TestImagPowN(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ImagPowN. calc iⁿ; i.e. (√-1)ⁿ")

	for n := 0; n < 21; n++ {
		res := ImagPowN(n)
		tol := 1e-15
		if n > 10 {
			tol = 1e-14
		}
		chk.AnaNumC(tst, io.Sf("pow(i,%2d)", n), tol, res, cmplx.Pow(1i, complex(float64(n), 0)), chk.Verbose)
	}
}

func TestImagXpowN(tst *testing.T) {

	//verbose()
	chk.PrintTitle("ImagXpowN. calc (x⋅i)ⁿ")

	x := 2.5
	for n := 0; n < 21; n++ {
		res := ImagXpowN(x, n)
		tol := 1e-15
		if n > 2 {
			tol = 1e-12
		}
		if n > 8 {
			tol = 1e-11
		}
		if n > 10 {
			tol = 1e-9
		}
		if n > 14 {
			tol = 1e-6
		}
		chk.AnaNumC(tst, io.Sf("pow(x*i,%2d)", n), tol, res, cmplx.Pow(complex(0, x), complex(float64(n), 0)), chk.Verbose)
	}
}

func TestPowP(tst *testing.T) {

	//verbose()
	chk.PrintTitle("PowP. calc xⁿ with positive integer n")

	res := PowP(0, 0)
	chk.Float64(tst, "0⁰", 1e-17, PowP(0, 0), math.Pow(0, 0))

	var n uint32
	x := 2.5
	for n = 0; n < 23; n++ {
		res = PowP(x, n)
		tol := 1e-17
		chk.AnaNum(tst, io.Sf("pow(x,%2d)", n), tol, res, math.Pow(x, float64(n)), chk.Verbose)
	}
}

func TestPow2(tst *testing.T) {
	//verbose()
	chk.PrintTitle("Pow2. calc x²")
	chk.Float64(tst, "2²", 1e-15, Pow2(2), 4)
	chk.Float64(tst, "3²", 1e-15, Pow2(3), 9)
	chk.Float64(tst, "4²", 1e-15, Pow2(4), 16)
	chk.Float64(tst, "10²", 1e-15, Pow2(10), 100)
}

func TestPow3(tst *testing.T) {
	//verbose()
	chk.PrintTitle("Pow3. calc x³")
	chk.Float64(tst, "2³", 1e-15, Pow3(2), 8)
	chk.Float64(tst, "3³", 1e-15, Pow3(3), 27)
	chk.Float64(tst, "4³", 1e-15, Pow3(4), 64)
	chk.Float64(tst, "10³", 1e-15, Pow3(10), 1000)
}
