# Chaas - a bot of fun
# Copyright (c) 2015 Garrett Woodworth (https://github.com/gwoo)

.PHONY: bin/terminal bin/terminal-linux-amd64 bin/terminal-darwin-amd64 test bin/slack bin/slack-linux-amd64 bin/slack-darwin-amd64
VERSION := $(shell git describe --always --dirty --tags)
VERSION_FLAGS := -ldflags "-X main.Version=$(VERSION)"

bin/terminal: bin
	go build -v $(VERSION_FLAGS) -o $@ ./cmd/terminal

bin/terminal-linux-amd64: bin
	GOOS=linux GOARCH=amd64 go build -v $(VERSION_FLAGS) -o $@ ./cmd/terminal

bin/terminal-darwin-amd64: bin
	GOOS=darwin GOARCH=amd64 go build -v $(VERSION_FLAGS) -o $@ ./cmd/terminal

bin/slack: bin
	go build -v $(VERSION_FLAGS) -o $@ ./cmd/slack

bin/slack-linux-amd64: bin
	GOOS=linux GOARCH=amd64 go build -v $(VERSION_FLAGS) -o $@ ./cmd/slack

bin/slack-darwin-amd64: bin
	GOOS=darwin GOARCH=amd64 go build -v $(VERSION_FLAGS) -o $@ ./cmd/slack


test:
	go test -v ./...

bin:
	mkdir bin