GOPATH ?= /go
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
PROJ := slackers

all: deps fmt test $(PROJ)

deps: $(DEPS)
	GOPATH=$(GOPATH) dep ensure

fmt:
	GOPATH=$(GOPATH) go fmt *.go
	GOPATH=$(GOPATH) go tool vet *.go

test: deps
		GOPATH=$(GOPATH) go test -v 

$(PROJ): deps 
	CGO_ENABLED=0 GOPATH=$(GOPATH) go build -a $(LDFLAGS) -o $@ -v *.go
	touch $@ && chmod 755 $@

linux: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build -a $(LDFLAGS) -o $(PROJ)-linux-amd64 -v *.go
	touch $(PROJ)-linux-amd64 && chmod 755 $(PROJ)-linux-amd64

windows: deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build -a $(LDFLAGS) -o $(PROJ)-windows-amd64.exe -v *.go
	touch $(PROJ)-windows-amd64.exe

darwin: deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GOPATH=$(GOPATH) go build -a $(LDFLAGS) -o $(PROJ)-darwin-amd64 -v *.go
	touch $(PROJ)-darwin-amd64 && chmod 755 $(PROJ)-darwin-amd64

.PHONY: $(DEPS) clean

clean:
		rm -rf slackers slackers-win-amd64.exe slackers-linux-amd64.bin vendor

