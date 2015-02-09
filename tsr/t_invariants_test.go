// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

const (
	SAVEPLOT = false
	//SAVEPLOT = true
)

func Test_invs01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs01")

	sig := [][]float64{
		{100, 150, 5},
		{150, 100, 10},
		{5, 10, 100},
	}
	σ := make([]float64, 6)
	s := make([]float64, 6)
	s_ := make([]float64, 6)
	Ten2Man(σ, sig) // σ := sig
	p := M_p(σ)
	q := M_q(σ)
	θ := M_θ(σ)
	sno, p_, q_ := M_devσ(s, σ)
	p1, q1, θ1 := M_pqθ(σ)
	la.MatVecMul(s_, 1, Psd, σ)
	la.PrintMat("sig", sig, "%8g", false)
	utl.Pf("σ   = %v\n", σ)
	utl.Pf("s   = %v\n", s)
	utl.Pf("s_  = %v\n", s_)
	utl.Pf("sno = %v\n", sno)
	utl.Pf("p   = %v\n", p)
	utl.Pf("q   = %v\n", q)
	utl.Pf("q_  = %v\n", q_)
	utl.Pf("θ   = %v\n", θ)
	utl.CheckScalar(tst, "p", 1e-17, p, p_)
	utl.CheckScalar(tst, "p", 1e-17, p, -100)
	utl.CheckScalar(tst, "q", 1e-17, q, 260.52830940226056)
	utl.CheckScalar(tst, "q", 1e-13, q, q_)
	utl.CheckVector(tst, "s", 1e-17, s, s_)
	utl.CheckScalar(tst, "p1", 1e-17, p, p1)
	utl.CheckScalar(tst, "q1", 1e-13, q, q1)
	utl.CheckScalar(tst, "θ1", 1e-17, θ, θ1)
}

func Test_invs02(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs02")

	eps := [][]float64{
		{100 / 200.0, 150 / 200.0, 5 / 200.0},
		{150 / 200.0, 100 / 200.0, 10 / 200.0},
		{5 / 200.0, 10 / 200.0, 100 / 200.0},
	}
	ε := make([]float64, 6)
	e := make([]float64, 6)
	e_ := make([]float64, 6)
	Ten2Man(ε, eps)
	εv := M_εv(ε)
	εd := M_εd(ε)
	eno, εv_, εd_ := M_devε(e, ε)
	la.MatVecMul(e_, 1, Psd, ε)
	la.PrintMat("eps", eps, "%8g", false)
	utl.Pf("ε   = %v\n", ε)
	utl.Pf("e   = %v\n", e)
	utl.Pf("e_  = %v\n", e_)
	utl.Pf("eno = %v\n", eno)
	utl.Pf("εv  = %v\n", εv)
	utl.Pf("εd  = %v\n", εd)
	utl.Pf("εd_ = %v\n", εd_)
	utl.CheckScalar(tst, "εv", 1e-17, εv, εv_)
	utl.CheckScalar(tst, "εv", 1e-17, εv, eps[0][0]+eps[1][1]+eps[2][2])
	utl.CheckScalar(tst, "εd", 1e-13, εd, εd_)
	utl.CheckVector(tst, "e", 1e-17, e, e_)
}

