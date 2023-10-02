#!/usr/bin/env bash
# -*- coding: utf-8 -*-

set -e

echo '# STEP: setup env/get inputs'
SCRIPT_ROOT=$(dirname $0)
SOFTWARE_VERSION=$(grep 'var Version' -r $SCRIPT_ROOT/security-goat/cmd/version.go | awk -F'"' '{print $2}')

if [[ -n "$1" ]]; then
  SOFTWARE_VERSION=$1
fi
echo "# STEP: check found version: $SOFTWARE_VERSION"


echo '# STEP: build binaries'
make

echo '# STEP: create release using github cli'
gh release create "v$SOFTWARE_VERSION" $SCRIPT_ROOT/build/* --generate-notes
