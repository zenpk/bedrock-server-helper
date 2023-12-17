#!/bin/bash

# $1 backup path
# $2 save data path

set -e
echo "making backup folder"
mkdir -p "$1"
echo "backing up"
cp -r "$2" "$1"
echo "done"
