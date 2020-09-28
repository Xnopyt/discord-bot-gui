GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
BINDATACMD = go-bindata
RSRCCMD = rsrc
BUILD_DIR=bin
APPIFYCMD = appify
BINARY_NAME_LINUX=$(BUILD_DIR)/discord-bot-gui_linux
BINARY_NAME_WIN32=$(BUILD_DIR)/discord-bot-gui_win32.exe
BINARY_NAME_WIN64=$(BUILD_DIR)/discord-bot-gui_win64.exe
BINARY_NAME_DARWIN=$(BUILD_DIR)/discord-bot-gui_darwin
OUT_NAME_DARWIN=$(BINARY_NAME_DARWIN).zip
WIN64_CROSSCOMPILE=x86_64-w64-mingw32-gcc
WIN32_CROSSCOMPILE=i686-w64-mingw32-gcc


.PHONY: all dist linux build test clean run win build-win dep x-win x-build-win darwin dep-darwin build-darwin dep-win

all: linux

dist: linux x-win x-win32

linux: dep build

build:
	@echo 'Packing static assets...'
	@$(BINDATACMD) ./ui/...
	@rm -f discord-bot-gui.syso
	@echo 'Building Linux 64 bit binary...'
	@GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_NAME_LINUX)

test:
	@echo 'Packing static assets...'
	@$(BINDATACMD) ./ui/...
	@echo 'Running Tests...'
	@GO111MODULE=on $(GOFMT) ./...
	@GO111MODULE=on $(GOVET) -v ./...

clean: 
	@echo 'Removing build files...'
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f ./bindata.go
	@rm -f discord-bot-gui.syso

run: dep build
	@echo 'Running build...'
	@./$(BINARY_NAME_LINUX)

x-win: dep dep-win test x-build-win

x-win32: dep dep-win test x-build-win32

x-build-win:
	@echo 'Packing static assets...'
	@$(BINDATACMD) ./ui/...
	@echo 'Generating windows icon information...'
	@$(RSRCCMD) -ico=discord-512.ico -arch=amd64 -o=discord-bot-gui.syso
	@echo 'Cross Compiling Windows 64 bit binary...'
	@GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WIN64_CROSSCOMPILE) $(GOBUILD) -v -ldflags='-H windowsgui' -o $(BINARY_NAME_WIN64) ./...

x-build-win32:
	@echo 'Packing static assets...'
	@$(BINDATACMD) ./ui/...
	@echo 'Generating windows icon information...'
	@$(RSRCCMD) -ico=discord-512.ico -arch=386 -o=discord-bot-gui.syso
	@echo 'Cross Compiling Windows 32 bit binary...'
	@GO111MODULE=on GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=$(WIN32_CROSSCOMPILE) $(GOBUILD) -v -ldflags='-H windowsgui' -o $(BINARY_NAME_WIN32) ./...

darwin: dep dep-darwin build-darwin

build-darwin:
	@echo 'Packing static assets...'
	@$(BINDATACMD) ./ui/...
	@echo 'Building MacOS 64 bit binary...'
	@GO111MODULE=on GOOS=darwin GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_NAME_DARWIN) ./...
	@echo 'Packing MacOS application...'
	@cd $(BUILD_DIR); $(APPIFYCMD) -name "Discord Bot GUI" -icon ../discord-512.png -author "Xnopyt" ../$(BINARY_NAME_DARWIN)
	@cd $(BUILD_DIR); zip -r ../$(OUT_NAME_DARWIN) "Discord Bot GUI.app"
	@rm -rf $(BUILD_DIR)/"Discord Bot GUI.app"

dep-darwin:
	@echo 'Downloading MacOS build dependencies...'
	@cd; GO111MODULE=on go get -u github.com/machinebox/appify/...

dep-win:
	@echo 'Downloading Windows build dependencies...'
	@cd; GO111MODULE=on go get -u github.com/akavel/rsrc/...

dep:
	@echo 'Downloading build dependencies...'
	@cd; GO111MODULE=on go get -u github.com/go-bindata/go-bindata/...

win: dep dep-win test build-win

build-win:
	@echo 'Packing static assets...'
	@$(BINDATACMD) ./ui/...
	@echo 'Generating windows icon information...'
	@$(RSRCCMD) -ico=discord-512.ico -arch=amd64 -o=discord-bot-gui.syso
	@echo 'Cross Compiling Windows 64 bit binary...'
	@GO111MODULE=on GOOS=windows GOARCH=amd64 $(GOBUILD) -v -ldflags='-H windowsgui' -o $(BINARY_NAME_WIN64) ./...