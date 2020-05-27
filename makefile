GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
BINDATACMD = go-bindata
RSRCCMD = rsrc
BUILD_DIR=bin
BINARY_NAME_LINUX=$(BUILD_DIR)/discord-bot-gui_linux
BINARY_NAME_WIN32=$(BUILD_DIR)/discord-bot-gui_win32.exe
BINARY_NAME_WIN64=$(BUILD_DIR)/discord-bot-gui_win64.exe
WIN64_CROSSCOMPILE=x86_64-w64-mingw32-gcc
WIN32_CROSSCOMPILE=i686-w64-mingw32-gcc


.PHONY: all dist linux build test clean run win build-win dep

all: linux

dist: linux win win32

linux: dep build

build:
	@$(BINDATACMD) ./ui/...
	@rm -f discord-bot-gui.syso
	@GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_NAME_LINUX)

test:
	@$(BINDATACMD) ./ui/...
	@GO111MODULE=on $(GOFMT) ./...
	@GO111MODULE=on $(GOVET) -v ./...

clean: 
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f ./bindata.go
	@rm -f discord-bot-gui.syso

run: dep build
	@./$(BINARY_NAME_LINUX)

win: dep test build-win

win32: dep test build-win32

build-win:
	@$(BINDATACMD) ./ui/...
	@$(RSRCCMD) -ico=discord-512.ico -arch=amd64 -o=discord-bot-gui.syso
	@GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WIN64_CROSSCOMPILE) $(GOBUILD) -v -ldflags='-H windowsgui' -o $(BINARY_NAME_WIN64) ./...

build-win32:
	@$(BINDATACMD) ./ui/...
	@$(RSRCCMD) -ico=discord-512.ico -arch=386 -o=discord-bot-gui.syso
	@GO111MODULE=on GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=$(WIN32_CROSSCOMPILE) $(GOBUILD) -v -ldflags='-H windowsgui' -o $(BINARY_NAME_WIN32) ./...


dep:
	@cd; GO111MODULE=on go get -u github.com/go-bindata/go-bindata/...
	@cd; GO111MODULE=on go get -u github.com/akavel/rsrc/...