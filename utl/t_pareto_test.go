// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_pareto01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pareto01. compare vectors: Pareto-optimal")

	u := []float64{1, 2, 3, 4, 5, 6}
	v := []float64{1, 2, 3, 4, 5, 6}
	io.Pforan("u = %v\n", u)
	io.Pfblue2("v = %v\n", v)
	u_dominates, v_dominates := DblsParetoMin(u, v)
	io.Pfpink("u_dominates = %v\n", u_dominates)
	io.Pfpink("v_dominates = %v\n", v_dominates)
	if u_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if v_dominates {
		tst.Errorf("test failed\n")
		return
	}

	v = []float64{1, 1.8, 3, 4, 5, 6}
	io.Pforan("\nu = %v\n", u)
	io.Pfblue2("v = %v\n", v)
	u_dominates, v_dominates = DblsParetoMin(u, v)
	io.Pfpink("u_dominates = %v\n", u_dominates)
	io.Pfpink("v_dominates = %v\n", v_dominates)
	if u_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if !v_dominates {
		tst.Errorf("test failed\n")
		return
	}

	v = []float64{1, 2.1, 3, 4, 5, 6}
	io.Pforan("\nu = %v\n", u)
	io.Pfblue2("v = %v\n", v)
	u_dominates, v_dominates = DblsParetoMin(u, v)
	io.Pfpink("u_dominates = %v\n", u_dominates)
	io.Pfpink("v_dominates = %v\n", v_dominates)
	if !u_dominates {
		tst.Errorf("test failed\n")
		return
	}
	if v_dominates {
		tst.Errorf("test failed\n")
		return
	}
}

