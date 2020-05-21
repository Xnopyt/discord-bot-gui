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
WIN_CROSSCOMPILE=x86_64-w64-mingw32-gcc

.PHONY: all dist linux build test clean run win build-win dep

all: linux

dist: linux win

linux: dep build

build:
	@$(BINDATACMD) ./ui/...
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

run: dep dep-linux build
	@./$(BINARY_NAME_LINUX)

win: dep test build-win

build-win:
	@$(BINDATACMD) ./ui/...
	@$(RSRCCMD) -ico=discord-512.ico -arch=amd64 -o=discord-bot-gui.syso
	@GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WIN_CROSSCOMPILE) $(GOBUILD) -v -ldflags='-H windowsgui' -o $(BINARY_NAME_WIN64) ./...
	@rm -rf vendor/github.com/zserge/webview/dll


dep:
	@cd; GO111MODULE=on go get -u github.com/go-bindata/go-bindata/...
	@cd; GO111MODULE=on go get -u github.com/akavel/rsrc/...