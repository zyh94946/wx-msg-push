NAME := wx-msg-push
VERSION := `git describe --dirty`
COMMIT := `git rev-parse HEAD`
BUILDDATE=`date "+%Y-%m-%d"`
CFGPATH='config.toml'

BUILD_DIR := build
VAR_SETTING := -X main.v=$(VERSION) -X main.c=$(COMMIT) -X main.d=${BUILDDATE}
GOBUILD = CGO_ENABLED=0 $(GO_DIR)go build -ldflags="-s -w -buildid= $(VAR_SETTING)" -o $(BUILD_DIR)

.PHONY: build release

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD)

runs:
	go run -ldflags "-s -w $(VAR_SETTING)" main.go server -c ${CFGPATH}

%.zip: %
	@zip -du $(NAME)-$@ -j $(BUILD_DIR)/$</*
	@zip -du $(NAME)-$@ ${CFGPATH}
	@echo "<<< ---- $(NAME)-$@"

release: darwin-amd64.zip linux-amd64.zip freebsd-amd64.zip windows-amd64.zip

darwin-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=darwin $(GOBUILD)/$@

linux-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD)/$@

freebsd-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=freebsd $(GOBUILD)/$@

windows-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=windows $(GOBUILD)/$@

