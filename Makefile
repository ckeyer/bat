PWD := $(shell pwd)
PKG := github.com/ckeyer/bat
APP := bat

DEV_IMAGE := ckeyer/dev
PORT := 8001

OS := $(shell go env GOOS)-$(shell go env GOARCH)
VERSION := $(shell cat VERSION.txt)

REDIS_ADDR := 172.16.1.10:6379
REDIS_AUTH := rh4WlOiNx3zBhiiU11wU

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
	REDIS_ADDR=$(REDIS_ADDR) \
	REDIS_AUTH=$(REDIS_AUTH) \
	DEBUG=true \
	go run cli/main.go

test: 
	docker run --rm \
	 --name bat-testing \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) make unit-test

unit-test:
	go test $$(go list ./... |grep -v "vendor")

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
