VERSION := $(shell git tag --points-at=HEAD || git rev-parse --short HEAD)
GOBUILD_ARGS := -ldflags "-X main.Version $(VERSION)"
OS := $(shell go env GOOS)
ARCH := $(shell go env GOHOSTARCH)

release: bin/smartling-linux-amd64 bin/smartling

bin/smartling-linux-amd64:
	@mkdir -p bin
	docker run -it -v $$GOPATH:/go library/golang go build $(GOBUILD_ARGS) -o /go/src/github.com/99designs/smartling/$@ github.com/99designs/smartling/cli/smartling

bin/smartling:
	@mkdir -p bin
	go build $(GOBUILD_ARGS) -o bin/smartling-$(OS)-$(ARCH) ./cli/smartling

clean:
	rm -f bin/*
