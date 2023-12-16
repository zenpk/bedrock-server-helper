#!/bin/bash

# $1 filename
# $2 destination
# unzips a file and removes the zip
unzip -q "$1" -d "$2"
rm -f "$1"
