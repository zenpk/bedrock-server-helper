#!/bin/bash

# $1 server path
# $2 log path
# $3 properties
# $4 allow list

set -e
cd "$1" || exit
echo "writing properties"
echo "$3" > "./server.properties"
if [ -n "$4" ]; then
  echo "writing allow list"
  echo "$4" > "./allowlist.json"
fi
echo "starting server"
LD_LIBRARY_PATH=. ./bedrock_server >> "$2" 2>&1
