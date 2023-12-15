#!/bin/bash

if [ -d "$1" ]; then
    echo "deleting backup: $1"
    rm -rf "$1"
else
    echo "backup does not exist: $1"
fi
