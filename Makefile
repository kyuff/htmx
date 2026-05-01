.PHONY: test test-coverage cover lint gen ci

test:
	go test -count 1 -race ./...

test-coverage:
	go test -coverprofile=coverage.txt -count 1 -race ./...

cover:
	go test -count 1 -race -cover ./...

lint:
	go vet ./...

gen:
	go generate ./...

ci: lint test
