PWD := $(shell pwd)
TMPDIR := bundles/tmp
PKG := github.com/ckeyer/bat
APP := bat

DEV_IMAGE := ckeyer/dev
PORT := 8001

NET := $(shell docker network inspect cknet > /dev/zero && echo "--net cknet --ip 172.16.1.7" || echo "")

try:
	echo "hello$(NET)"

build:
	docker run --rm \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) make local

local:
	go build -o $(TMPDIR)/$(APP) cli/main.go

dev:
	docker run --rm -it \
	 $(NET) \
	 --name bat-deving \
	 -e ADDR=":$(PORT)" \
	 -p $(PORT):$(PORT) \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) bash

run:
	go run cli/main.go -D 