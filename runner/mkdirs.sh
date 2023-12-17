#!/bin/bash

# mkdir for world, backup and server. Then create log file
set -e
echo "creating folders"
mkdir -p "$1"
mkdir -p "$2"
mkdir -p "$3"
echo "creating log file"
touch "$4"
echo "done"
