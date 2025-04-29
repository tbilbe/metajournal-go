#!/bin/bash

mkdir -p build

GOOS=darwin  GOARCH=amd64 go build -o build/metajournal-darwin-amd64
GOOS=darwin  GOARCH=arm64 go build -o build/metajournal-darwin-arm64
GOOS=linux   GOARCH=amd64 go build -o build/metajournal-linux-amd64
GOOS=windows GOARCH=amd64 go build -o build/metajournal-windows-amd64.exe

echo "✅ Build complete! Files are in ./build"