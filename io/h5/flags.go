// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!darwin

package h5

/*
#cgo CFLAGS: -I/usr/include/hdf5/serial -D_LARGEFILE64_SOURCE -D_LARGEFILE_SOURCE -D_FORTIFY_SOURCE=2 -g -O2 -Wformat -Werror=format-security
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/hdf5/serial /usr/lib/x86_64-linux-gnu/hdf5/serial/libhdf5_hl.a /usr/lib/x86_64-linux-gnu/hdf5/serial/libhdf5.a -Wl,-Bsymbolic-functions -Wl,-z,relro -lpthread -lsz -lz -ldl -lm -Wl,-rpath -Wl,/usr/lib/x86_64-linux-gnu/hdf5/serial
*/
import "C"

// NOTE: get flags with:
//
//   h5cc -show
