# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

bench:
	go test -run=XXX -bench=. -benchtime=5s

debug:
	go test -gcflags "-N -l" -c -o /tmp/test_gosl_utl

list02: debug
	gdb --args /tmp/test_gosl_utl -test.run=list02
