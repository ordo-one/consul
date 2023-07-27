#!/bin/bash

set -xeuo pipefail

package_dir=$1

#export GODEBUG=gocachehash=1 
go test -tags="${GOTAGS:-}" -c -o "$package_dir/test.bin" "$package_dir" 2>&1 # | tee "$package_dir/test.bin.buildlog" | grep -qv '^HASH'
# grep -q '^HASH' "$package_dir/test.bin.buildlog" > "$package_dir/test.bin.hashlog"