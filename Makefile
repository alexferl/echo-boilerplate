.PHONY: dev run test cover fmt mock openapi-lint pre-commit docker-build docker-run

.DEFAULT: help
help:
	@echo "make dev"
	@echo "	setup development environment"
	@echo "make run"
	@echo "	run app"
	@echo "make test"
	@echo "	run go test"
	@echo "make cover"
	@echo "	run go test with -cover"
	@echo "make cover-html"
	@echo "	run go test with -cover and show HTML"
	@echo "make tidy"
	@echo "	run go mod tidy"
	@echo "make fmt"
	@echo "	run gofumpt"
	@echo "make mock"
	@echo "	run mockery"
	@echo "make openapi-lint"
	@echo "	lint openapi spec"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"
	@echo "make docker-build"
	@echo "	build docker image"
	@echo "make docker-run"
	@echo "	run docker image"

check-gofumpt:
ifeq (, $(shell which gofumpt))
	$(error "gofumpt not in $(PATH), gofumpt (https://pkg.go.dev/mvdan.cc/gofumpt) is required")
endif

check-pre-commit:
ifeq (, $(shell which pre-commit))
	$(error "pre-commit not in $(PATH), pre-commit (https://pre-commit.com) is required")
endif

check-redocly:
ifeq (, $(shell which redocly))
	$(error "redocly not in $(PATH), redocly (https://redocly.com/docs/cli/installation/) is required")
endif

dev: check-pre-commit
ifeq (,$(wildcard ./private-key.pem))
	@echo "No private key file, generating one..."
	openssl genrsa -out private-key.pem 4096
endif
	pre-commit install

run:
	go build -o server-bin ./cmd/server && ./server-bin

build:
	go build -o server-bin ./cmd/server

test:
	go test -v ./...

cover:
	go test -cover -v ./...

cover-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

tidy:
	go mod tidy

fmt: check-gofumpt
	gofumpt -l -w .

mock:
	mockery

openapi-lint: check-redocly
	redocly lint openapi/openapi.yaml

pre-commit: check-pre-commit
	pre-commit

docker-build:
	docker build -t echo-boilerplate .

docker-run:
	docker run -p 1323:1323 --rm echo-boilerplate --http-bind-address 0.0.0.0
