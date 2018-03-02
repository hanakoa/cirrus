.PHONY: all
all:
	@$(MAKE) grpc
	@$(MAKE) build
	@$(MAKE) run

.PHONY: grpc
grpc:
	protoc pb/heartbeat.proto --go_out=plugins=grpc:.

.PHONY: build
build:
	go build -o ./bin/cirrus .

.PHONY: run
run:
	./bin/cirrus

.PHONY: fmt
fmt:
	go fmt .

.PHONY: vet
vet:
	go tool vet .
