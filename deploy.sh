#!/usr/bin/env bash
set -euo pipefail

echo "Building anna from local source..."

go build -o anna

echo "Running anna..."
./anna