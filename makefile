GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
BINDATACMD = go-bindata
BUILD_DIR=bin
BINARY_NAME_LINUX=$(BUILD_DIR)/discord-bot-gui_linux

.PHONY: all linux build test clean run dep

all: linux

linux: dep build

build:
	@$(BINDATACMD) ./ui/...
	@GO111MODULE=on GOOS=linux $(GOBUILD) -v -o $(BINARY_NAME_LINUX)

test:
	@$(BINDATACMD) ./ui/...
	@GO111MODULE=on $(GOFMT) ./...
	@GO111MODULE=on $(GOVET) -v ./...

clean: 
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f ./bindata.go

run: dep dep-linux build
	@./$(BINARY_NAME_LINUX)


dep:
	@cd; GO111MODULE=on go get -u github.com/go-bindata/go-bindata/...