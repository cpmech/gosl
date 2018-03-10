#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="LinReg01"
    #go test -test.run="LogReg00"
    #go test -test.run="LogReg01"
    #go test -test.run="LogReg02"
    go test
done
