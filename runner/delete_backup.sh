#!/bin/bash

# $1 backup path
set -e
if [ -d "$1" ]; then
    echo "deleting backup: $1"
    rm -rf "$1"
else
    echo "backup does not exist: $1"
fi
echo "done"
