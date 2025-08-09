#!/usr/bin/env bash

NAME="$1"

if [ -z "$1" ] 
then
    NAME="you"
fi

echo "One for $NAME, one for me."