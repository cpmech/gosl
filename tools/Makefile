# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: NurbsDraw GoEqs2TeX
.PHONY: NurbsDraw GoEqs2TeX

NurbsDraw: NurbsDraw.go
	go build NurbsDraw.go
	mv NurbsDraw $(GOPATH)/bin/

GoEqs2TeX: GoEqs2TeX.go
	go build GoEqs2TeX.go
	mv GoEqs2TeX $(GOPATH)/bin/
