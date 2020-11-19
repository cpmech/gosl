// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

/*
#cgo CFLAGS: -O2 -I/usr/include/suitesparse -I/usr/local/include/suitesparse
#cgo LDFLAGS: -L/usr/lib -L/usr/local/lib
#cgo LDFLAGS: -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig -lopenblas -lgfortran
#cgo LDFLAGS: -ldmumps -lzmumps -lmumps_common -lpord
#cgo LDFLAGS: -lptesmumps -lptscotch -lptscotcherr -lparmetis -lmetis -lscalapack-openmpi
#cgo LDFLAGS: -lm -ldl -lgfortran
*/
import "C"
