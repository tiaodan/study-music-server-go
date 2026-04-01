@echo off
chcp 65001 >nul

echo === Start building Ubuntu version ===

del /f /q music-server 2>nul

echo Compiling...
set GOOS=linux
set GOARCH=amd64
go build -o music-server .

echo.
echo === Build complete ===
dir music-server | findstr "music-server"