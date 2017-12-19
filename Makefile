CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-select
	cp -r cache src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r criteria src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r http src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r parser src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r query src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r reader src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r response src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r results src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r utils src/github.com/whosonfirst/go-whosonfirst-select/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"
	@GOPATH=$(GOPATH) go get -u "github.com/patrickmn/go-cache"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/sjson"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-sanitize"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cache/*.go
	go fmt cmd/*.go
	go fmt criteria/*.go
	go fmt http/*.go
	go fmt parser/*.go
	go fmt query/*.go
	go fmt reader/*.go
	go fmt response/*.go
	go fmt results/*.go
	go fmt utils/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-select cmd/wof-select.go
	@GOPATH=$(GOPATH) go build -o bin/wof-selectd cmd/wof-selectd.go
