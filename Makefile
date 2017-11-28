GOPATH ?= /go
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
PROJ := lusers

all: deps fmt test $(PROJ)

deps: $(DEPS)
	GOPATH=$(GOPATH) glide install

fmt:
	GOPATH=$(GOPATH) go fmt $(glide novendor)
	GOPATH=$(GOPATH) go tool vet *.go

test: deps
		GOPATH=$(GOPATH) go test -v $(shell glide novendor)

$(PROJ): deps 
	GOPATH=$(GOPATH) go build $(LDFLAGS) -o $@ -v $(glide novendor)
	touch $@ && chmod 755 $@

linux: deps
	GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build $(LDFLAGS) -o $(PROJ)-linux-amd64 -v $(glide novendor)
	touch $(PROJ)-linux-amd64 && chmod 755 $(PROJ)-linux-amd64

windows: deps
	GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build $(LDFLAGS) -o $(PROJ)-windows-amd64.exe -v $(glide novendor)
	touch $(PROJ)-windows-amd64.exe

darwin: deps
	GOOS=darwin GOARCH=amd64 GOPATH=$(GOPATH) go build -o $(PROJ)-darwin-amd64 -v $(glide novendor)
	touch $(PROJ)-darwin-amd64 && chmod 755 $(PROJ)-darwin-amd64

.PHONY: $(DEPS) clean

clean:
		rm -rf lusers lusers-win-amd64.exe lusers-linux-amd64.bin .glide vendor

