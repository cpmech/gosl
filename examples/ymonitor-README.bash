#!/bin/bash

FILE="README.md"
GOCODE="README-to-html.go"

refresh(){
    CURRENT=`xdotool getwindowfocus`
    BROWSER=`xdotool search --name "Gosl Examples - Google Chrome"`
    echo "CURRENT = $CURRENT"
    echo "BROWSER = $BROWSER"
    go run $GOCODE
    xdotool windowactivate $BROWSER
    xdotool key "CTRL+R"
    xdotool windowactivate $CURRENT
}

while true; do
    inotifywait -q -e modify $FILE
    refresh
done
