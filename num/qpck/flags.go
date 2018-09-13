// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin

package qpck

/*
#cgo linux LDFLAGS: -lopenblas -llapack -lgfortran -lm

#cgo windows LDFLAGS: -lopenblas -lgfortran -lm
*/
import "C"
