#!/bin/bash

FILES="*.go"

echo
echo "monitoring:"
echo $FILES
echo
echo

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -test.run="EncDec02"
    #go test -test.run="Outputter05"
    #go test -test.run="Mylab01"
    go test -test.run="Deep04"
done
