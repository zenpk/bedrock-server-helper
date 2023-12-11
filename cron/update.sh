#!/bin/bash
. ./vars.sh
wget "https://minecraft.azureedge.net/bin-linux/bedrock-server-$THIS_VERSION.zip"
unzip "bedrock-server-$THIS_VERSION.zip" -d "bedrock-server-$THIS_VERSION"
rm "bedrock-server-$THIS_VERSION.zip"
mkdir "bedrock-server-$THIS_VERSION/worlds"
mv "bedrock-server-$LAST_VERSION/worlds/newworld" "bedrock-server-$THIS_VERSION/worlds"
cp "start.sh" "bedrock-server-$THIS_VERSION"
cp "server.properties" "bedrock-server-$THIS_VERSION"
cp "allowlist.json" "bedrock-server-$THIS_VERSION"
chmod +x "bedrock-server-$THIS_VERSION/bedrock_server"
chmod +x "bedrock-server-$THIS_VERSION/start.sh"
