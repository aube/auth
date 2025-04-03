.PHONY: build
auth:
	go run -v ./cmd/auth

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

#.DEFAULT_GOAL := auth
