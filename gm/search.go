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
// xf   -- [ndim] final positions
// ndiv -- number of divisions for the maximun length
func (o *Bins) Init(xi, xf []float64, ndiv int) (err error) {

	// check for out-of-range values
	o.Ndim = len(xi)
	o.Xi = xi
	o.Xf = xf
	if len(xi) != len(xf) || len(xi) < 2 || len(xi) > 3 {
		return utl.Err("sizes of xi and l must be the same and equal to either 2 or 3")
	}

	// allocate lentgth and number of division slices
	o.L = make([]float64, o.Ndim)
	o.N = make([]int, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		o.L[k] = o.Xf[k] - o.Xi[k]
	}

	// maximun length
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

	// allocate slices
	o.All = make([]*Bin, nbins)
	o.tmp = make([]int, o.Ndim)
	return
}

// Append adds a new entry {x, id} to the bins structure
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

// Clear clears all bins
func (o *Bins) Clear() {
	o.All = make([]*Bin, 0)
}

// Find returns the stored id of the entry whose coordinates are closest to x
// returns -1 if out of range or not found
func (o Bins) Find(x []float64) int {

	// index and check
	idx := o.CalcIdx(x)
	if idx < 0 {
		return -1 // out-of-range
	}

	// search for the closest point
	bin := o.FindBinByIndex(idx)
	dmin := math.MaxFloat64
	id_closest := -1
	var entry *BinEntry
	for _, entry = range bin.Entries {
		var d float64
		for k := 0; k < o.Ndim; k++ {
			d += math.Pow(x[k]-entry.x[k], 2)
		}
		if d < dmin {
			dmin = d
			id_closest = entry.Id
		}
	}
	return id_closest
}

func (o Bins) FindBinByIndex(idx int) *Bin {

	// check
	if idx < 0 || idx >= len(o.All) {
		return nil
	}

	// allocate new bin if necessary
	if o.All[idx] == nil {
		o.All[idx] = new(Bin)
		o.All[idx].Idx = idx
	}
	return o.All[idx]
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

// FindAlongLine gets the ids of entries that lie close to a line
func (o Bins) FindAlongLine(xi, xf []float64, tol float64) []int {

	// essential data
	btol := 0.9 * o.S
	var sbins []*Bin // selected bins
	var p, pi, pf Point
	var x, y, z float64
	pi.X = xi[0]
	pf.X = xf[0]
	pi.Y = xi[1]
	pf.Y = xf[1]
	if o.Ndim == 3 {
		pi.Z = xi[2]
		pf.Z = xf[2]
	}

	// loop along all bins
	for i, bin := range o.All {
		if bin == nil {
			continue
		}
		// coordinates of bin center
		x = float64(i % o.N[0])
		y = float64(i % (o.N[0] * o.N[1]) / o.N[0])
		z = float64(i / (o.N[0] * o.N[1]))
		x = (x + 0.5) * o.S
		y = (y + 0.5) * o.S
		z = (z + 0.5) * o.S

		p = Point{x, y, z}

		d := DistPointLine(&p, &pi, &pf, tol, false)
		if d <= btol {
			sbins = append(sbins, bin)
		}
	}

	// entries ids
	var ids []int

	// find closest points
	for _, bin := range sbins {
		for _, entry := range bin.Entries {
			x = entry.x[0]
			y = entry.x[1]
			if o.Ndim == 3 {
				z = entry.x[0]
			}
			p := Point{x, y, z}

			d := DistPointLine(&p, &pi, &pf, tol, false)
			if d <= tol {
				ids = append(ids, entry.Id)
			}
		}
	}
	return ids
}

func (o Bin) String() string {
	l := utl.Sf("{\"idx\":%d, \"entries\":[", o.Idx)
	for i, entry := range o.Entries {
		if i > 0 {
			l += ", "
		}
		l += utl.Sf("{\"id\":%d, \"x\":[%g,%g", entry.Id, entry.x[0], entry.x[1])
		if len(entry.x) > 2 {
			l += utl.Sf(",%g", entry.x[2])
		}
		l += "]}"
	}
	l += "]}"
	return l
}

func (o Bins) String() string {
	l := "[\n"
	k := 0
	for _, bin := range o.All {
		if bin != nil {
			if k > 0 {
				l += ",\n"
			}
			l += utl.Sf("  %v", bin)
			k += 1
		}
	}
	l += "\n]"
	return l
}
