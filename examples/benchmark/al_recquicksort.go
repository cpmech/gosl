// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"sort"
	"time"

	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/rnd"
	"github.com/cpmech/gosl/utl"
	"github.com/cpmech/gosl/utl/al"
)

func run(N []float64, nRnd int, sorter func(A []float64, compare func(a, b float64) int)) (Tfwd, Tbwd, Trnd []float64) {
	Tfwd = make([]float64, len(N))
	Tbwd = make([]float64, len(N))
	Trnd = make([]float64, len(N))
	for i, nn := range N {
		n := int(nn)

		Afwd := utl.LinSpace(0, nn-1, n)
		t0 := time.Now()
		sorter(Afwd, al.Float64Comparator)
		Tfwd[i] = time.Now().Sub(t0).Seconds()

		Abwd := utl.LinSpace(nn-1, 0, n)
		t0 = time.Now()
		sorter(Abwd, al.Float64Comparator)
		Tbwd[i] = time.Now().Sub(t0).Seconds()

		sum := 0.0
		for j := 0; j < nRnd; j++ {
			Arnd := utl.GetCopy(Afwd)
			rnd.Shuffle(Arnd)
			t0 = time.Now()
			sorter(Arnd, al.Float64Comparator)
			del := float64(time.Now().Sub(t0).Nanoseconds())
			sum += del
		}
		Trnd[i] = (sum / 1e9) / float64(nRnd)
	}
	return
}

func main() {
	rnd.Init(0)
	plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150, Prop: 1.7})

	nRnd := 10
	N := utl.LinSpace(2, 1e3, 11)

	TfwdRC, TbwdRC, TrndRC := run(N, nRnd, al.Float64RecQuickSort)
	TfwdNO, TbwdNO, TrndNO := run(N, nRnd, al.Float64RecQuickSortNonOpt)
	TfwdQS, TbwdQS, TrndQS := run(N, nRnd, func(A []float64, dummy func(a, b float64) int) { utl.Qsort(A) })
	TfwdGO, TbwdGO, TrndGO := run(N, nRnd, func(A []float64, dummy func(a, b float64) int) { sort.Float64s(A) })

	plt.Subplot(2, 1, 1)

	plt.Plot(N, TfwdRC, &plt.A{C: plt.C(0, 0), L: "Recursive: fwd", M: plt.M(0, 0), Ls: "-", NoClip: true})
	plt.Plot(N, TbwdRC, &plt.A{C: plt.C(0, 0), L: "Recursive: bwd", M: plt.M(1, 0), Ls: "-", NoClip: true})
	plt.Plot(N, TrndRC, &plt.A{C: plt.C(0, 0), L: "Recursive: rnd", M: plt.M(2, 0), Ls: "-", NoClip: true})

	plt.Plot(N, TfwdNO, &plt.A{C: plt.C(1, 0), L: "RecNonOpt: fwd", M: plt.M(0, 0), Ls: "-", NoClip: true})
	plt.Plot(N, TbwdNO, &plt.A{C: plt.C(1, 0), L: "RecNonOpt: bwd", M: plt.M(1, 0), Ls: "-", NoClip: true})
	plt.Plot(N, TrndNO, &plt.A{C: plt.C(1, 0), L: "RecNonOpt: rnd", M: plt.M(2, 0), Ls: "-", NoClip: true})

	plt.Plot(N, TfwdGO, &plt.A{C: plt.C(3, 0), L: "Go native: fwd", M: plt.M(0, 0), Ls: ":", NoClip: true})
	plt.Plot(N, TbwdGO, &plt.A{C: plt.C(3, 0), L: "Go native: bwd", M: plt.M(1, 0), Ls: ":", NoClip: true})
	plt.Plot(N, TrndGO, &plt.A{C: plt.C(3, 0), L: "Go native: rnd", M: plt.M(2, 0), Ls: ":", NoClip: true})

	plt.Plot(N, TfwdQS, &plt.A{C: plt.C(2, 0), L: "utl.Qsort: fwd", M: plt.M(0, 0), Ls: "--", NoClip: true})
	plt.Plot(N, TbwdQS, &plt.A{C: plt.C(2, 0), L: "utl.Qsort: bwd", M: plt.M(1, 0), Ls: "--", NoClip: true})
	plt.Plot(N, TrndQS, &plt.A{C: plt.C(2, 0), L: "utl.Qsort: rnd", M: plt.M(2, 0), Ls: "--", NoClip: true})

	plt.Gll("$N$", "$time\\;[s]$", nil)
	plt.HideTRborders()

	plt.Subplot(2, 1, 2)

	plt.Plot(N, TrndRC, &plt.A{C: plt.C(0, 0), L: "Recursive: rnd", M: plt.M(2, 0), Ls: "-", NoClip: true})

	plt.Plot(N, TfwdGO, &plt.A{C: plt.C(3, 0), L: "Go native: fwd", M: plt.M(0, 0), Ls: ":", NoClip: true})
	plt.Plot(N, TbwdGO, &plt.A{C: plt.C(3, 0), L: "Go native: bwd", M: plt.M(1, 0), Ls: ":", NoClip: true})
	plt.Plot(N, TrndGO, &plt.A{C: plt.C(3, 0), L: "Go native: rnd", M: plt.M(2, 0), Ls: ":", NoClip: true})

	plt.Plot(N, TfwdQS, &plt.A{C: plt.C(2, 0), L: "utl.Qsort: fwd", M: plt.M(0, 0), Ls: "--", NoClip: true})
	plt.Plot(N, TbwdQS, &plt.A{C: plt.C(2, 0), L: "utl.Qsort: bwd", M: plt.M(1, 0), Ls: "--", NoClip: true})
	plt.Plot(N, TrndQS, &plt.A{C: plt.C(2, 0), L: "utl.Qsort: rnd", M: plt.M(2, 0), Ls: "--", NoClip: true})

	plt.Gll("$N$", "$time\\;[s]$", nil)
	plt.HideTRborders()
	plt.Save("/tmp/gosl", "al_recquicksort")
}
