#!/bin/bash

GP="${GOPATH%:*}"

go build -o /tmp/NurbsDraw NurbsDraw.go && mv /tmp/NurbsDraw $GP/bin/
                                        
go build -o /tmp/GoEqs2TeX GoEqs2TeX.go && mv /tmp/GoEqs2TeX $GP/bin/
