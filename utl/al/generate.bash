#!/bin/bash

FNKEYS="
linkedlist \
queue \
recquicksort \
"

FNKEYS="queue"

TYPES[1]="int"
NAMES[1]="Int"

TYPES[2]="float64"
NAMES[2]="Float64"

TYPES[3]="string"
NAMES[3]="String"

for key in $FNKEYS; do
    fnt="template_$key.got"
    for i in 1 2 3; do
        TYPE=${TYPES[i]}
        NAME=${NAMES[i]}
        fng="$key"_"$TYPE".go
        sed -e 's,DATATYPE,'"$TYPE"',g' \
            -e 's,TYPENAME,'"$NAME"',g' $fnt > $fng
        gofmt -s -w $fng
    done
done
