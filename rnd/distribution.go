// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "github.com/cpmech/gosl/chk"

// Distribution defines a probability distribution
type Distribution interface {
	Name() string
	Init(prms *Variable)
	Pdf(x float64) float64
	Cdf(x float64) float64
}

// factory
var distallocators = make(map[string]func() Distribution)

// GetDistrib returns a distribution from factory
func GetDistrib(dtype string) (d Distribution) {
	allocator, ok := distallocators[dtype]
	if !ok {
		chk.Panic("cannot find %q distribution\n", dtype)
	}
	return allocator()
}
