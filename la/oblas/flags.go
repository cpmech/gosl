// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

/*
#cgo linux CFLAGS: -O2 -I/usr/local/include
#cgo linux LDFLAGS: -lopenblas -llapacke -L/local/lib

#cgo windows CFLAGS: -O2
#cgo windows LDFLAGS: -lopenblas -lgfortran

#cgo darwin CFLAGS: -I/usr/local/opt/openblas/include
#cgo darwin LDFLAGS: -lopenblas -L/usr/local/opt/openblas/lib
*/
import "C"
