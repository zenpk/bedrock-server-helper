#!/bin/bash

# $1 backup path
# $2 backup path with world name
# $3 save data path

set -e
echo "making backup folder"
mkdir -p "$1"
echo "backing up"
cp -r "$3" "$2"
echo "done"
