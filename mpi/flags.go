// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package mpi

/*
#cgo linux CFLAGS: -I/usr/lib/x86_64-linux-gnu/openmpi/include/openmpi -I/usr/lib/x86_64-linux-gnu/openmpi/include -pthread
#cgo linux LDFLAGS: -pthread -L/usr/lib/x86_64-linux-gnu/openmpi/lib -lmpi

#cgo darwin CFLAGS: -I/usr/local/Cellar/open-mpi/4.0.1_2/include
#cgo darwin LDFLAGS: -L/usr/local/opt/libevent/lib -L/usr/local/Cellar/open-mpi/4.0.1_2/lib -lmpi
*/
import "C"
