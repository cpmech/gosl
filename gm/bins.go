// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"gosl/chk"
	"gosl/io"
	"gosl/utl"
)

// consts
var (
	XDELZERO = 1e-10 // minimum distance between coordinates; i.e. xmax[i]-xmin[i] mininum
)

// BinEntry holds data of an entry to bin
type BinEntry struct {
	ID    int         // object Id
	X     []float64   // entry coordinate (read only)
	Extra interface{} // any entity attached to this entry
}

// Bin defines one bin in Bins (holds entries for search)
type Bin struct {
	Index   int         // index of bin
	Entries []*BinEntry // entries
}

// Bins implements a set of bins holding entries and is used to fast search entries by given coordinates.
type Bins struct {
	Ndim int       // space dimension
	Xmin []float64 // [ndim] left/lower-most point
	Xmax []float64 // [ndim] right/upper-most point
	Xdel []float64 // [ndim] the lengths along each direction (whole box)
	Size []float64 // size of bins
	Ndiv []int     // [ndim] number of divisions along each direction
	All  []*Bin    // [nbins] all bins (there will be an extra "ghost" bin along each dimension)
	tmp  []int     // [ndim] temporary (auxiliary) slice
}

// Init initialise Bins structure
//   xmin -- [ndim] min/initial coordinates of the whole space (box/cube)
//   xmax -- [ndim] max/final coordinates of the whole space (box/cube)
//   ndiv -- [ndim] number of divisions for xmax-xmin
func (o *Bins) Init(xmin, xmax []float64, ndiv []int) {

	// check for out-of-range values
	o.Ndim = len(xmin)
	o.Xmin = xmin
	o.Xmax = xmax
	chk.IntAssert(len(xmin), len(xmax))
	chk.IntAssert(len(xmin), len(ndiv))
	if len(xmin) < 2 || len(xmin) > 3 {
		chk.Panic("sizes of xmin and xmax must be the same and equal to either 2 or 3\n")
	}

	// allocate slices with max lengths and number of division
	o.Xdel = make([]float64, o.Ndim)
	o.Size = make([]float64, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		if ndiv[k] < 1 {
			ndiv[k] = 1
		}
		o.Xdel[k] = o.Xmax[k] - o.Xmin[k]
		o.Size[k] = o.Xdel[k] / float64(ndiv[k])
		if o.Xdel[k] < XDELZERO {
			chk.Panic("xmax[%d]-xmin[%d]=%g must be greater than %g\n", k, k, o.Xdel[k], XDELZERO)
		}
	}

	// number of divisions
	o.Ndiv = make([]int, o.Ndim)
	nbins := 1
	for k := 0; k < o.Ndim; k++ {
		o.Ndiv[k] = int(o.Xdel[k] / o.Size[k])
		nbins *= o.Ndiv[k]
	}

	// other slices
	o.All = make([]*Bin, nbins)
	o.tmp = make([]int, o.Ndim)
}

// Append adds a new entry {x, id, something} into the Bins structure
func (o *Bins) Append(x []float64, id int, extra interface{}) {
	idx := o.CalcIndex(x)
	if idx < 0 {
		chk.Panic("coordinates %v are out of range\n", x)
	}
	bin := o.FindBinByIndex(idx)
	if bin == nil {
		chk.Panic("bin index %v is out of range\n", idx)
	}
	xcopy := utl.GetCopy(x)
	entry := BinEntry{id, xcopy, extra}
	bin.Entries = append(bin.Entries, &entry)
}

// Clear clears all bins
func (o *Bins) Clear() {
	o.All = make([]*Bin, 0)
}

// FindBinByIndex finds or allocate new bin corresponding to index idx
func (o Bins) FindBinByIndex(idx int) *Bin {

	// check
	if idx < 0 || idx >= len(o.All) {
		return nil
	}

	// allocate new bin if necessary
	if o.All[idx] == nil {
		o.All[idx] = new(Bin)
		o.All[idx].Index = idx
	}
	return o.All[idx]
}

