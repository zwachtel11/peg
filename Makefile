GOCMD=go
GOBUILD=$(GOCMD) build -v 
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

PEG=bin/peg

PKG := 

all: format peg

clean:
	rm -rf ${PEG} 
peg:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${PEG} github.com/zwachtel11/peg 

.PHONY: vendor
vendor:
	go mod tidy

format:
	gofmt -s -w cmd/ 