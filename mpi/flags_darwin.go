// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!linux

package mpi

/*
#cgo CFLAGS: -O2 -I/usr/lib/openmpi/include/openmpi/opal/mca/event/libevent2021/libevent -I/usr/lib/openmpi/include/openmpi/opal/mca/event/libevent2021/libevent/include -I/usr/lib/openmpi/include -I/usr/lib/openmpi/include/openmpi -I/usr/include/openmpi -pthread
#cgo LDFLAGS: -pthread -lmpi
*/
import "C"

// NOTE: get flags with:
//
//   mpicc hello_c.c -showme:compile
//
//   mpicc hello_c.c -showme:link