// CalcIndex calculates the bin index where the point x is
// returns -1 if out-of-range
func (o Bins) CalcIndex(x []float64) int {
	for k := 0; k < o.Ndim; k++ {
		if x[k] < o.Xmin[k] || x[k] > o.Xmax[k] {
			return -1
		}
		o.tmp[k] = int((x[k] - o.Xmin[k]) / o.Size[k])
		if o.tmp[k] == o.Ndiv[k] { // it's exactly on max edge => select inner bin
			o.tmp[k]-- // move to the inside
		}
	}
	idx := o.tmp[0] + o.tmp[1]*o.Ndiv[0]
	if o.Ndim > 2 {
		idx += o.tmp[2] * o.Ndiv[0] * o.Ndiv[1]
	}
	return idx
}

// high-level functions ///////////////////////////////////////////////////////////////////////////

// FindClosest returns the id of the entry whose coordinates are closest to x
//   idClosest -- the id of the closest entity. return -1 if out-of-range or not found
//   sqDistMin -- the minimum distance (squared) between x and the closest entity in the same bin
//
//  NOTE: FindClosest does search the whole area.
//        It only locates neighbours in the same bin where the given x is located.
//        So, if there area no entries in the bin containing x, no entry will be found.
//
func (o Bins) FindClosest(x []float64) (idClosest int, sqDistMin float64) {

	// set "not-found" results
	idClosest = -1
	sqDistMin = math.Inf(+1)

	// index and check
	idx := o.CalcIndex(x)
	if idx < 0 { // out-of-range
		return
	}

	// search for the closest point
	bin := o.FindBinByIndex(idx)
	idClosest = -1
	var entry *BinEntry
	for _, entry = range bin.Entries {
		var d float64
		for k := 0; k < o.Ndim; k++ {
			d += math.Pow(x[k]-entry.X[k], 2)
		}
		if d < sqDistMin {
			idClosest = entry.ID
			sqDistMin = d
		}
	}
	return
}

// FindClosestAndAppend finds closest point and, if not found, append to bins with a new Id
//   Input:
//     nextId -- is the Id of the next point. Will be incremented if x is a new point to be added.
//     x      -- is the point to be added
//     extra  -- extra information attached to point
//     radTol -- is the tolerance for the radial distance (i.e. NOT squared) to decide
//               whether a new point will be appended or not.
//     diff   -- [optional] a function for further check that the new and an eventual existent
//               points are really different, even after finding that they coincide (within tol)
//   Output:
//     id       -- the id attached to x
//     existent -- flag telling if x was found, based on given tolerance
func (o *Bins) FindClosestAndAppend(nextID *int, x []float64, extra interface{}, radTol float64, diff func(idOld int, xNew []float64) bool) (id int, existent bool) {

	// try to find another close point
	idClosest, sqDistMin := o.FindClosest(x)

	// new point for sure; i.e no other point was found
	id = *nextID
	if idClosest < 0 {
		o.Append(x, id, extra)
		(*nextID)++
		return
	}

	// new point, distant from the point just found
	dist := math.Sqrt(sqDistMin)
	if dist > radTol {
		o.Append(x, id, extra)
		(*nextID)++
		return
	}

	// further check
	if diff != nil {
		if diff(idClosest, x) {
			o.Append(x, id, extra)
			(*nextID)++
			return
		}
	}

	// existent point
	id = idClosest
	existent = true
	return
}

// FindAlongSegment gets the ids of entries that lie close to a segment
//  Note: the initial (xi) and final (xf) points on segment define a bounding box to filter points
func (o Bins) FindAlongSegment(xi, xf []float64, tol float64) []int {

	// auxiliary variables
	var sbins []*Bin // selected bins
	lmax := utl.Max(o.Size[0], o.Size[1])
	if o.Ndim == 3 {
		lmax = utl.Max(lmax, o.Size[2])
	}
	btol := 0.9 * lmax // tolerance for bins
	var p, pi, pf Point
	pi.X = xi[0]
	pf.X = xf[0]
	pi.Y = xi[1]
	pf.Y = xf[1]
	if o.Ndim == 3 {
		pi.Z = xi[2]
		pf.Z = xf[2]
	} else {
		xi = []float64{xi[0], xi[1], -1}
		xf = []float64{xf[0], xf[1], 1}
	}

	// loop along all bins
	var i, j, k int
	var x, y, z float64
	nxy := o.Ndiv[0] * o.Ndiv[1]
	for idx, bin := range o.All {

		// skip empty bins
		if bin == nil {
			continue
		}

		// coordinates of bin center
		i = idx % o.Ndiv[0] // indices representing bin
		j = (idx % nxy) / o.Ndiv[0]
		x = o.Xmin[0] + float64(i)*o.Size[0] // coordinates of bin corner
		y = o.Xmin[1] + float64(j)*o.Size[1]
		x += o.Size[0] / 2.0
		y += o.Size[1] / 2.0
		if o.Ndim == 3 {
			k = idx / nxy
			z = o.Xmin[2] + float64(k)*o.Size[2]
			z += o.Size[2] / 2.0
		}

		// check if bin is near line
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
			x = entry.X[0]
			y = entry.X[1]
			if o.Ndim == 3 {
				z = entry.X[0]
			}
			p := Point{x, y, z}
			d := DistPointLine(&p, &pi, &pf, tol, false)
			if d <= tol {
				if IsPointIn(&p, xi, xf, tol) {
					ids = append(ids, entry.ID)
				}
			}
		}
	}
	return ids
}

