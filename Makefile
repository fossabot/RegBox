PROJECTNAME	:=	regbox
GOBIN		:=	$(shell pwd)/bin
LDFLAGS		:=	-ldflags "-s -w"

.PHONY: build clean mongo

build: a_main-packr.go pb/RegBoxService.pb.go
	go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) 

pb/RegBoxService.pb.go: pb/RegBoxService.proto
	protoc --go_out=plugins=grpc:pb/ -I pb/ RegBoxService.proto

assets/regbox.crt: assets/regbox.key
	@echo Generating selfsigned certificate...
	@openssl req -x509 -new -key $< \
		-out $@ -nodes -days 3650 \
		-config tools/req.conf -extensions v3_req

assets/regbox.key:
	@openssl genrsa -out $@ 4096

a_main-packr.go: assets/regbox.key assets/regbox.crt
	@echo Embedding static files...
	@packr

clean:
	rm -f $(GOBIN)/$(PROJECTNAME)
	
mongo:
	docker build -t mongo-regbox .
