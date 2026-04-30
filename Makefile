.PHONY: test lint ci

test:
	go test -count 1 ./...

lint:
	go vet ./...

ci: lint test
