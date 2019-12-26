.DEFAULT_GOAL := run
SHELL := /bin/bash
APP ?= $(shell basename $$(pwd))
COMMIT_SHA = $(shell git rev-parse --short HEAD)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs main.go with the race detector
run:
	source .env; source .env_*; go run -race main.go

.PHONY: gin
## gin: runs main.go via gin (hot reloading)
gin:
	gin --all --immediate run main.go

.PHONY: build
## build: builds the application
build: clean
	@echo "Building binary ..."
	go build -o ${APP}

.PHONY: clean
## clean: cleans up binary files
clean:
	@echo "Cleaning up ..."
	@go clean

.PHONY: test
## test: runs go test with the race detector
test:
	@source .env; GOARCH=amd64 GOOS=linux go test -v -race ./...

.PHONY: init
## init: sets up go modules
init:
	@echo "Setting up modules ..."
	@go mod init 2>/dev/null; go mod tidy && go mod vendor

.PHONY: provision-instance
## provision-instance: creates an example service instance
provision-instance:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc?accepts_incomplete=true \
		-X PUT -H "Content-Type: application/json" \
		-d '{ "service_id":"e27ea95a-3883-44f2-8ca4-01101f39d50c", "plan_id":"355ef4a4-08f5-4764-b4ed-8353812b6963" }'

.PHONY: poll-instance
## poll-instance: queries last_operation of example service instance
poll-instance:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc/last_operation \
		-X GET

.PHONY: fetch-instance
## fetch-instance: queries example service instance
fetch-instance:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc \
		-X GET

.PHONY: update-instance
## update-instance: updates example service instance
update-instance:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc?accepts_incomplete=true \
		-X PATCH -H "Content-Type: application/json" \
		-d '{ "service_id":"e27ea95a-3883-44f2-8ca4-01101f39d50c", "plan_id":"ae2bda53-fe15-4335-9422-774aae3e7e32" }'

.PHONY: deprovision-instance
## deprovision-instance: deletes example service instance
deprovision-instance:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc?accepts_incomplete=true \
		-X DELETE

.PHONY: create-binding
## create-binding: creates binding for example service instance
create-binding:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc/service_bindings/deadbeef \
		-X PUT

.PHONY: fetch-binding
## fetch-binding: queries binding for example service instance
fetch-binding:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc/service_bindings/deadbeef \
		-X GET

.PHONY: remove-binding
## remove-binding: removes binding for example service instance
remove-binding:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc/service_bindings/deadbeef \
		-X DELETE
