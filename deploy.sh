#!/usr/bin/env bash
set -euo pipefail

echo "Installing Go..."

GO_VERSION=1.22.0
curl -LO https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar -C /tmp -xzf go${GO_VERSION}.linux-amd64.tar.gz
export PATH="/tmp/go/bin:$PATH"

go version

echo "Building anna from local source..."
go build -o anna

echo "Running anna..."
./anna