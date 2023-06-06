.PHONY: build clean win linux mac mac2

DIST=go-unpkg-downloader

build: clean win linux mac mac2
clean:
	rm -rf dist
	mkdir -p dist
win:
	GOOS=windows GOARCH=amd64 go build -o ./dist/$(DIST).exe main.go
	cd dist && 7z a -sdel $(DIST)-win32-x64.zip $(DIST).exe
linux:
	GOOS=linux GOARCH=amd64 go build -o ./dist/$(DIST) main.go
	cd dist && 7z a -sdel $(DIST)-linux-x64.tar $(DIST)
	cd dist && 7z a -sdel $(DIST)-linux-x64.tar.gz $(DIST)-linux-x64.tar
mac:
	GOOS=darwin GOARCH=amd64 go build -o ./dist/$(DIST) main.go
	cd dist && 7z a -sdel $(DIST)-darwin-x64.tar $(DIST)
	cd dist && 7z a -sdel $(DIST)-darwin-x64.tar.gz $(DIST)-darwin-x64.tar
mac2:
	GOOS=darwin GOARCH=arm64 go build -o ./dist/$(DIST) main.go
	cd dist && 7z a -sdel $(DIST)-darwin-arm64.tar $(DIST)
	cd dist && 7z a -sdel $(DIST)-darwin-arm64.tar.gz $(DIST)-darwin-arm64.tar
