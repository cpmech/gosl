#!/bin/bash

GP="${GOPATH%:*}"

echo "...installing tools to $GP"

go build -o /tmp/NurbsDraw NurbsDraw.go && mv /tmp/NurbsDraw $GP/bin/
echo "......NurbsDraw installed"
                                        
go build -o /tmp/GoEqs2TeX GoEqs2TeX.go && mv /tmp/GoEqs2TeX $GP/bin/
echo "......GoEqs2TeX installed"
