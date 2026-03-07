APP_NAME=cli-sorter
LDFLAGS=-ldflags="-s -w"

build: build-linux build-windows build-macos

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(APP_NAME)-linux-amd64
build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/$(APP_NAME)-windows-amd64.exe
build-macos:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(APP_NAME)-macos-amd64