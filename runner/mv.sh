#!/bin/bash

# $1 old path
# $2 new path

set -e
echo "moving files"
mv "$1" "$2"
echo "done"
