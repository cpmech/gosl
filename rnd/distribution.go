// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import "github.com/cpmech/gosl/chk"

// Distribution defines a probability distribution
type Distribution interface {
	Init(prms *VarData) error
	Pdf(x float64) float64
	Cdf(x float64) float64
}

// factory
var distallocators = make(map[string]func() Distribution)

// GetDistrib returns a distribution from factory
func GetDistrib(name string) (d Distribution, err error) {
	allocator, ok := distallocators[name]
	if !ok {
		return nil, chk.Err("cannot find distribution named: %s", name)
	}
	return allocator(), nil
}
