#!/usr/bin/env bash
# -*- coding: utf-8 -*-

set -e

echo "🟢 STEP: setup env/get inputs"
SOFTWARE_VERSION=$(grep 'var Version' -r cmd/version.go | awk -F'"' '{print $2}')

echo '🟢 STEP: build binaries using version: $SOFTWARE_VERSION'
make release-build

echo '🟢 STEP: create release using GitHub CLI'
gh release create "v$SOFTWARE_VERSION" ../build/* --generate-notes

echo '🟢 STEP: commit and push version'
git add .
git commit -m "release: v$SOFTWARE_VERSION" --allow-empty
git push