func Test_invs03(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs03")

	// square with vertical stress only under plane-strain
	E, ν := 210000.0, 0.49999
	qY := 240.0
	σx := 0.0
	σy := -qY / math.Sqrt(ν*ν-ν+1.0)
	σz := ν * (σx + σy)
	εx := -(ν*σz + ν*σy - σx) / E
	εy := -(ν*σz - σy + ν*σx) / E
	εz := 0.0

	// check
	c := E / ((1.0 + ν) * (1.0 - 2.0*ν))
	De := [][]float64{
		{c * (1.0 - ν), c * ν, c * ν, 0.0},
		{c * ν, c * (1.0 - ν), c * ν, 0.0},
		{c * ν, c * ν, c * (1.0 - ν), 0.0},
		{0.0, 0.0, 0.0, c * (1.0 - 2.0*ν)},
	}
	ε := [][]float64{
		{εx, 0, 0},
		{0, εy, 0},
		{0, 0, εz},
	}
	εm := make([]float64, 4)
	σm := make([]float64, 4)
	Ten2Man(εm, ε)
	la.MatVecMul(σm, 1, De, εm)
	q := M_q(σm)
	θ := M_θ(σm)
	utl.Pfcyan("σm = %v\n", σm)
	utl.Pfcyan("q  = %v\n", q)
	utl.Pfcyan("θ  = %v\n", θ)
	utl.CheckScalar(tst, "q", 1e-10, q, qY)
	utl.CheckScalar(tst, "θ", 1e-3, θ, 0)
}

func run_invs_tests(tst *testing.T, a []float64, ver bool) {
	at := Alloc2()
	Man2Ten(at, a) // at := TensorVersionOf(a)
	na := M_Norm(a)
	tra := M_Tr(a)
	deva := M_Dev(a)
	deta := M_Det(a)
	w := M_w(a)
	trat := Tr(at)
	detat := Det(at)
	p := M_p(a)
	q := M_q(a)
	θ := M_θ(a)
	s := make([]float64, len(a))
	sX := make([]float64, len(a))
	p_, q_, θ_ := M_pqθ(a)
	pp, qq, rr := M_pqw(a)
	pX, qX, rX := M_pqws(sX, a)
	sno, p1, q1 := M_devσ(s, a)
	λ0, λ1, λ2, err := M_PrincValsNum(a)
	if err != nil {
		utl.Panic("PrincValsNum failed:\n%v", err)
	}
	I1, I2, I3 := M_CharInvs(a)
	utl.CheckScalar(tst, "tr(a)", 1e-17, tra, trat)
	utl.CheckScalar(tst, "det(a)", 1e-14, deta, detat)
	utl.CheckScalar(tst, "p", 1e-14, p, p_)
	utl.CheckScalar(tst, "q", 1e-14, q, q_)
	utl.CheckScalar(tst, "sno", 1e-14, sno, q/SQ3by2)
	utl.CheckScalar(tst, "p1", 1e-14, p, p1)
	utl.CheckScalar(tst, "q1", 1e-14, q, q1)
	utl.CheckScalar(tst, "θ", 1e-14, θ, θ_)
	utl.CheckScalar(tst, "pp", 1e-14, p, pp)
	utl.CheckScalar(tst, "qq", 1e-14, q, qq)
	utl.CheckScalar(tst, "rr", 1e-14, θ, math.Asin(rr)*180.0/(3.0*math.Pi))
	utl.CheckScalar(tst, "pX", 1e-14, p, pX)
	utl.CheckScalar(tst, "qX", 1e-14, q, qX)
	utl.CheckScalar(tst, "rX", 1e-14, θ, math.Asin(rX)*180.0/(3.0*math.Pi))
	utl.CheckScalar(tst, "I1", 1e-17, I1, tra)
	utl.CheckScalar(tst, "I3", 1e-17, I3, deta)
	utl.CheckScalar(tst, "I1", 1e-14, I1, λ0+λ1+λ2)
	utl.CheckScalar(tst, "I2", 1e-12, I2, λ0*λ1+λ1*λ2+λ2*λ0)
	utl.CheckScalar(tst, "I3", 1e-12, I3, λ0*λ1*λ2)
	if ver {
		utl.Pf("θ    = %v\n", θ)
		utl.Pf("na   = %v\n", na)
		utl.Pf("tra  = %v\n", tra)
		utl.Pf("deva = %v\n", deva)
		utl.Pf("deta = %v\n", deta)
		utl.Pf("w    = %v\n", w)
	}
	devat := Alloc2()
	deva_ := Alloc2()
	Man2Ten(devat, deva)
	Add(deva_, 1, at, -(at[0][0]+at[1][1]+at[2][2])/3.0, It) // deva_ := at - tr(at) * It / 3
	utl.CheckMatrix(tst, "deva", 1e-17, devat, deva_)
	utl.CheckVector(tst, "s", 1e-14, s, deva)
	utl.CheckVector(tst, "sX", 1e-14, s, sX)
	// octahedral invariants
	σa, σb, σc := L2O(λ0, λ1, λ2)
	if σa > 0 {
		σa = -σa
	}
	if σb < 0 {
		σb = -σb
	}
	σa_, σb_, σc_ := PQW2O(p, q, w)
	Σa, Σb, Σc := M_oct(a)
	utl.CheckScalar(tst, "σa", 1e-13, σa, σa_)
	utl.CheckScalar(tst, "σb", 1e-13, σb, σb_)
	utl.CheckScalar(tst, "σc", 1e-13, σc, σc_)
	utl.CheckScalar(tst, "Σa", 1e-13, σa, Σa)
	utl.CheckScalar(tst, "Σb", 1e-13, σb, Σb)
	utl.CheckScalar(tst, "Σc", 1e-13, σc, Σc)
	if ver {
		utl.Pforan("λ0 = %v\n", λ0)
		utl.Pforan("λ1 = %v\n", λ1)
		utl.Pforan("λ2 = %v\n", λ2)
		utl.Pforan("σa = %v (%v)\n", σa, σa_)
		utl.Pforan("σb = %v (%v)\n", σb, σb_)
		utl.Pforan("σc = %v (%v)\n", σc, σc_)
	}
}

