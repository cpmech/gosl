// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
    "bytes"
)

// GenIntArray generates the NumPy text in 'b' corresponding to an array of integers
func GenIntArray(b *bytes.Buffer, name string, u []int) {
    Ff(b, "%s=array([", name)
    for i, _ := range u {
        Ff(b, "%d,", u[i])
    }
    Ff(b, "],dtype=int)\n")
}

// GenArray generates the NumPy text in 'b' corresponding to an array of float point numbers
func GenArray(b *bytes.Buffer, name string, u []float64) {
    Ff(b, "%s=array([", name)
    for i, _ := range u {
        Ff(b, "%g,", u[i])
    }
    Ff(b, "],dtype=float)\n")
}

func GenMat(buf *bytes.Buffer, name string, a [][]float64) {
    Ff(buf, "%s=array([", name)
    for i, _ := range a {
        Ff(buf, "[")
        for j, _ := range a[i] {
            Ff(buf, "%g,", a[i][j])
        }
        Ff(buf, "],")
    }
    Ff(buf, "],dtype=float)\n")
}

// Gen2Arrays generates the NumPy text in 'b' corresponding to 2 arrays of float point numbers
func Gen2Arrays(buf *bytes.Buffer, nameA, nameB string, a, b []float64) {
    GenArray(buf, nameA, a)
    GenArray(buf, nameB, b)
}

// Gen3Arrays generates the NumPy text in 'b' corresponding to 3 arrays of float point numbers
func Gen3Arrays(buf *bytes.Buffer, nameA, nameB, nameC string, a, b, c []float64) {
    GenArray(buf, nameA, a)
    Gen2Arrays(buf, nameB, nameC, b, c)
}

// Gen4Arrays generates the NumPy text in 'b' corresponding to 4 arrays of float point numbers
func Gen4Arrays(buf *bytes.Buffer, nameA, nameB, nameC, nameD string, a, b, c, d []float64) {
    GenArray(buf, nameA, a)
    Gen3Arrays(buf, nameB, nameC, nameD, b, c, d)
}

// Gen5Arrays generates the NumPy text in 'b' corresponding to 5 arrays of float point numbers
func Gen5Arrays(buf *bytes.Buffer, nameA, nameB, nameC, nameD, nameE string, a, b, c, d, e []float64) {
    GenArray(buf, nameA, a)
    Gen4Arrays(buf, nameB, nameC, nameD, nameE, b, c, d, e)
}

// Gen6Arrays generates the NumPy text in 'b' corresponding to 6 arrays of float point numbers
func Gen6Arrays(buf *bytes.Buffer, nameA, nameB, nameC, nameD, nameE, nameF string, a, b, c, d, e, f []float64) {
    GenArray(buf, nameA, a)
    Gen5Arrays(buf, nameB, nameC, nameD, nameE, nameF, b, c, d, e, f)
}

// Gen7Arrays generates the NumPy text in 'b' corresponding to 7 arrays of float point numbers
func Gen7Arrays(buf *bytes.Buffer, nameA, nameB, nameC, nameD, nameE, nameF, nameG string, a, b, c, d, e, f, g []float64) {
    GenArray(buf, nameA, a)
    Gen6Arrays(buf, nameB, nameC, nameD, nameE, nameF, nameG, b, c, d, e, f, g)
}

// Gen8Arrays generates the NumPy text in 'b' corresponding to 8 arrays of float point numbers
func Gen8Arrays(buf *bytes.Buffer, nameA, nameB, nameC, nameD, nameE, nameF, nameG, nameH string, a, b, c, d, e, f, g, h []float64) {
    GenArray(buf, nameA, a)
    Gen7Arrays(buf, nameB, nameC, nameD, nameE, nameF, nameG, nameH, b, c, d, e, f, g, h)
}

// GenFuncT generates the NumPy text in 'b' corresponding to a time series by calling the 'F' function
func GenFuncT(b *bytes.Buffer, nameT, nameF string, t0, tf, Δt float64, F func(t float64) float64) {
    // t
    Ff(b, "%s=array([", nameT)
    t := t0
    for t < tf {
        Ff(b, "%g,", t)
        t += Δt
    }
    Ff(b, "],dtype=float)\n")
    // F(t)
    Ff(b, "%s=array([", nameF)
    t = t0
    for t < tf {
        Ff(b, "%g,", F(t))
        t += Δt
    }
    Ff(b, "],dtype=float)\n")
}
