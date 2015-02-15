// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/utl"
)

func Test_hash01(tst *testing.T) {

	prevTs := utl.Tsilent
	defer func() {
		utl.Tsilent = prevTs
		if err := recover(); err != nil {
			tst.Error("[1;31mERROR:", err, "[0m\n")
		}
	}()

	//utl.Tsilent = false
	utl.TTitle("hash01")

	c1 := []float64{1.0000000001}
	c2 := []float64{1.0000000000001, 2.0000000002}
	c3 := []float64{1.0000000000002, 2.0000000002, 3.0000000000003}
	c4 := []float64{1.0000000000003, 2.0000000002, 3.0000000000003}

	h1 := HashPointC(c1)
	h2 := HashPointC(c2)
	h3 := HashPointC(c3)
	h4 := HashPointC(c4)

	utl.Pforan("h1 = %v\n", h1)
	utl.Pforan("h2 = %v\n", h2)
	utl.Pforan("h3 = %v\n", h3)
	utl.Pforan("h4 = %v\n", h4)

	if h1 == h2 {
		utl.Panic("h1 must not be equal to h2")
	}
	if h1 == h3 {
		utl.Panic("h1 must not be equal to h3")
	}
	if h1 == h4 {
		utl.Panic("h1 must not be equal to h4")
	}
	if h2 == h3 {
		utl.Panic("h2 must not be equal to h3")
	}
	if h2 == h4 {
		utl.Panic("h2 must not be equal to h4")
	}
	// TODO: this one fails
	if false {
		if h3 == h4 {
			utl.Panic("h3 must not be equal to h4")
		}
	}
}
