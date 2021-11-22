.PHONY: clean
.ONESHELL:

build:
	rm -rf dist
	mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build -o ./dist/go-unpkg-downloader     *.go
	GOOS=windows GOARCH=amd64 go build -o ./dist/go-unpkg-downloader.exe *.go
	cd dist
	7z a       go-unpkg-downloader_linux_amd64.tar     go-unpkg-downloader
	7z a -sdel go-unpkg-downloader_linux_amd64.tar.gz  go-unpkg-downloader_linux_amd64.tar
	7z a       go-unpkg-downloader_windows_amd64.zip   go-unpkg-downloader.exe


