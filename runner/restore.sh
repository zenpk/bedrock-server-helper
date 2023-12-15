#!/bin/bash

# $1 backup path
# $2 save data path

rm -rf "$2"
cp -r "$1" "$2"
