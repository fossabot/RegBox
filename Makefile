PROJECTNAME	:=	$(shell basename "$(PWD)")
GOBASE		:=	$(shell pwd)
GOBIN		:=	$(GOBASE)/bin
GOSOURCE	:=	$(wildcard *.go)
LDFLAGS		:=	-ldflags "-s -w"

.PHONY: rpc build clean

build: rpc
	go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOSOURCE)

rpc: rpc/service.proto
	protoc --go_out=plugins=grpc:. -I rpc/ service.proto

clean:
	go clean
	rm -f $(GOBIN)/$(PROJECTNAME)
