#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    echo
    echo
    #go test -run MultiInterp01
    #go test -run Pow2
    #go test -run Pow3
    #go test -run InterpQuad01
    #go test -run InterpQuad02
    #go test -run InterpQuad03
    #go test -run InterpQuad04
    #go test -run InterpCubic01
    #go test -run InterpCubic02
    #go test -run InterpCubic03
    go test -run InterpCubic04
done
