.PHONY: build
auth:
	go run -v ./cmd/auth

.PHONY: mocks
mocks:
	docker run -v "$PWD":/src -w /src vektra/mockery --all

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

#.DEFAULT_GOAL := auth
