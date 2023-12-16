#!/bin/bash

# $1 backup path
# $2 save data path

echo "making backup folder"
mkdir -p "$1"
echo "backing up"
cp -r "$1" "$2"
echo "done"
