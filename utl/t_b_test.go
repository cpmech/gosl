// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"math/rand"
	"testing"

	"github.com/cpmech/gosl/io"
)

var (
	benchStrindexsmall []string
	benchStrindexmap   map[string]int
	benchResult        int
)

func init() {
	rand.Seed(13)
	answers := []string{"66", "644", "666", "653", "10", "0", "1", "1", "1", "9", "A", ""}
	//n := 1000000
	//m := 500000
	//n := 100
	//m := 50
	n := 20
	m := 10
	benchStrindexsmall = make([]string, n)
	benchStrindexmap = make(map[string]int)
	for i := 0; i < n; i++ {
		j := rand.Intn(len(answers))
		benchStrindexsmall[i] = answers[j]
		benchStrindexmap[io.Sf("%s_%d", answers[j], i)] = i
		if i == m {
			benchStrindexsmall[i] = "user"
			benchStrindexmap["user"] = i
		}
	}
	//io.Pforan("%v\n", ___benchmarking_strindexsmall___)
	//io.Pfcyan("%v\n", ___benchmarking_strindexmap___)
}

func BenchmarkStrIndexSmall(b *testing.B) {
	var idx int
	for i := 0; i < b.N; i++ {
		idx = StrIndexSmall(benchStrindexsmall, "user")
	}
	benchResult = idx
}

func BenchmarkStrIndexMap(b *testing.B) {
	var idx int
	for i := 0; i < b.N; i++ {
		idx = benchStrindexmap["user"]
	}
	benchResult = idx
}
