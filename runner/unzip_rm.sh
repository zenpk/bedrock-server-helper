#!/bin/bash

# $1 filename
# $2 destination
# unzips a file and removes the zip
echo "unzipping"
unzip -q "$1" -d "$2"
echo "clearing"
rm -f "$1"
echo "done"
