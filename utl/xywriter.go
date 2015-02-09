// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"bytes"
	"math"
)

// Ycorrect is a callback function used to calculate correct Y values
type Ycorrect func(y []float64, x float64)

// XYwriter holds data to write a file with X and Y values
type XYwriter struct {
	k  int
	b  bytes.Buffer
	yc []float64
}

// Write writes a file with X and Y values
func (w *XYwriter) Write(x float64, y []float64, Yc Ycorrect) {

	// write header
	if w.k == 0 {
		Ff(&w.b, "%23s ", "x")
		for i := 0; i < len(y); i++ {
			Ff(&w.b, "%23s ", Sf("y%d", i))
		}
		if Yc != nil {
			Ff(&w.b, "%23s", "max_err_y")
		}
		Ff(&w.b, "\n")
	}

	// write results
	Ff(&w.b, "%23.15E ", x)
	var max_err_y float64
	if Yc != nil {
		if len(w.yc) != len(y) {
			w.yc = make([]float64, len(y))
		}
		Yc(w.yc, x)
	}
	for i := 0; i < len(y); i++ {
		if Yc != nil {
			err := math.Abs(y[i] - w.yc[i])
			if err > max_err_y {
				max_err_y = err
			}
		}
		Ff(&w.b, "%23.15E ", y[i])
	}

	// write error
	if Yc != nil {
		Ff(&w.b, "%23.15E\n", max_err_y)
	} else {
		Ff(&w.b, "\n")
	}
	w.k++
}

// Save save a file with X and Y values
func (w *XYwriter) Save(fn string) {
	WriteFile(fn, &w.b)
}
