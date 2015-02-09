// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"strings"
)

// XYreader holds data for reading X and Y values from a file
type XYreader struct {
	X []float64   // X   [n]: n=number of data rows
	Y [][]float64 // Y[m][n]: m=number of y columns, n=number of data rows
}

// ReadFile reads X and Y values from a file
func (o *XYreader) ReadFile(fn string) (oserr error) {
	f, oserr := OpenFileR(fn)
	if oserr != nil {
		return
	}
	header := true
	ReadLinesFile(f, func(idx int, line string) (stop bool) {
		r := strings.Fields(line)
		if len(r) == 0 { // skip empty lines
			return
		}
		ny := len(r) - 1
		if ny < 2 {
			oserr = Err("ReadFile: number of columns must be at least 2 ('x' and 'y0')")
			return true
		}
		if header {
			if r[0] != "x" {
				oserr = Err("ReadFile: the first column on the header of file <%s> must 'x'", f.Name())
				return true
			}
			for i := 0; i < ny; i++ {
				if r[1+i] != Sf("y%d", i) {
					oserr = Err("ReadFile: the column #%d in the header of file <%s> must be 'y%d' and not '%s'", i+2, f.Name(), i, r[1+i])
					return true
				}
			}
			o.X = make([]float64, 0)
			o.Y = make([][]float64, ny)
			for i := 0; i < ny; i++ {
				o.Y[i] = make([]float64, 0)
			}
			header = false
		} else {
			o.X = append(o.X, Atof(r[0]))
			for i := 0; i < len(o.Y); i++ {
				o.Y[i] = append(o.Y[i], Atof(r[1+i]))
			}
		}
		return
	})
	return
}
