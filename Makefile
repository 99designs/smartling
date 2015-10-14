# the current tag else the current git sha
VERSION := $(shell git tag --points-at=HEAD | grep . || git rev-parse --short HEAD)

GOBUILD_ARGS := -ldflags "-X main.Version=$(VERSION)"
OS := $(shell go env GOOS)
ARCH := $(shell go env GOHOSTARCH)

# To create a new release:
#  $ git tag vx.x.x
#  $ git push --tags
#  $ make clean
#  $ make release     # this will create 2 binaries in ./bin - darwin and linux
#
#  Next, go to https://github.com/99designs/smartling/releases/new
#  - select the tag version you just created
#  - Attach the binaries from ./bin/*

release: bin/smartling-linux-amd64 bin/smartling
	gzip bin/*

bin/smartling-linux-amd64:
	@mkdir -p bin
	docker run -it -v $$GOPATH:/go library/golang go build $(GOBUILD_ARGS) -o /go/src/github.com/99designs/smartling/$@ github.com/99designs/smartling/cli/smartling

bin/smartling:
	@mkdir -p bin
	go build $(GOBUILD_ARGS) -o bin/smartling-$(OS)-$(ARCH) ./cli/smartling

clean:
	rm -f bin/*
