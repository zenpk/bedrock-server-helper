#!/bin/bash

# $1 backup path
# $2 save data path

set -e
echo "removing current save data"
rm -rf "$2"
echo "restoring backup"
cp -r "$1" "$2"
echo "done"
