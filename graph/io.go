// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"strings"

	"gosl/io"
)

// ReadGraphTable reads data and allocate graph
func ReadGraphTable(fname string, bargera bool) *Graph {

	// data
	var ne int
	var edges [][]int
	var weights []float64

	// Bar-Gera format files from: http://www.bgu.ac.il/~bargera/tntp/
	if bargera {
		k := 0
		readingMeta := true
		io.ReadLines(fname, func(idx int, line string) (stop bool) {
			if len(line) < 1 {
				return false
			}
			line = strings.TrimSpace(line)
			if line[0] == '~' {
				return false
			}
			if readingMeta {
				switch {
				case strings.HasPrefix(line, "<NUMBER OF LINKS>"):
					res := strings.Split(line, "<NUMBER OF LINKS>")
					ne = io.Atoi(strings.TrimSpace(res[1]))
					edges = make([][]int, ne)
					weights = make([]float64, ne)
				case strings.HasPrefix(line, "<END OF METADATA>"):
					readingMeta = false
				}
				return false
			}
			l := strings.Fields(line)
			edges[k] = []int{io.Atoi(l[0]) - 1, io.Atoi(l[1]) - 1}
			weights[k] = io.Atof(l[4])
			k++
			return false
		})
	} else {
		_, dat := io.ReadTable(fname)
		ne = len(dat["from"]) // number of edges
		edges = make([][]int, ne)
		weights = make([]float64, ne)
		for i := 0; i < ne; i++ {
			edges[i] = []int{int(dat["from"][i]) - 1, int(dat["to"][i]) - 1}
			weights[i] = dat["cost"][i]
		}
	}

	// graph
	var G Graph
	G.Init(edges, weights, nil, nil)
	return &G
}
