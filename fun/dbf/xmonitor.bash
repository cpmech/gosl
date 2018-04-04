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
    #go test -run Params03
    #go test -run Params04
    #go test -run Params05
    #go test -run Params06
    #go test -run Params07
    #go test -run Params08
    #go test -run Params09
    #go test -run Params10
    #go test -run Params11
    #go test -run Params12
    #go test -run Params13
    #go test -run Params14
    #go test -run Params15
    #go test -run Params16
    go test -run Params17
done
