#!/bin/bash

# $1 output path
# $2 version
set -e
echo "downloading server"
wget --output-document "$1" "https://minecraft.azureedge.net/bin-linux/bedrock-server-$2.zip"
echo "done"
