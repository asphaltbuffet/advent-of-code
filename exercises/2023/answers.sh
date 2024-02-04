#!/usr/bin/zsh

fd info -X \
    jq -r '.day as $day | .data.answers | to_entries | .[] | "\($day)." + (if .key == "a" then "1" else "2" end) + "=\"" + .value + "\""'| \
    sort -n
