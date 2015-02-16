// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/utl"
)

// HashPoint returns a unique id of a point
func HashPoint(x, y, z float64) int {
	return int(x*100001 + y*10000001 + z*1000000001)
}

// HashPoint returns a unique id of a point (given slice of coordinates)
func HashPointC(c []float64) int {
	var x, y, z float64
	n := len(c)
	if n > 0 {
		x = c[0]
	}
	if n > 1 {
		y = c[1]
	}
	if n > 2 {
		z = c[2]
	}
	return HashPoint(x, y, z)
}

type BinEntry struct {
	Id int       // object Id
	x  []float64 // entry coordinate
}

type Bin struct {
	Idx     int         // index of bin
	Entries []*BinEntry // entries
}

type Bins struct {
	Ndim int       // space dimension
	Xi   []float64 // [ndim] left/lower-most point
	Xf   []float64 // [ndim] right/upper-most point
	L    []float64 // [ndim] whole box lengths
	N    []int     // [ndim] number of divisions
	S    float64   // size of bins
	All  []*Bin    // [nbins] all bins
	tmp  []int     // [ndim] temporary (auxiliary) slice
}

// xi   -- [ndim] initial positions
// l    -- [ndim] whole box lengths
// ndiv -- number of divisions for the maximun length
func (o *Bins) Init(xi, xf []float64, ndiv int) (err error) {
	o.Ndim = len(xi)
	o.Xi = xi
	o.Xf = xf
	if len(xi) != len(xf) || len(xi) < 2 || len(xi) > 3 {
		return utl.Err("sizes of xi and l must be the same and equal to either 2 or 3")
	}

	o.L = make([]float64, o.Ndim)
	o.N = make([]int, o.Ndim)

	for k := 0; k < o.Ndim; k++ {
		o.L[k] = o.Xf[k] - o.Xi[k]
	}

	// maximun lenght
	lmax := math.Max(o.L[0], o.L[1])
	if o.Ndim == 3 {
		lmax = math.Max(lmax, o.L[2])
	}
	o.S = lmax / float64(ndiv)

	// number of divisions
	nbins := 1
	for k := 0; k < o.Ndim; k++ {
		o.N[k] = int(o.L[k]/o.S) + 1
		nbins *= o.N[k]
	}

	o.All = make([]*Bin, nbins)
	o.tmp = make([]int, o.Ndim)

	return
}

func (o *Bins) Append(x []float64, id int) (err error) {
	idx := o.CalcIdx(x)
	if idx < 0 {
		return utl.Err("point %v is out of range", x)
	}
	bin := o.FindBinByIndex(idx)
	if bin == nil {
		return utl.Err("bin index %v is out of range", idx)
	}

	entry := BinEntry{id, x}
	bin.Entries = append(bin.Entries, &entry)
	return
}

func (o *Bins) Clear() {
	o.All = make([]*Bin, 0)
}

func (o Bins) FindBinByIndex(idx int) *Bin {
	if idx < 0 || idx >= len(o.All) {
		return nil
	}
	return o.All[idx]
}

func (o Bins) Find(x []float64) *BinEntry {
	idx := o.CalcIdx(x)
	if idx < 0 {
		return nil // out-of-range
	}
	bin := o.FindBinByIndex(idx)
	dmin := math.MaxFloat64
	var entry, closest *BinEntry
	for _, entry = range bin.Entries {
		var d float64
		for k := 0; k < o.Ndim; k++ {
			d += math.Pow(x[k]-entry.x[k], 2)
		}
		if d < dmin {
			dmin = d
			closest = entry
		}
	}
	return closest
}

// CalcIdx calculates the bin index where the point x is
// returns -1 if out-of-range
func (o Bins) CalcIdx(x []float64) int {
	for k := 0; k < o.Ndim; k++ {
		if x[k] < o.Xi[k] || x[k] > o.Xf[k] {
			return -1
		}
		o.tmp[k] = int((x[k] - o.Xi[k]) / o.S)
	}
	idx := o.tmp[0] + o.tmp[1]*o.N[0]
	if o.Ndim > 2 {
		idx += o.tmp[2] * o.N[0] * o.N[1]
	}
	return idx
}
