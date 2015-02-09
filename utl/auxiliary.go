// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// BestSquare finds the best square for given size=Nrows * Ncolumns
func BestSquare(size int) (nrow, ncol int) {
    nrow = -1 // not found
    ncol = -1 // not found
    for x := 1; x <= size; x++ {
        if (x * x) >= size {
            if (x * x) == size {
                nrow = x
                ncol = x
                return
            } else {
                for y := x; y >= 1; y-- {
                    if (x * y) == size {
                        nrow = x
                        ncol = y
                        return
                    }
                }
            }
        }
    }
    return
}

func IntPrint(a []int, printfkey string) (res string) {
    res = "["
    for i, val := range a {
        if i > 0 {
            res += Sf(", ")
        }
        res += Sf(printfkey, val)
    }
    res += "]"
    return
}

func DblPrint(a []float64, printfkey string) (res string) {
    res = "["
    for i, val := range a {
        if i > 0 {
            res += Sf(", ")
        }
        res += Sf(printfkey, val)
    }
    res += "]"
    return
}
