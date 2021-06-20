// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

/*
#cgo CFLAGS: -O2 -I/usr/include/mumps -I/usr/include/suitesparse
#cgo LDFLAGS: -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig -lopenblas -lgfortran
#cgo LDFLAGS: -ldmumps_seq -lzmumps_seq -lmumps_common_seq -lpord
*/
import "C"
