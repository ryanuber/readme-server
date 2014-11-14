GOPATH=$(shell /bin/bash -c "echo $${GOPATH%%:*}")

all: deps build

deps:
	@echo Installing deps...
	@go get -v ./...

build:
	@echo Building...
	@go build \
		-v -o $(GOPATH)/bin/go-readme
	@cp scripts/readme.sh $(GOPATH)/bin/readme

test:
	@echo Running tests...
	@go test -cover ./...
