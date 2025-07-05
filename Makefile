.PHONY: run
run:
	go run -v ./cmd/auth

.PHONY: test
test:
	go test -v -timeout 30s ./...

.PHONY: race
race:
	go test -v -race -timeout 30s ./...

.PHONY: cover
cover:
	go test -v -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.html
	rm coverage.out

.DEFAULT_GOAL := run
