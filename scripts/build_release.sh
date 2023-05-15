#!/bin/bash

mkdir -p .dist/x86_64_unknown-linux/
mkdir -p .dist/x86_64_pc-windows/
mkdir -p .dist/x86_64_apple-darwin/

go build -o .dist/x86_64_unknown-linux/docbook
GOOS=windows GOARCH=amd64 go build -o .dist/x86_64_pc-windows/docbook.exe
GOOS=darwin GOARCH=amd64 go build -o .dist/x86_64_apple-darwin/docbook

cd .dist/x86_64_unknown-linux/
tar -czvf docbook-x86_64_unknown-linux.tar.gz docbook
cd ../x86_64_pc-windows
tar -czvf docbook-x86_64_pc-windows.tar.gz docbook.exe
cd ../x86_64_apple-darwin
tar -czvf docbook-x86_64_apple-darwin.tar.gz docbook

cd ../..

