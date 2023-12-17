#!/bin/bash

# mkdir for world, backup and server
set -e
echo "creating folders"
mkdir -p "$1"
mkdir -p "$2"
mkdir -p "$3"
echo "done"
