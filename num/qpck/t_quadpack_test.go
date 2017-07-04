// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func TestQags01a(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Qags01a.")

	y := func(x float64) (res float64) {
		return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
	}

	fid := 0
	A, abserr, neval, last, err := Qagse(fid, y, 0, 1, 0, 0, nil, nil, nil, nil, nil)
	if err != nil {
		tst.Errorf("%v\n", err)
		return
	}
	io.Pforan("A      = %v\n", A)
	io.Pforan("abserr = %v\n", abserr)
	io.Pforan("neval  = %v\n", neval)
	io.Pforan("last   = %v\n", last)
	chk.Scalar(tst, "A", 1e-12, A, 1.08268158558)
}

func TestQags01b(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Qags01b. goroutines")

	y := func(x float64) (res float64) {
		return math.Sqrt(1.0 + math.Pow(math.Sin(x), 3.0))
	}

	// channels
	nch := 3
	done := make(chan int, nch)

	// run all
	for ich := 0; ich < nch; ich++ {
		go func(fid int) {
			A, _, _, _, _ := Qagse(fid, y, 0, 1, 0, 0, nil, nil, nil, nil, nil)
			chk.Scalar(tst, "A", 1e-12, A, 1.08268158558)
			done <- 1
		}(ich)
	}

	// wait
	for i := 0; i < nch; i++ {
		<-done
	}
}
