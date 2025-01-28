#!/bin/bash

source ./scripts/download_bbook.sh

mkdir -p .dist/web/

cd docs/
../bin/bbook build
cp -r .book/* ../.dist/web/
cd ../..
