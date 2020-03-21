GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
BINDATACMD = go-bindata
RSRCBIN = rsrc
APPIFYBIN = appify
BUILD_DIR=bin
BINARY_NAME_LINUX=$(BUILD_DIR)/discord-bot-gui_linux
BINARY_NAME_WIN32=$(BUILD_DIR)/discord-bot-gui_win32.exe
BINARY_NAME_WIN64=$(BUILD_DIR)/discord-bot-gui_win64.exe
BINARY_NAME_DARWIN=$(BUILD_DIR)/discord-bot-gui_darwin
OUT_NAME_DARWIN=$(BINARY_NAME_DARWIN).zip

.PHONY: all dist linux build test clean run win64 build-win64 win32 build-win32 darwin build-darwin dep dep-go dep-linux dep-win32 dep-win64 dep-darwin

all: linux

dist: linux darwin win32 win64

linux: dep dep-linux test build

build:
	$(BINDATACMD) ./ui/...
	GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) -v -mod=vendor -o $(BINARY_NAME_LINUX)

test:
	$(BINDATACMD) ./ui/...
	GO111MODULE=on $(GOFMT) ./...
	GO111MODULE=on $(GOVET) -v ./...

clean: 
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f ./bindata.go
	rm -f ./ui/astilectron.zip
	rm -f ui/electron-windows_386.zip
	rm -f ui/electron-windows_amd64.zip
	rm -f ui/electron-darwin_amd64.zip
	rm -f ui/electron-linux_amd64.zip
	rm -f discord-bot-gui.syso

run: dep dep-linux build
	./$(BINARY_NAME_LINUX)

win64: dep dep-win64 test build-win64

build-win64:
	$(BINDATACMD) ./ui/...
	$(RSRCBIN) -ico=discord-512.ico -arch=amd64 -o=discord-bot-gui.syso
	GO111MODULE=on GOOS=windows GOARCH=amd64 $(GOBUILD) -v -mod=vendor -o $(BINARY_NAME_WIN64) ./...

win32: dep dep-win32 test build-win32

build-win32:
	$(BINDATACMD) ./ui/...
	$(RSRCBIN) -ico=discord-512.ico -arch=386 -o=discord-bot-gui.syso
	GO111MODULE=on GOOS=windows GOARCH=386 $(GOBUILD) -v -mod=vendor -o $(BINARY_NAME_WIN32) ./...

darwin: dep dep-darwin test build-darwin

build-darwin:
	$(BINDATACMD) ./ui/...
	GO111MODULE=on GOOS=darwin GOARCH=amd64 $(GOBUILD) -v -mod=vendor -o $(BINARY_NAME_DARWIN) ./...
	cd $(BUILD_DIR); $(APPIFYBIN) -name "Discord Bot GUI" -icon ../discord-512.png ../$(BINARY_NAME_DARWIN)
	cd $(BUILD_DIR); zip -r ../$(OUT_NAME_DARWIN) "Discord Bot GUI.app"
	rm -rf $(BUILD_DIR)/"Discord Bot GUI.app"


dep: ui/astilectron.zip dep-go

ui/astilectron.zip:
	test -f $@ || wget -O ui/astilectron.zip https://github.com/asticode/astilectron/archive/v0.35.1.zip

dep-go:
	export GO111MODULE=on 
	cd; go get -u github.com/go-bindata/go-bindata/...
	cd; go get -u github.com/machinebox/appify/...
	cd; go get -u github.com/akavel/rsrc/...

dep-linux:
	rm -f ui/electron-windows_386.zip
	rm -f ui/electron-windows_amd64.zip
	rm -f ui/electron-darwin_amd64.zip
	test -f ui/electron-linux_amd64.zip || wget -O ui/electron-linux_amd64.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-linux-x64.zip

dep-win32:
	rm -f ui/electron-windows_amd64.zip
	rm -f ui/electron-darwin_amd64.zip
	rm -f ui/electron-linux_amd64.zip
	test -f ui/electron-windows_386.zip || wget -O ui/electron-windows_386.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-win32-ia32.zip

dep-win64:
	rm -f ui/electron-windows_386.zip
	rm -f ui/electron-darwin_amd64.zip
	rm -f ui/electron-linux_amd64.zip
	test -f ui/electron-windows_amd64.zip || wget -O ui/electron-windows_amd64.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-win32-x64.zip


dep-darwin:
	rm -f ui/electron-windows_386.zip
	rm -f ui/electron-windows_amd64.zip
	rm -f ui/electron-linux_amd64.zip
	test -f ui/electron-darwin_amd64.zip || wget -O ui/electron-darwin_amd64.zip https://github.com/electron/electron/releases/download/v4.0.1/electron-v4.0.1-darwin-x64.zip