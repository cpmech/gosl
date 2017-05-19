// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

// NOTE: -lblas must come after -lumfpack on Windows
//       There is no need for -I and -L flags on Wndows
//       Because include and lib files will be installed
//       in C:\TDM-GCC-64, the default place for those files

/*
#cgo linux CFLAGS: -O2 -I/usr/include/suitesparse -I/usr/local/include/suitesparse
#cgo linux LDFLAGS: -L/usr/lib -L/usr/local/lib
#cgo linux LDFLAGS: -lm -llapack -lgfortran -lblas -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig
#cgo linux LDFLAGS: -ldmumps -lzmumps -lmumps_common -lpord

#cgo windows CFLAGS: -O2
#cgo windows LDFLAGS: -llapack -lgfortran -lumfpack -lblas -lamd -lcholmod -lcolamd -lsuitesparseconfig

#cgo darwin LDFLAGS: -L/usr/local/lib
#cgo darwin LDFLAGS: -lm -llapack -lblas -lumfpack -lamd -lcholmod -lcolamd -lsuitesparseconfig
*/
import "C"
