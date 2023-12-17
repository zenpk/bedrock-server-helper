#!/bin/bash

# $1 directory path
set -e
if [ -d "$1" ]; then
    echo "deleting directory: $1"
    rm -rf "$1"
else
    echo "directory does not exist: $1"
fi
echo "done"
