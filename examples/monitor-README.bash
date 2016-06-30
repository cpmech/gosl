#!/bin/bash

FILE="README.md"
GOCODE="README-to-html.go"

refresh(){
    CURRENT=`xdotool getwindowfocus`
    BROWSER=`xdotool search --name "Google Chrome"`
    go run $GOCODE
    xdotool windowactivate $BROWSER
    xdotool key "CTRL+R"
    xdotool windowactivate $CURRENT
}

while true; do
    inotifywait -q -e modify $FILE
    refresh
done
