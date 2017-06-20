// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oblas

/*
#cgo linux   CFLAGS: -DOPENBLAS_USE64BITINT -O2 -I/usr/include
#cgo linux   CFLAGS: -DOPENBLAS_USE64BITINT -O2 -I/usr/local/include
#cgo windows CFLAGS: -DOPENBLAS_USE64BITINT -O2 -IC:/Gosl/include

#cgo linux   LDFLAGS: -lopenblas -L/local/lib
#cgo darwin  LDFLAGS: -lopenblas -L/usr/local/lib
#cgo windows LDFLAGS: -lopenblas -LC:/Gosl/lib
*/
import "C"
