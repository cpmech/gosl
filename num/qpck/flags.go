// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qpck

/*
#cgo linux LDFLAGS: -lopenblas -llapack -lgfortran -lm

#cgo windows LDFLAGS: -lopenblas -lgfortran -lm

#cgo darwin LDFLAGS: -L/usr/local/opt/openblas/lib -L/usr/local/Cellar/gcc/7.3.0_1/lib/gcc/7/ -lopenblas -lgfortran -lm
*/
import "C"
