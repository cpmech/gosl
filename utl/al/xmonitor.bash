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
    echo
    echo
    echo
    echo
    #go test -run "IntQueue01"
    #go test -run "Float64Queue01"
    #go test -run "StringQueue01"
    #go test -run "IntLinkedList"
    go test -run "Float64LinkedList"
    #go test -run "StringLinkedList"
done
