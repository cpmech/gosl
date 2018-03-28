// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package h5

/*
#cgo linux CFLAGS: -I/usr/include/hdf5/serial -D_LARGEFILE64_SOURCE -D_LARGEFILE_SOURCE -D_FORTIFY_SOURCE=2 -g -O2 -Wformat -Werror=format-security
#cgo linux LDFLAGS: -L/usr/lib/x86_64-linux-gnu/hdf5/serial -lhdf5_hl -lhdf5 -lpthread -lz -ldl -lm

#cgo darwin CFLAGS: -I/usr/local/opt/szip/include
#cgo darwin LDFLAGS: -L/usr/local/opt/szip/lib -L/usr/local/Cellar/hdf5/1.10.1_2/lib -lhdf5_hl -lhdf5 -lz -ldl -lm
*/
import "C"

// NOTE: get flags with:
//
//   h5cc -show
