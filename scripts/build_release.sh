#!/bin/bash

mkdir -p .dist/x86_64_unknown-linux/
mkdir -p .dist/x86_64_pc-windows/
mkdir -p .dist/x86_64_apple-darwin/

go build -o .dist/x86_64_unknown-linux/bbook
GOOS=windows GOARCH=amd64 go build -o .dist/x86_64_pc-windows/bbook.exe
GOOS=darwin GOARCH=amd64 go build -o .dist/x86_64_apple-darwin/bbook

cd .dist/x86_64_unknown-linux/
tar -czvf bbook-x86_64_unknown-linux.tar.gz bbook
cd ../x86_64_pc-windows
tar -czvf bbook-x86_64_pc-windows.tar.gz bbook.exe
cd ../x86_64_apple-darwin
tar -czvf bbook-x86_64_apple-darwin.tar.gz bbook

cd ../..

