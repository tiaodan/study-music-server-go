#!/bin/bash
# Cross compile for Ubuntu only (no package)

set -e

echo "=== Start building Ubuntu version ==="

rm -f music-server

echo "Compiling..."
GOOS=linux GOARCH=amd64 go build -o music-server .

echo ""
echo "=== Build complete ==="
ls -lh music-server