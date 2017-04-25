// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"bytes"
	"encoding/json"
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/utl"
)

// NurbsExchangeData holds all data required to exchange NURBS; e.g. read/save files
type NurbsExchangeData struct {
	Id    int         `json:"i"` // id of Nurbs
	Gnd   int         `json:"g"` // 1: curve, 2:surface, 3:volume (geometry dimension)
	Ords  []int       `json:"o"` // order along each x-y-z direction [gnd]
	Knots [][]float64 `json:"k"` // knots along each x-y-z direction [gnd][m]
	Ctrls []int       `json:"c"` // global ids of control points
}

// NurbsExchangeDataSet defines a set of nurbs exchange data
type NurbsExchangeDataSet []*NurbsExchangeData

// NurbsPatch defines a patch of many NURBS'
type NurbsPatch struct {

	// input
	Data NurbsExchangeDataSet `json:"patch"`
}

// WriteMsh writes .msh file
// Input:
//   vtagged -- maps hashed id of control point to vertex tag
//   ctagged -- maps idOfNurbs_localIdOfElem to cell tag
//   tol     -- tolerance for normalized numbers comparison when generating hashes, e.g. 1e-7
func WriteMsh(dirout, fnk string, nurbss []*Nurbs, vtagged map[int]int, ctagged map[string]int, tol float64) {

	// compute limits
	xmin, xmax, xdel := nurbss[0].GetLimitsQ()
	for r := 1; r < len(nurbss); r++ {
		xmi, xma, _ := nurbss[r].GetLimitsQ()
		for i := 0; i < 3; i++ {
			xmin[i] = utl.Min(xmin[i], xmi[i])
			xmax[i] = utl.Min(xmax[i], xma[i])
		}
	}
	if len(nurbss) > 1 {
		for i := 0; i < 3; i++ {
			xdel[i] = xmax[i] - xmin[i]
		}
	}

	var buf bytes.Buffer
	io.Ff(&buf, "{\n  \"verts\" : [\n")
	verts := make(map[int]int)
	vid := 0
	for _, o := range nurbss {
		for k := 0; k < o.n[2]; k++ {
			for j := 0; j < o.n[1]; j++ {
				for i := 0; i < o.n[0]; i++ {
					x := o.GetQ(i, j, k)
					hsh := HashPoint(x, xmin, xdel, tol)
					if _, ok := verts[hsh]; !ok {
						tag := 0
						if vtagged != nil {
							if val, tok := vtagged[hsh]; tok {
								tag = val
							}
						}

						// TODO: remove this
						if math.Abs(x[0]) < 1e-7 {
							tag = -100 // vertical
						}
						if math.Abs(x[1]) < 1e-7 {
							tag = -200 // horizontal
						}

						if len(verts) > 0 {
							io.Ff(&buf, ",\n")
						}
						io.Ff(&buf, "    { \"id\":%3d, \"tag\":%3d, \"c\":[%24.17e,%24.17e,%24.17e,%24.17e] }", vid, tag, x[0], x[1], x[2], x[3])
						verts[hsh] = vid
						vid += 1
					}
				}
			}
		}
	}
	io.Ff(&buf, "\n  ],\n  \"nurbss\" : [\n")
	for sid, o := range nurbss {
		if sid > 0 {
			io.Ff(&buf, ",\n")
		}
		io.Ff(&buf, "    { \"id\":%d, \"gnd\":%d, \"ords\":[%d,%d,%d],\n", sid, o.gnd, o.p[0], o.p[1], o.p[2])
		io.Ff(&buf, "      \"knots\":[\n")
		for d := 0; d < o.gnd; d++ {
			if d > 0 {
				io.Ff(&buf, ",\n")
			}
			io.Ff(&buf, "        [")
			for i, t := range o.b[d].T {
				if i > 0 {
					io.Ff(&buf, ",")
				}
				io.Ff(&buf, "%24.17e", t)
			}
			io.Ff(&buf, "]")
		}
		io.Ff(&buf, "\n      ],\n      \"ctrls\":[")
		first := true
		for k := 0; k < o.n[2]; k++ {
			for j := 0; j < o.n[1]; j++ {
				for i := 0; i < o.n[0]; i++ {
					if !first {
						io.Ff(&buf, ",")
					}
					x := o.GetQ(i, j, k)
					hsh := HashPoint(x, xmin, xdel, tol)
					io.Ff(&buf, "%d", verts[hsh])
					if first {
						first = false
					}
				}
			}
		}
		io.Ff(&buf, "]\n    }")
	}
	io.Ff(&buf, "\n  ],\n  \"cells\" : [\n")
	ndim := nurbss[0].gnd
	bry := make([]bool, 2*ndim)
	fti2d := []int{2, 1, 3, 0}
	ftags := []int{-10, -11, -20, -21, -30, -31}
	cid := 0
	for sid, o := range nurbss {
		spanmin := make([]int, ndim)
		spanmax := make([]int, ndim)
		for i := 0; i < ndim; i++ {
			spanmin[i] = o.p[i]
			spanmax[i] = o.n[i] // TODO: check this
		}
		elems := o.Elements()
		for eid, e := range elems {
			ibasis := o.IndBasis(e)
			if cid > 0 {
				io.Ff(&buf, ",\n")
			}
			tag := -1
			if ctagged != nil {
				if val, tok := ctagged[io.Sf("%d_%d", sid, eid)]; tok {
					tag = val
				}
			}

			// TODO: find a better way to tag cells
			if e[1] == spanmax[0] && e[2] == spanmin[1] {
				tag = -2
			}

			io.Ff(&buf, "    { \"id\":%3d, \"tag\":%2d, \"nrb\":%d, \"part\":0, \"type\":\"nurbs\",", cid, tag, sid)
			io.Ff(&buf, " \"span\":[")
			for k, idx := range e {
				if k > 0 {
					io.Ff(&buf, ",")
				}
				io.Ff(&buf, "%d", idx)
			}
			io.Ff(&buf, "], \"verts\":[")
			for i, l := range ibasis {
				if i > 0 {
					io.Ff(&buf, ",")
				}
				x := o.GetQl(l)
				hsh := HashPoint(x, xmin, xdel, tol)
				io.Ff(&buf, "%d", verts[hsh])
			}
			var onbry bool
			for i := 0; i < 2*ndim; i++ {
				bry[i] = false
			}
			for i := 0; i < ndim; i++ {
				smin, smax := e[i*ndim], e[i*ndim+1]
				if smin == spanmin[i] {
					bry[i*ndim] = true
					onbry = true
				}
				if smax == spanmax[i] {
					bry[i*ndim+1] = true
					onbry = true
				}
			}
			io.Ff(&buf, "]")
			if onbry && ndim > 1 {
				io.Ff(&buf, ", \"ftags\":[")
				for i := 0; i < 2*ndim; i++ {
					tag := 0
					I := i
					if ndim == 2 {
						I = fti2d[i]
					}
					if bry[I] {
						if ndim == 2 {
							tag = ftags[fti2d[i]]

							// TODO: replace this by a better approach
							if i == 2 {
								r := spanmin[0] + (spanmax[0]-spanmin[0])/2
								//io.Pforan("e[0]=%v r=%v\n", e[0], r)
								if e[0] >= r {
									tag = -24
								}
							}

						} else {
							tag = ftags[i]
						}
					}
					if i > 0 {
						io.Ff(&buf, ",")
					}
					io.Ff(&buf, "%d", tag)
				}
				io.Ff(&buf, "]")
			}
			io.Ff(&buf, " }")
			cid += 1
		}
	}
	io.Ff(&buf, "\n  ]\n}")
	io.WriteFileVD(dirout, fnk+".msh", &buf)
}

