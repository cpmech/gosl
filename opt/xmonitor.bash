#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    echo
    echo
    #go test -run Powell01
    go test -run Powell02
done
