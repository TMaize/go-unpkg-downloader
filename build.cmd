@echo off

echo clean....
if exist dist\go-unpkg-downloader del dist\go-unpkg-downloader
if exist dist\go-unpkg-downloader.exe del dist\go-unpkg-downloader.exe
if exist dist\go-unpkg-downloader_linux_amd64.tar.gz del dist\go-unpkg-downloader_linux_amd64.tar.gz
if exist dist\go-unpkg-downloader_windows_amd64.zip del dist\go-unpkg-downloader_windows_amd64.zip

echo build linux....
set GOOS=linux
set GOARCH=amd64
go build -o dist/go-unpkg-downloader main.go
if %errorlevel% neq 0 goto BUILDFAIL


cd dist
7z -sdel a go-unpkg-downloader_linux_amd64.tar go-unpkg-downloader
if %errorlevel% neq 0 goto BUILDFAIL
7z -sdel a go-unpkg-downloader_linux_amd64.tar.gz go-unpkg-downloader_linux_amd64.tar
if %errorlevel% neq 0 goto BUILDFAIL
cd ..

echo build windows....
set GOOS=windows
set GOARCH=amd64
go build -o dist/go-unpkg-downloader.exe main.go
if %errorlevel% neq 0 goto BUILDFAIL

cd dist
7z -sdel a go-unpkg-downloader_windows_amd64.zip go-unpkg-downloader.exe
if %errorlevel% neq 0 goto BUILDFAIL
cd ..

echo OK...................


goto END

:BUILDFAIL
echo.
echo FAIL................

:END