func Test_pareto02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pareto02. probabilistic Pareto-optimal")

	rand.Seed(time.Now().UnixNano())

	Φ := []float64{0, 0.25, 0.5, 0.75, 1}
	u := 0.4
	v := 0.6
	for _, φ := range Φ {
		p := ProbContestSmall(u, v, φ)
		io.Pf("u=%v v=%v p(u,v,%5.2f) = %.8f\n", u, v, φ, p)
	}

	ntrials := 1000
	doplot := chk.Verbose
	var buf bytes.Buffer
	var Zu []float64

	U := []float64{1, 2, 3, 4, 5, 6}
	V := []float64{1, 2, 3, 4, 5, 6}
	zu, _ := run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{1, 2, 3, 4, 5, 6}
	V = []float64{2, 2, 3, 4, 5, 6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{1, 2, 3, 4, 5, 6}
	V = []float64{2, 3, 3, 4, 5, 6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{1, 2, 3, 4, 5, 6}
	V = []float64{2, 3, 4, 4, 5, 6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{1, 2, 3, 4, 5, 6}
	V = []float64{2, 3, 4, 5, 5, 6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{1, 2, 3, 4, 5, 6}
	V = []float64{2, 3, 4, 5, 6, 6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{1, 2, 3, 4, 5, 6}
	V = []float64{2, 3, 4, 5, 6, 7}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	if doplot {
		plot_pareto_test(&buf, "test_pareto02", Φ, Zu, false, false)
	}
}

func Test_pareto03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("pareto03. probabilistic Pareto-optimal")

	rand.Seed(time.Now().UnixNano())

	Φ := []float64{0, 0.25, 0.5, 0.75, 1}
	u := 0.4
	v := 0.6
	for _, φ := range Φ {
		p := ProbContestSmall(u, v, φ)
		io.Pf("u=%v v=%v p(u,v,%5.2f) = %.8f\n", u, v, φ, p)
	}

	ntrials := 1000
	doplot := chk.Verbose
	var buf bytes.Buffer
	var Zu []float64

	U := []float64{-1, -2, -3, -4, -5, -6}
	V := []float64{-1, -2, -3, -4, -5, -6}
	zu, _ := run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{-2, -2, -3, -4, -5, -6}
	V = []float64{-1, -2, -3, -4, -5, -6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{-2, -3, -3, -4, -5, -6}
	V = []float64{-1, -2, -3, -4, -5, -6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{-2, -3, -4, -4, -5, -6}
	V = []float64{-1, -2, -3, -4, -5, -6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{-2, -3, -4, -5, -5, -6}
	V = []float64{-1, -2, -3, -4, -5, -6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{-2, -3, -4, -5, -6, -6}
	V = []float64{-1, -2, -3, -4, -5, -6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	U = []float64{-2, -3, -4, -5, -6, -7}
	V = []float64{-1, -2, -3, -4, -5, -6}
	zu, _ = run_pareto_test(U, V, Φ, ntrials)
	Zu = append(Zu, zu...)

	if doplot {
		plot_pareto_test(&buf, "test_pareto03", Φ, Zu, false, true)
	}
}

func run_pareto_test(U, V, Φ []float64, ntrials int) (zu, zv []float64) {
	io.Pforan("\nu = %v\n", U)
	io.Pforan("v = %v\n", V)
	io.PfWhite("%5s%13s%13s\n", "φ", "u wins", "v wins")
	zu = make([]float64, len(Φ))
	zv = make([]float64, len(Φ))
	for i, φ := range Φ {
		u_wins := 0
		v_wins := 0
		for j := 0; j < ntrials; j++ {
			u_dominates := DblsParetoMinProb(U, V, φ)
			if u_dominates {
				u_wins++
			} else {
				v_wins++
			}
		}
		zu[i] = 100 * float64(u_wins) / float64(ntrials)
		zv[i] = 100 * float64(v_wins) / float64(ntrials)
		io.Pf("%5.2f%12.2f%%%12.2f%%\n", φ, zu[i], zv[i])
	}
	return
}

func write_python_array(buf *bytes.Buffer, name string, x []float64) {
	io.Ff(buf, "%s=np.array([", name)
	for i := 0; i < len(x); i++ {
		io.Ff(buf, "%g,", x[i])
	}
	io.Ff(buf, "])\n")
}

func plot_pareto_test(buf *bytes.Buffer, fnkey string, Φ, Zu []float64, show, negative bool) {
	io.Ff(buf, "from gosl import SetForEps, Save\n")
	io.Ff(buf, "from mpl_toolkits.mplot3d import Axes3D\n")
	io.Ff(buf, "import matplotlib.pyplot as plt\n")
	io.Ff(buf, "import numpy as np\n")
	io.Ff(buf, "fig = plt.figure()\n")
	io.Ff(buf, "ax = fig.add_subplot(111, projection='3d')\n")
	write_python_array(buf, "phi", Φ)
	write_python_array(buf, "dz", Zu)
	io.Ff(buf, "SetForEps(0.75, 455, mplclose=0, text_usetex=0)\n")
	io.Ff(buf, "n=len(phi)\nx,y=np.meshgrid(phi,np.linspace(0,1,7))\nx=x.flatten()\ny=y.flatten()\nz=np.zeros(n*7)\n")
	io.Ff(buf, "ax.bar3d(x,y,z,dx=0.1*np.ones(n*7),dy=0.1*np.ones(n*7),dz=dz, color='#cee9ff')\n")
	io.Ff(buf, "ax.set_xlabel('$\\phi$')\n")
	io.Ff(buf, "ax.set_zlabel('u-wins [%%]')\n")
	io.Ff(buf, "ax.set_xticks([0,0.25,0.5,0.75,1])\n")
	io.Ff(buf, "ax.set_xticklabels(['0.0','0.25','0.5','0.75','1.0'])\n")
	if negative {
		io.Ff(buf, "ax.set_yticklabels(['u=[-1 -2 -3 -4 -5 -6]', 'u=[-2 -2 -3 -4 -5 -6]', 'u=[-2 -3 -3 -4 -5 -6]', 'u=[-2 -3 -4 -4 -5 -6]', 'u=[-2 -3 -4 -5 -5 -6]', 'u=[-2 -3 -4 -5 -6 -6]', 'u=[-2 -3 -4 -5 -6 -7]'], rotation=-15,verticalalignment='baseline',horizontalalignment='left')\n")
	} else {
		io.Ff(buf, "ax.set_yticklabels(['v=[1 2 3 4 5 6]', 'v=[2 2 3 4 5 6]', 'v=[2 3 3 4 5 6]', 'v=[2 3 4 4 5 6]', 'v=[2 3 4 5 5 6]', 'v=[2 3 4 5 6 6]', 'v=[2 3 4 5 6 7]'], rotation=-15,verticalalignment='baseline',horizontalalignment='left')\n")
	}
	io.Ff(buf, "import matplotlib.patheffects as path_effects\n")
	io.Ff(buf, "for i, xval in enumerate(x): ax.text(xval,y[i],dz[i],'%%.2f'%%dz[i],color='#bf0000',fontsize=10, path_effects=[path_effects.withSimplePatchShadow(offset=(1,-1),shadow_rgbFace='white')])\n")
	if show {
		io.Ff(buf, "plt.show()\n")
	} else {
		io.Ff(buf, "plt.savefig('/tmp/gosl/%s.eps')\n", fnkey)
		io.Ff(buf, "print 'file </tmp/gosl/%s.eps> witten'\n", fnkey)
	}
	io.WriteFileVD("/tmp/gosl", io.Sf("%s.py", fnkey), buf)
}
