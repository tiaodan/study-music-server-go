#!/bin/bash
# Cross compile for Ubuntu and package

set -e

echo "=== Start building Ubuntu version ==="

rm -f music-server
rm -rf release-linux
rm -f music-server-linux.zip

echo "Compiling..."
GOOS=linux GOARCH=amd64 go build -o music-server .

echo "Packaging..."

mkdir -p release-linux
cp music-server release-linux/
cp config.yaml release-linux/
cp -r img release-linux/ 2>/dev/null || true
mkdir -p release-linux/logout

cd release-linux
zip -r ../music-server-linux.zip .
cd ..

rm -rf release-linux

echo ""
echo "=== Build complete ==="
ls -lh music-server-linux.zip