#!/bin/bash
rm -rf ./dist
mkdir ./dist

cp ../client/src/* ./dist

export HOSTNAME=0.0.0.0:3000
export LOG_LEVEL=debug

go run src/main.go