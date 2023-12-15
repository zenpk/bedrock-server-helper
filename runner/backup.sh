#!/bin/bash

# $1 backup path
# $2 save data path

mkdir -p "$1"
cp -r "$1" "$2"
