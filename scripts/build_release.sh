#!/bin/bash

mkdir -p .dist/x86_64_unknown-linux/
mkdir -p .dist/x86_64_pc-windows/
mkdir -p .dist/x86_64_apple-darwin/

go build -o .dist/x86_64_unknown-linux/bfbook
GOOS=windows GOARCH=amd64 go build -o .dist/x86_64_pc-windows/bfbook.exe
GOOS=darwin GOARCH=amd64 go build -o .dist/x86_64_apple-darwin/bfbook

cd .dist/x86_64_unknown-linux/
tar -czvf bfbook-x86_64_unknown-linux.tar.gz bfbook
cd ../x86_64_pc-windows
tar -czvf bfbook-x86_64_pc-windows.tar.gz bfbook.exe
cd ../x86_64_apple-darwin
tar -czvf bfbook-x86_64_apple-darwin.tar.gz bfbook

cd ../..

