#!/bin/bash

mkdir -p web/static/css/vendor
mkdir -p web/static/js/vendor

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/static/js/vendor/

