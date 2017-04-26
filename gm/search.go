// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// consts
var (
	XDELZERO = 1e-10 // minimum distance between coordinates; i.e. xmax[i]-xmin[i] mininum
)

// BinEntry holds data of an entry to bin
type BinEntry struct {
	Id    int         // object Id
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
	Npts []int     // [ndim] number of points along each direction; i.e. ndiv + 1
	All  []*Bin    // [nbins] all bins (there will be an extra "ghost" bin along each dimension)
	tmp  []int     // [ndim] temporary (auxiliary) slice
}

// Init initialise Bins structure
//   xmin -- [ndim] min/initial coordinates of the whole space (box/cube)
//   xmax -- [ndim] max/final coordinates of the whole space (box/cube)
//   ndiv -- number of divisions for xmax-xmin
func (o *Bins) Init(xmin, xmax []float64, ndiv int) (err error) {

	// check for out-of-range values
	o.Ndim = len(xmin)
	o.Xmin = xmin
	o.Xmax = xmax
	if len(xmin) != len(xmax) || len(xmin) < 2 || len(xmin) > 3 {
		return chk.Err("sizes of xmin and xmax must be the same and equal to either 2 or 3")
	}
	if ndiv < 1 {
		ndiv = 1
	}

	// allocate slices with max lengths and number of division
	o.Xdel = make([]float64, o.Ndim)
	o.Size = make([]float64, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		o.Xdel[k] = o.Xmax[k] - o.Xmin[k]
		o.Size[k] = o.Xdel[k] / float64(ndiv)
		if o.Xdel[k] < XDELZERO {
			return chk.Err("xmax[%d]-xmin[%d]=%g must be greater than %g", k, k, o.Xdel[k], XDELZERO)
		}
	}

	// number of divisions
	o.Npts = make([]int, o.Ndim)
	nbins := 1
	for k := 0; k < o.Ndim; k++ {
		o.Npts[k] = int(o.Xdel[k]/o.Size[k]) + 1
		nbins *= o.Npts[k]
	}

	// other slices
	o.All = make([]*Bin, nbins)
	o.tmp = make([]int, o.Ndim)
	return
}

// Append adds a new entry {x, id, something} into the Bins structure
func (o *Bins) Append(x []float64, id int, extra interface{}) (err error) {
	idx := o.CalcIndex(x)
	if idx < 0 {
		return chk.Err("coordinates %v are out of range", x)
	}
	bin := o.FindBinByIndex(idx)
	if bin == nil {
		return chk.Err("bin index %v is out of range", idx)
	}
	xcopy := utl.GetCopy(x)
	entry := BinEntry{id, xcopy, extra}
	bin.Entries = append(bin.Entries, &entry)
	return
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

// CalcIdx calculates the bin index where the point x is
// returns -1 if out-of-range
func (o Bins) CalcIndex(x []float64) int {
	for k := 0; k < o.Ndim; k++ {
		if x[k] < o.Xmin[k] || x[k] > o.Xmax[k] {
			return -1
		}
		o.tmp[k] = int((x[k] - o.Xmin[k]) / o.Size[k])
	}
	idx := o.tmp[0] + o.tmp[1]*o.Npts[0]
	if o.Ndim > 2 {
		idx += o.tmp[2] * o.Npts[0] * o.Npts[1]
	}
	return idx
}

// high-level functions ///////////////////////////////////////////////////////////////////////////

// FindClosest returns the id of the entry whose coordinates are closest to x
//   idClosest -- the id of the closest entity. return -1 if out-of-range or not found
//   sqDistMin -- the minimum distance (squared) between x and the closest entity in the same bin
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
			idClosest = entry.Id
			sqDistMin = d
		}
	}
	return
}

// FindClosestAndAppend finds closest point and, if not found, append to bins with given Id
//   radTol -- is the tolerance for the radial distance (i.e. NOT squared) to decide
//             whether a new entry will be appended or not.
//   returns the next Id which will be either the input Id, or the input Id incremented by one.
func (o *Bins) FindClosestAndAppend(id int, x []float64, extra interface{}, radTol float64) int {

	// try to find another close point
	idClosest, sqDistMin := o.FindClosest(x)

	// new point for sure; i.e no other point was found
	if idClosest < 0 {
		o.Append(x, id, extra)
		return id + 1
	}

	// new point, distant from found one by radTol
	dist := math.Sqrt(sqDistMin)
	if dist > radTol {
		o.Append(x, id, extra)
		return id + 1
	}

	// existent point, within tolerance
	return id
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
	nxy := o.Npts[0] * o.Npts[1]
	for idx, bin := range o.All {

		// skip empty bins
		if bin == nil {
			continue
		}

		// coordinates of bin center
		i = idx % o.Npts[0] // indices representing bin
		j = (idx % nxy) / o.Npts[0]
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
					ids = append(ids, entry.Id)
				}
			}
		}
	}
	return ids
}

// plotting ///////////////////////////////////////////////////////////////////////////////////////

