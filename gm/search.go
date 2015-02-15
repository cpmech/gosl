// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import "github.com/cpmech/gosl/utl"

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

/*
//    FAST_32_hash
//    A very fast 2D hashing function.  Requires 32bit support.
//
//    The hash formula takes the form....
//    hash = mod( coord.x * coord.x * coord.y * coord.y, SOMELARGEFLOAT ) / SOMELARGEFLOAT
//    We truncate and offset the domain to the most interesting part of the noise.
//
vec4 FAST_32_hash( vec2 gridcell )
{
//    gridcell is assumed to be an integer coordinate
const vec2 OFFSET = vec2( 26.0, 161.0 );
const float DOMAIN = 71.0;
const float SOMELARGEFLOAT = 951.135664;
vec4 P = vec4( gridcell.xy, gridcell.xy + 1.0.xx );
P = P - floor(P * ( 1.0 / DOMAIN )) * DOMAIN;    //    truncate the domain
P += OFFSET.xyxy;                                //    offset to interesting part of the noise
P *= P;                                          //    calculate and return the hash
return fract( P.xzxz * P.yyww * ( 1.0 / SOMELARGEFLOAT.x ).xxxx );
}
*/

type BinEntry struct {
	Id  int // object Id
	Idx int // index in Entries slice
}

type Bins struct {
	Ndim      int         // space dimension
	Xi        []float64   // [ndim] left/lower-most point
	Xf        []float64   // [ndim] right/upper-most point
	L         []float64   // [ndim] whole box lengths
	N         []int       // [ndim] number of divisions
	S         float64     // size of bins
	Entries   []*BinEntry // entries
	Idx2entry []int       // maps idx to entry in Entries
	tmp       []int       // [ndim] temporary (auxiliary) slice
}

// xi   -- [ndim] initial positions
// l    -- [ndim] whole box lengths
// ndiv -- number of divisions
func (o *Bins) Init(xi, l []float64, ndiv int) (err error) {
	o.Ndim = len(xi)
	o.Xi, o.L = xi, l
	if len(xi) != len(l) || len(l) < 2 || len(l) > 3 {
		return utl.Err("sizes of xi and l must be the same and equal to either 2 or 3")
	}
	o.Xf = make([]float64, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		o.Xf[k] = o.Xi[k] + o.L[k]
		o.N[k] = int(o.L[k])/ndiv + 1
	}
	return
}

func (o *Bins) Append(x []float64, id int) (err error) {
	idx := o.CalcIdx(x)
	if idx < 0 {
		return utl.Err("point %v is out of range", x)
	}
	entry := o.FindByIndex(idx)
	if entry == nil {
		o.Entries = append(o.Entries, &BinEntry{id, idx})
		return
	}
	entry.Id = id
	return
}

func (o *Bins) Clear() {
	o.Entries = make([]*BinEntry, 0)
	o.Idx2entry = make([]int, 0)
}

func (o Bins) FindByIndex(idx int) *BinEntry {
	if idx < 0 || idx >= len(o.Idx2entry) {
		return nil
	}
	i := o.Idx2entry[idx]
	return o.Entries[i]
}

func (o Bins) FindByCoords(x []float64) *BinEntry {
	idx := o.CalcIdx(x)
	if idx < 0 {
		return nil // out-of-range
	}
	i := o.Idx2entry[idx]
	return o.Entries[i]
}

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
