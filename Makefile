PROJECTNAME	:=	$(shell basename "$(PWD)")
GOBASE		:=	$(shell pwd)
GOBIN		:=	$(GOBASE)/bin
GOSOURCE	:=	$(wildcard *.go)
LDFLAGS		:=	-ldflags "-s -w"

.PHONY: rpc build clean cert

build: rpc
	go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOSOURCE)

rpc: rpc/service.proto
	protoc --go_out=plugins=grpc:. -I rpc/ service.proto

clean:
	go clean
	rm -f $(GOBIN)/$(PROJECTNAME)

cert:
	openssl req -x509 \
		-newkey rsa:4096 -keyout assets/regbox.key \
		-out assets/regbox.crt -nodes -days 365 \
		-config tools/req.conf -extensions v3_req

image-db:
	docker build -t postgres-twchd .