// Draw draws bins; i.e. grid
func (o *Bins) Draw(withEntry, withGrid, withEntryTxt, withGridTxt bool, argsEntry, argsGrid, argsTxtEntry, argsTxtGrid *plt.A, selBins map[int]bool) {

	// grid
	if withGrid {

		// configuration
		if argsGrid == nil {
			argsGrid = &plt.A{C: "#427ce5", Lw: 0.8, NoClip: true}
		}

		// x-y coordinates
		X := make([][]float64, o.Npts[0])
		Y := make([][]float64, o.Npts[0])
		for i := 0; i < o.Npts[0]; i++ {
			X[i] = make([]float64, o.Npts[1])
			Y[i] = make([]float64, o.Npts[1])
			for j := 0; j < o.Npts[1]; j++ {
				X[i][j] = o.Xmin[0] + float64(i)*o.Size[0]
				Y[i][j] = o.Xmin[1] + float64(j)*o.Size[1]
			}
		}

		// draw grid
		if o.Ndim == 2 {
			plt.Grid2d(X, Y, false, argsGrid, nil)
		} else {
			Zlevels := make([]float64, o.Npts[2])
			for k := 0; k < o.Npts[2]; k++ {
				Zlevels[k] = o.Xmin[2] + float64(k)*o.Size[2]
			}
			plt.Grid3d(X, Y, Zlevels, argsGrid)
		}
	}

	// selected bins
	if o.Ndim == 2 {
		nxy := o.Npts[0] * o.Npts[1]
		for idx, _ := range selBins {
			i := idx % o.Npts[0] // indices representing bin
			j := (idx % nxy) / o.Npts[0]
			x := o.Xmin[0] + float64(i)*o.Size[0] // coordinates of bin corner
			y := o.Xmin[1] + float64(j)*o.Size[1]
			plt.Polyline([][]float64{
				{x, y},
				{x + o.Size[0], y},
				{x + o.Size[0], y + o.Size[1]},
				{x, y + o.Size[1]},
			}, &plt.A{Fc: "#fbefdc", Ec: "#8e8371", Lw: 0.5, Closed: true, NoClip: true})
		}
	}

	// plot items
	if withEntry {

		// configuration
		if argsEntry == nil {
			argsEntry = &plt.A{C: "r", M: ".", NoClip: true}
			if o.Ndim == 3 {
				argsEntry.M = "o"
			}
		}
		if argsTxtEntry == nil {
			argsTxtEntry = &plt.A{C: "g", Fsz: 8}
		}
		argsTxtEntry.Ha = "right"

		// draw markers indicating entries
		nentries := o.Nentries()
		X := make([]float64, nentries)
		Y := make([]float64, nentries)
		var Z []float64
		if o.Ndim == 3 {
			Z = make([]float64, nentries)
		}
		k := 0
		for _, bin := range o.All {
			if bin != nil {
				for _, entry := range bin.Entries {
					X[k] = entry.X[0]
					Y[k] = entry.X[1]
					if o.Ndim == 3 {
						Z[k] = entry.X[2]
						if withEntryTxt {
							plt.Text3d(X[k], Y[k], Z[k], io.Sf("%d", entry.Id), argsTxtEntry)
						}
					} else {
						if withEntryTxt {
							plt.Text(X[k], Y[k], io.Sf("%d", entry.Id), argsTxtEntry)
						}
					}
					k++
				}
			}
		}
		if o.Ndim == 2 {
			argsEntry.Ls = "none"
			plt.Plot(X, Y, argsEntry)
		} else {
			plt.Plot3dPoints(X, Y, Z, argsEntry)
		}
	}

	// grid txt
	if withGridTxt {

		// configuration
		if argsTxtGrid == nil {
			argsTxtGrid = &plt.A{C: "orange", Fsz: 8}
		}

		// add text
		n2 := 1
		if o.Ndim == 3 {
			n2 = o.Npts[2]
		}
		for k := 0; k < n2; k++ {
			z := 0.0
			if o.Ndim == 3 {
				z = o.Xmin[2] + float64(k)*o.Size[2] + 0.02*o.Size[2]
			}
			for j := 0; j < o.Npts[1]; j++ {
				for i := 0; i < o.Npts[0]; i++ {
					idx := i + j*o.Npts[0]
					if o.Ndim == 3 {
						idx += k * o.Npts[0] * o.Npts[1]
					}
					x := o.Xmin[0] + float64(i)*o.Size[0] + 0.02*o.Size[0]
					y := o.Xmin[1] + float64(j)*o.Size[1] + 0.02*o.Size[1]
					txt := io.Sf("%d", idx)
					if o.Ndim == 3 {
						plt.Text3d(x, y, z, txt, argsTxtGrid)
					} else {
						plt.Text(x, y, txt, argsTxtGrid)
					}
				}
			}
		}
	}
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
	l += io.Sf("Npts        = %v\n", o.Npts)
	l += io.Sf("Nalll(+ext) = %d\n", len(o.All))
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
		l += io.Sf("{\"id\":%d, \"x\":[%g,%g", entry.Id, entry.X[0], entry.X[1])
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
			k += 1
		}
	}
	l += "\n]"
	return l
}
