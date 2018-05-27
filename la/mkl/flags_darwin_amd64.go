// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mkl

/*
#cgo CFLAGS: -O2  -DMKL_ILP64 -m64 -I/opt/intel/mkl/include
#cgo LDFLAGS: -L/opt/intel/lib /opt/intel/mkl/lib/libmkl_intel_ilp64.a /opt/intel/mkl/lib/libmkl_intel_thread.a /opt/intel/mkl/lib/libmkl_core.a /opt/intel/lib/libiomp5.a -lpthread -lm -ldl

*/
import "C"
