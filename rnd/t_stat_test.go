// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

func Test_stat01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("stat01")

	x := []float64{100, 100, 102, 98, 77, 99, 70, 105, 98}

	sum, mean, adev, sdev, vari, skew, kurt, err := Moments(x)
	if err != nil {
		chk.Panic("Moments failed:\n%v", err)
	}
	io.Pforan("x    = %v\n", x)
	io.Pforan("sum  = %v\n", sum)
	io.Pforan("mean = %v\n", mean)
	io.Pforan("adev = %v\n", adev)
	io.Pforan("sdev = %v\n", sdev)
	io.Pforan("vari = %v\n", vari)
	io.Pforan("skew = %v\n", skew)
	io.Pforan("kurt = %v\n", kurt)
	chk.Scalar(tst, "sum ", 1e-17, sum, 849)
	chk.Scalar(tst, "mean", 1e-17, mean, 849.0/9.0)
	chk.Scalar(tst, "sdev", 1e-17, sdev, 12.134661099511597)
	chk.Scalar(tst, "vari", 1e-17, vari, 147.25)

	// TODO: add checks for adev, skew and kurt
}
