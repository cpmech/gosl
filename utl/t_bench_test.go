// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math/rand"
	"testing"

	"github.com/cpmech/gosl/io"
)

var ___benchmarking_strindexsmall___ []string
var ___benchmarking_strindexmap___ map[string]int
var ___benchmarking_result___ int

func init() {
	rand.Seed(13)
	answers := []string{"66", "644", "666", "653", "10", "0", "1", "1", "1", "9", "A", ""}
	//n := 1000000
	//m := 500000
	//n := 100
	//m := 50
	n := 20
	m := 10
	___benchmarking_strindexsmall___ = make([]string, n)
	___benchmarking_strindexmap___ = make(map[string]int)
	for i := 0; i < n; i++ {
		j := rand.Intn(len(answers))
		___benchmarking_strindexsmall___[i] = answers[j]
		___benchmarking_strindexmap___[io.Sf("%s_%d", answers[j], i)] = i
		if i == m {
			___benchmarking_strindexsmall___[i] = "user"
			___benchmarking_strindexmap___["user"] = i
		}
	}
	//io.Pforan("%v\n", ___benchmarking_strindexsmall___)
	//io.Pfcyan("%v\n", ___benchmarking_strindexmap___)
}

func BenchmarkStrIndexSmall(b *testing.B) {
	var idx int
	for i := 0; i < b.N; i++ {
		idx = StrIndexSmall(___benchmarking_strindexsmall___, "user")
	}
	___benchmarking_result___ = idx
}

func BenchmarkStrIndexMap(b *testing.B) {
	var idx int
	for i := 0; i < b.N; i++ {
		idx = ___benchmarking_strindexmap___["user"]
	}
	___benchmarking_result___ = idx
}
