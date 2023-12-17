#!/bin/bash

# $1 server path
# $2 log path
cd "$1" || exit
LD_LIBRARY_PATH=. ./bedrock_server >> "$2" 2>&1
