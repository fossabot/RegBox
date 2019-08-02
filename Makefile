PROJECTNAME	:=	$(shell basename "$(PWD)")
GOBASE		:=	$(shell pwd)
GOBIN		:=	$(GOBASE)/bin
GOSOURCE	:=	$(filter-out $(wildcard *_generate.go),$(wildcard *.go))
LDFLAGS		:=	-ldflags "-s -w"

.PHONY: rpc build clean cert generate image

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

image:
	docker build -t postgres-twchd .

generate: assets
	go run assets_generate.go