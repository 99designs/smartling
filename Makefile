export GO111MODULE=on
VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w

release: bin/smartling-linux-amd64 bin/smartling-darwin-amd64 bin/smartling-linux-arm64 bin/smartling-darwin-arm64 bin/smartling-windows-386.exe

bin/smartling-linux-amd64:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

bin/smartling-darwin-amd64:
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

bin/smartling-windows-386.exe:
	@mkdir -p bin
	GOOS=windows GOARCH=386 go build -o $@ -ldflags="$(FLAGS)" .

bin/smartling-linux-arm64:
	@mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -o $@ -ldflags="$(FLAGS)" .

bin/smartling-darwin-arm64:
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -o $@ -ldflags="$(FLAGS)" .

clean:
	rm -f bin/*
