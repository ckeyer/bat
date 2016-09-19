PWD := $(shell pwd)
TMPDIR := bundles/tmp
PKG := github.com/ckeyer/bat
APP := bat

DEV_IMAGE := ckeyer/dev
PORT := 8001

build:
	docker run --rm \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) make local

local:
	go build -o $(TMPDIR)/$(APP) cli/main.go

dev:
	docker run --rm -it \
	 -e PORT=$(PORT) \
	 -p $(PORT):$(PORT) \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) bash