// GetLimits returns limigs of a specific bin
func (o *Bins) GetLimits(idxBin int) (xmin, xmax []float64) {
	nxy := o.Ndiv[0] * o.Ndiv[1]
	i := idxBin % o.Ndiv[0]
	j := (idxBin % nxy) / o.Ndiv[0]
	if o.Ndim == 2 {
		xmin = []float64{o.Xmin[0] + float64(i+0)*o.Size[0], o.Xmin[1] + float64(j+0)*o.Size[1]}
		xmax = []float64{o.Xmin[0] + float64(i+1)*o.Size[0], o.Xmin[1] + float64(j+1)*o.Size[1]}
	} else {
		k := idxBin / nxy
		xmin = []float64{o.Xmin[0] + float64(i+0)*o.Size[0], o.Xmin[1] + float64(j+0)*o.Size[1], o.Xmin[2] + float64(k+0)*o.Size[2]}
		xmax = []float64{o.Xmin[0] + float64(i+1)*o.Size[0], o.Xmin[1] + float64(j+1)*o.Size[1], o.Xmin[2] + float64(k+1)*o.Size[2]}
	}
	return
}

// information ////////////////////////////////////////////////////////////////////////////////////

// Nactive returns the number of active bins; i.e. non-nil bins
func (o *Bins) Nactive() (nactive int) {
	for _, bin := range o.All {
		if bin != nil {
			nactive++
		}
	}
	return
}

// Nentries returns the total number of entries (in active bins)
func (o *Bins) Nentries() (nentries int) {
	for _, bin := range o.All {
		if bin != nil {
			nentries += len(bin.Entries)
		}
	}
	return
}

// Summary returns the summary of this Bins' data
func (o *Bins) Summary() (l string) {
	l += io.Sf("Ndim        = %v\n", o.Ndim)
	l += io.Sf("Xmin        = %v\n", o.Xmin)
	l += io.Sf("Xmax        = %v\n", o.Xmax)
	l += io.Sf("Xdel        = %v\n", o.Xdel)
	l += io.Sf("Size        = %v\n", o.Size)
	l += io.Sf("Ndiv        = %v\n", o.Ndiv)
	l += io.Sf("Nalll       = %d\n", len(o.All))
	l += io.Sf("NactiveBins = %d\n", o.Nactive())
	l += io.Sf("Nentries    = %d\n", o.Nentries())
	return
}

// String returns the string representation of one Bin
func (o Bin) String() string {
	l := io.Sf("{\"idx\":%d, \"entries\":[", o.Index)
	for i, entry := range o.Entries {
		if i > 0 {
			l += ", "
		}
		l += io.Sf("{\"id\":%d, \"x\":[%g,%g", entry.ID, entry.X[0], entry.X[1])
		if len(entry.X) > 2 {
			l += io.Sf(",%g", entry.X[2])
		}
		l += "]"
		if entry.Extra != nil {
			l += io.Sf(", \"extra\":true")
		}
		l += "}"
	}
	l += "]}"
	return l
}

// String returns the string representation of a set of Bins
func (o Bins) String() string {
	l := "[\n"
	k := 0
	for _, bin := range o.All {
		if bin != nil {
			if k > 0 {
				l += ",\n"
			}
			l += io.Sf("  %v", bin)
			k++
		}
	}
	l += "\n]"
	return l
}
