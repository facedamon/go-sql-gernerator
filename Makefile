# Go paramters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME= go-sql-generator
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
	$(GOCMD) env -w GO111MODULE=on
	$(GOCMD) env -w GOPROXY=https://goproxy.io,direct
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v
clean:
	$(GOCLEAN)
	rm -fr $(BINARY_NAME)
	rm -fr $(BINARY_UNIX)

# Cross compilation
build-linux:
	$(GOCMD) env -w GO111MODULE=on
	$(GOCMD) env -w GOPROXY=https://goproxy.io,direct
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
