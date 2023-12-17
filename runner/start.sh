#!/bin/bash

# $1 server path
# $2 log path
LD_LIBRARY_PATH=. "$1/bedrock_server" >> "$2" 2>&1
