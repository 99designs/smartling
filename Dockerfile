FROM golang:1.12

COPY . /smartling
RUN cd /smartling && go mod download && go install ./...
RUN mkdir /work
WORKDIR /work

ENTRYPOINT ["smartling"]
