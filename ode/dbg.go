// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ode

import (
	"fmt"
)

func pv(name string, v []float64) {
	fmt.Printf(name + " =")
	for _, val := range v {
		//fmt.Printf("%23.15e", val)
		fmt.Printf("%20.10e", val)
	}
	fmt.Println()
}