func Test_invs04(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs04")

	a := []float64{-10.0, -20.0, -30.0, 4.0 * SQ2, 5.0 * SQ2, 6.0 * SQ2}
	at := Alloc2()
	Man2Ten(at, a)
	utl.Pf("a = %v\n", a)
	utl.CheckMatrix(tst, "Man2Ten", 1e-17, at, [][]float64{
		{-10, 4, 6},
		{4, -20, 5},
		{6, 5, -30},
	})

	b := []float64{-88, -77, -55, -3 * SQ2}
	bt := Alloc2()
	Man2Ten(bt, b)
	utl.Pf("b = %v\n", b)
	utl.CheckMatrix(tst, "Man2Ten", 1e-17, bt, [][]float64{
		{-88, -3, 0},
		{-3, -77, 0},
		{0, 0, -55},
	})

	ver := true
	run_invs_tests(tst, a, ver)
	run_invs_tests(tst, b, ver)
}

func Test_invs05(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs05")

	if SAVEPLOT {
		plt.Reset()
		plt.SetForPng(1, 500, 125)
		PlotRosette(1.1, true, true, true, 7)
	}

	addtoplot := func(σa, σb float64, σ []float64) {
		plt.PlotOne(σa, σb, "'ro', ms=5")
		plt.Text(σa, σb, utl.Sf("$\\sigma_{123}=(%g,%g,%g)$", σ[0], σ[1], σ[2]), "size=8")
	}

	dotest := func(σ []float64, σacor, σbcor, σccor, θcor, tolσ float64) {
		w := M_w(σ)
		θ2 := math.Asin(w) * 180.0 / (3.0 * math.Pi)
		θ3 := M_θ(σ)
		σa, σb, σc := L2O(σ[0], σ[1], σ[2])
		σ0, σ1, σ2 := O2L(σa, σb, σc)
		σI, σA := make([]float64, 3), []float64{σa, σb, σc}
		la.MatVecMul(σI, 1, O2Lmat(), σA) // σI := L * σA
		utl.Pf("σa σb σc = %v %v %v\n", σa, σb, σc)
		utl.Pf("w        = %v\n", w)
		utl.Pf("θ2, θ3   = %v, %v\n", θ2, θ3)
		utl.CheckScalar(tst, "σa", 1e-17, σa, σacor)
		utl.CheckScalar(tst, "σb", 1e-17, σb, σbcor)
		utl.CheckScalar(tst, "σc", 1e-17, σc, σccor)
		utl.CheckScalar(tst, "σ0", tolσ, σ0, σ[0])
		utl.CheckScalar(tst, "σ1", tolσ, σ1, σ[1])
		utl.CheckScalar(tst, "σ2", tolσ, σ2, σ[2])
		utl.CheckScalar(tst, "σI0", tolσ, σI[0], σ[0])
		utl.CheckScalar(tst, "σI1", tolσ, σI[1], σ[1])
		utl.CheckScalar(tst, "σI2", tolσ, σI[2], σ[2])
		utl.CheckScalar(tst, "θ2", 1e-6, θ2, θcor)
		utl.CheckScalar(tst, "θ3", 1e-17, θ3, θ2)
		addtoplot(σa, σb, σ)
	}

	dotest([]float64{-1, 0, 0, 0}, 0, 2.0/SQ6, 1.0/SQ3, 30, 1e-15)
	dotest([]float64{0, -1, 0, 0}, 1.0/SQ2, -1.0/SQ6, 1.0/SQ3, 30, 1e-15)
	dotest([]float64{0, 0, -1, 0}, -1.0/SQ2, -1.0/SQ6, 1.0/SQ3, 30, 1e-15)

	if SAVEPLOT {
		plt.Gll("$\\sigma_a$", "$\\sigma_b$", "")
		plt.Equal()
		plt.SaveD("/tmp/gosl", "fig_invs05.png")
	}
}