// Vert holds data for a vertex => control point
type Vert struct {
	Id  int       // id
	Tag int       // tag
	C   []float64 // coordinates (size==4)
}

// Data holds all geometry data
type Data struct {
	Verts  []Vert              // vertices
	Nurbss []NurbsExchangeData // NURBSs
}

// ReadMsh reads .msh file
func ReadMsh(fnk string) (nurbss []*Nurbs) {

	// read file
	fn := fnk + ".msh"
	buf, err := io.ReadFile(fn)
	if err != nil {
		chk.Panic("ReadMsh cannot read file = '%s', %v'", fn, err)
	}

	// decode
	var dat Data
	err = json.Unmarshal(buf, &dat)
	if err != nil {
		chk.Panic("ReadMsh cannot unmarshal file = '%s', %v'", fn, err)
	}

	// list of vertices
	verts := make([][]float64, len(dat.Verts))
	for i, _ := range dat.Verts {
		verts[i] = make([]float64, 4)
		for j := 0; j < 4; j++ {
			verts[i][j] = dat.Verts[i].C[j]
		}
	}

	// allocate NURBSs
	nurbss = make([]*Nurbs, len(dat.Nurbss))
	for i, v := range dat.Nurbss {
		nurbss[i] = new(Nurbs)
		nurbss[i].Init(v.Gnd, v.Ords, v.Knots)
		nurbss[i].SetControl(verts, v.Ctrls)
	}
	return
}

func tag_verts(b *Nurbs, tol float64) (vt map[int]int) {
	xmin, _, xdel := b.GetLimitsQ()
	vt = make(map[int]int)
	n0, n1 := b.NumBasis(0), b.NumBasis(1)
	for j := 0; j < n1; j++ {
		for i := 0; i < n0; i++ {
			x := b.GetQ(i, j, 0)
			if math.Abs(x[0]) < 1e-7 { // right
				vt[HashPoint(x, xmin, xdel, tol)] = -1
			}
			if math.Abs(x[1]) < 1e-7 { // bottom
				vt[HashPoint(x, xmin, xdel, tol)] = -2
			}
			if math.Abs(x[0]+4.0) < 1e-7 { // left
				vt[HashPoint(x, xmin, xdel, tol)] = -3
			}
			if math.Abs(x[0]+4.0) < 1e-7 && math.Abs(x[1]) < 1e-7 { // left-bottom
				vt[HashPoint(x, xmin, xdel, tol)] = -4
			}
		}
	}
	return
}
