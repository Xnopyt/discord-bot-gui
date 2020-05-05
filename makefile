GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
BINDATACMD = go-bindata
RSRCCMD = rsrc
APPIFYCMD = appify
BUILD_DIR=bin
BINARY_NAME_LINUX=$(BUILD_DIR)/discord-bot-gui_linux
BINARY_NAME_WIN32=$(BUILD_DIR)/discord-bot-gui_win32.exe
BINARY_NAME_WIN64=$(BUILD_DIR)/discord-bot-gui_win64.exe
BINARY_NAME_DARWIN=$(BUILD_DIR)/discord-bot-gui_darwin
OUT_NAME_DARWIN=$(BINARY_NAME_DARWIN).zip
WIN_CROSSCOMPILE=x86_64-w64-mingw32-gcc
OSX_CROSS_PATH=/opt/osxcross
OSX_SDK=MacOSX10.15.sdk
OSX_CROSSCOMPILE=$(OSX_CROSS_PATH)/bin/o64-clang

.PHONY: all dist linux build test clean run win build-win darwin build-darwin dep

all: linux

dist: linux darwin win

linux: dep test build

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


darwin: dep test build-darwin

build-darwin:
	@$(BINDATACMD) ./ui/...
	@GO111MODULE=on GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=$(OSX_CROSSCOMPILE) CGO_CFLAGS="-nostdinc -I${OSX_CROSS_PATH}/SDK/${OSX_SDK}/usr/include -I${OSX_CROSS_PATH}/SDK/${OSX_SDK}/System/Library/Frameworks/Kernel.framework/Versions/A/Headers -I${OSX_CROSS_PATH}/SDK/${OSX_SDK}/System/Library/Frameworks/CoreGraphics.framework/Headers -I${OSX_CROSS_PATH}/SDK/${OSX_SDK}/System/Library/Frameworks/Kernel.framework/Versions/A/Headers" $(GOBUILD) -v -o $(BINARY_NAME_DARWIN) ./...
	@cd $(BUILD_DIR); $(APPIFYCMD) -name "Discord Bot GUI" -icon ../discord-512.png ../$(BINARY_NAME_DARWIN)
	@cd $(BUILD_DIR); zip -r ../$(OUT_NAME_DARWIN) "Discord Bot GUI.app"
	@rm -rf $(BUILD_DIR)/"Discord Bot GUI.app"

dep:
	@cd; GO111MODULE=on go get -u github.com/go-bindata/go-bindata/...
	@cd; GO111MODULE=on go get -u github.com/machinebox/appify/...
	@cd; GO111MODULE=on go get -u github.com/akavel/rsrc/...