#!/bin/bash

# $1 old save data path
# $2 new server path
# $3 properties
# $4 allow list

cp -r "$1" "$2/worlds/"
echo "$3" > "$2/server.properties"
if [ -n "$4" ]; then
  echo "$4" > "$2/allowlist.json"
fi
