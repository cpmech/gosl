// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rw

import "github.com/cpmech/gosl/io"

func atob(s string) bool {
	if s == ".t." || s == ".T." {
		return true
	}
	if s == ".f." || s == ".F." {
		return false
	}
	return io.Atob(s)
}
