#!/bin/bash

# $1 old save data path
# $2 new server path
# $3 properties
# $4 allow list
# $5 world name

set -e
echo "debug1"
if [ -d "$2/worlds/$5" ]; then
  echo "removing old save data"
  rm -rf "$2/worlds/$5"
fi
echo "moving current save data to the new server"
cp -r "$1" "$2/worlds/$5"
echo "writing properties"
echo "$3" > "$2/server.properties"
if [ -n "$4" ]; then
  echo "writing allow list"
  echo "$4" > "$2/allowlist.json"
fi
echo "done"
