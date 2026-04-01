@echo off
chcp 65001 >nul

echo === Start building Ubuntu version ===

del /f /q music-server 2>nul
rmdir /s /q release-linux 2>nul
del /f /q music-server-linux.zip 2>nul

echo Compiling...
set GOOS=linux
set GOARCH=amd64
go build -o music-server .

echo Packaging...

mkdir release-linux
mkdir release-linux\logout

copy music-server release-linux\
copy config.yaml release-linux\
xcopy /E /I /Y img release-linux\img >nul 2>&1

powershell -Command "Compress-Archive -Path release-linux\* -DestinationPath music-server-linux.zip -Force"

rmdir /s /q release-linux

echo.
echo === Build complete ===
dir music-server-linux.zip | findstr "music-server-linux.zip"