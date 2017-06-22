// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

/*
#cgo linux CFLAGS: -O2 -I/usr/include/suitesparse -I/usr/local/include/suitesparse
#cgo linux LDFLAGS: -L/usr/lib -L/usr/local/lib
#cgo linux LDFLAGS: -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig -lopenblas -lgfortran
#cgo linux LDFLAGS: -ldmumps -lzmumps -lmumps_common -lpord

#cgo windows CFLAGS: -O2
#cgo windows LDFLAGS: -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig -lopenblas -lgfortran

#cgo darwin CFLAGS: -I/usr/local/opt/openblas/include
#cgo darwin LDFLAGS: -L/usr/local/opt/openblas/lib
#cgo darwin LDFLAGS: -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig -lopenblas
*/
import "C"
