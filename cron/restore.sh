#!/bin/bash

# default: no input parameter, recover yesterday's backup
# ./recover.sh 2 => recover the day before yesterday's backup

if [ -n "$1" ]; then
  days="$1"
else
  days=0
fi
. ./vars.sh
target=$(date -d "$date_now - $days day" +%Y-%m-%d)
cd ./backups
if [ -d temp ]; then # delete temp folder
  rm -r temp
  echo "temp backup deleted"
else
  echo "temp backup does not exist"
fi
mkdir temp
cd ..
mv -v "./bedrock-server-$THIS_VERSION/worlds/newworld" "./backups/temp"
cp -r "./backups/$target/newworld" "./bedrock-server-$THIS_VERSION/worlds"
echo "recover $days finished"
