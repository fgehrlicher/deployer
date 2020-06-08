SHELL := /bin/bash

.PHONY: help # Generate list of targets with descriptions
help:
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20

.PHONY: build-all # Builds all available binaries
build-all: build-container build-windows-amd64 build-darvin-amd64 build-linux-amd64

.PHONY: build-container # Builds the docker version
build-container:
	docker build -t docker.pkg.github.com/fgehrlicher/deployer/deployer:latest .

.PHONY: build-windows-amd64 # Builds the windows amd64 exe
build-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -a -tags netgo -ldflags '-w' -o deployer-winAmd64.exe .

.PHONY: build-darvin-amd64 # Builds the osx i368 binary
build-darvin-386:
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 \
	go build -a -tags netgo -ldflags '-w' -o deployer-darvin386 .

.PHONY: build-darvin-amd64 # Builds the osx amd64 binary
build-darvin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
	go build -a -tags netgo -ldflags '-w' -o deployer-darvinAmd64 .

.PHONY: build-linux-amd64 # Builds the linux amd64 binary
build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -a -tags netgo -ldflags '-w' -o deployer-linuxAmd64 .
