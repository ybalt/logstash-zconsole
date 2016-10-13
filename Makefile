export GOPATH:=$(shell pwd)/vendor
export PATH:=$(PATH):$(GOPATH)/bin

BINARY=logstash-zconsole

.PHONY: deps clean bootstrap container print-%

$(BINARY): *.go
	go build -o $(BINARY)

deps:
	mkdir -p vendor
	go get github.com/pebbe/zmq4

clean:
	rm -f $(BINARY)
	go fmt *.go

bootstrap:
	go get github.com/tools/godep

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)

docker: linux
	docker build -t ybalt/logstash-zconsole .

print-%: ; @echo $*=$($*)