func Test_invs06(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs06")

	ver := true
	s := make([]float64, 4)
	cos30, sin30 := SQ3/2.0, 0.5
	for i := 0; i < 11; i++ {
		σa, σb, σc := -sin30*float64(i), cos30*float64(i), 1.0/SQ3
		σ0, σ1, σ2 := O2L(σa, σb, σc)
		σ := []float64{σ0, σ1, σ2, 0}
		θ := M_θ(σ)
		utl.Pf("σ = %v\n", σ)
		utl.Pf("θ = %v\n", θ)
		run_invs_tests(tst, σ, ver)
		sno, p, _ := M_devσ(s, σ)
		utl.CheckScalar(tst, "σc", 1e-15, p, σc/SQ3)
		utl.CheckScalar(tst, "sno", 1e-15, sno, float64(i))
	}
}

func Test_invs07(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mSome error has happened:[0m\n", err)
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("invs07")

	a, b, β, ϵ := -1.0, 0.0, 1.0, 1e-3

	σ := []float64{-1, -1, -1, 0}
	λ := []float64{1, 1, 1, 0}

	pcam, qcam, _ := M_pqw(σ)
	poct, qoct := pcam*SQ3, qcam*SQ2by3

	nold := make([]float64, 3)
	SmpUnitDirector(nold, λ, b)
	psmp1, qsmp1, err := GenInvs(λ, nold, 1)
	if err != nil {
		utl.Panic("M_GenInvs failed:\n%v", err)
	}

	psmp2, qsmp2, err := M_pq_smp(σ, a, b, β, ϵ)
	if err != nil {
		utl.Panic("M_pq_smp failed:\n%v", err)
	}
	utl.Pforan("pcam,  qcam  = %v, %v\n", pcam, qcam)
	utl.Pforan("poct,  qoct  = %v, %v\n", poct, qoct)
	utl.Pforan("psmp1, qsmp1 = %v, %v\n", psmp1, qsmp1)
	utl.Pforan("psmp2, qsmp2 = %v, %v\n", psmp2, qsmp2)
	utl.CheckScalar(tst, "p", 1e-15, psmp1, psmp2)
	utl.CheckScalar(tst, "q", 1e-15, qsmp1, qsmp2)
}
