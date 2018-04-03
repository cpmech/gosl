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
    #go test -run FactObjs01
    #go test -run FactObjs02
    #go test -run LineSearch01
    #go test -run Powell01
    #go test -run Powell02
    #go test -run Powell03
    #go test -run ConjGrad01
    #go test -run ConjGrad02
    #go test -run ConjGrad03
    #go test -run ConjGrad04
    #go test -run GradDesc01
    #go test -run GradDesc02
    go test -run GradDesc03
done
