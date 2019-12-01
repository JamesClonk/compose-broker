.PHONY: run gin build test provision fetch-instance deprovision
SHELL := /bin/bash

all: run

run:
	source .env; source .env_*; go run main.go

gin:
	gin --all --immediate run main.go

build:
	rm -f compose-broker
	go build -o compose-broker

test:
	source .env && GOARCH=amd64 GOOS=linux go test -v ./...

provision:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc \
		-X PUT -H "Content-Type: application/json" \
		-d '{ "service_id":"8ff5d1c8-c6eb-4f04-928c-6a422e0ea330", "plan_id":"890d1ed6-0ff6-4c93-afc6-df753be6f1e3" }'

fetch-instance:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc \
		-X GET

deprovision:
	curl -v http://disco:dingo@localhost:9999/v2/service_instances/fe5556b9-8478-409b-ab2b-3c95ba06c5fc \
		-X DELETE
