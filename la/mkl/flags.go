// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mkl

/*
#cgo CFLAGS: -O2 -DMKL_ILP64 -m64 -I/opt/intel/mkl/include
#cgo LDFLAGS: -Wl,--start-group /opt/intel/mkl/lib/intel64/libmkl_intel_ilp64.a /opt/intel/mkl/lib/intel64/libmkl_intel_thread.a /opt/intel/mkl/lib/intel64/libmkl_core.a /opt/intel/lib/intel64/libiomp5.a -Wl,--end-group -lpthread -lm -ldl -L/opt/intel/lib/intel64
*/
import "C"
