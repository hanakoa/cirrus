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

.PHONY: docker-build
docker-build:
	docker build -t kevinmichaelchen/cirrus:latest .

.PHONY: docker-rebuild
docker-rebuild:
	docker build -t kevinmichaelchen/cirrus:latest . --no-cache

.PHONY: lint
lint:
	golint .

.PHONY: fmt
fmt:
	go fmt .

.PHONY: vet
vet:
	go tool vet .
