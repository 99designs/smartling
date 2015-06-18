# the current tag else the current git sha
VERSION := $(shell git tag --points-at=HEAD | grep . || git rev-parse --short HEAD)

GOBUILD_ARGS := -ldflags "-X main.Version $(VERSION)"
OS := $(shell go env GOOS)
ARCH := $(shell go env GOHOSTARCH)

# Release steps:
#  git tag vx.x.x
#  git push --tags
#  make clean release

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
