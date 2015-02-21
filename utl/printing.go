// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import "github.com/cpmech/gosl/io"

// Deep3Print prints an array of array of array
func Deep3Print(name string, A [][][]float64) {
	io.Pf("%s = [\n", name)
	for _, a := range A {
		io.Pf("  %v\n", a)
	}
	io.Pf("]\n")
}
