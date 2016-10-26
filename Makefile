PWD := $(shell pwd)
PKG := github.com/ckeyer/bat
APP := bat

DEV_IMAGE := ckeyer/dev
PORT := 8001

OS := $(shell go env GOOS)-$(shell go env GOARCH)
VERSION := $(shell cat VERSION.txt)

default:
	echo "hello$(NET)"

release: clean local
	cd bundles && tar zcf $(APP)$(VERSION).$(OS).tgz $(APP)

local:
	go build -v -o bundles/$(APP) cli/main.go
	echo "build Successful"

clean:
	-rm -rf bundles

run:
	go run cli/main.go

build:
	docker run --rm \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) make local

NET := $(shell docker network inspect cknet > /dev/zero && echo "--net cknet --ip 172.16.1.7" || echo "")
dev:
	docker run --rm -it \
	 $(NET) \
	 --name bat-deving \
	 -e ADDR=":$(PORT)" \
	 -p $(PORT):$(PORT) \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) bash
