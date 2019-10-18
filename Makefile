test:
	env "GORACE=halt_on_error=1" go test -v -benchtime 1ns -bench . -race ./...

test_coverage:
	env "GORACE=halt_on_error=1" go test -v -benchtime 1ns -bench . -covermode=count -coverprofile=coverage.out ./...

codecov_coverage: test_coverage
	curl -s https://raw.githubusercontent.com/codecov/codecov-bash/1044b7a243e0ea0c05ed43c2acd8b7bb7cef340c/codecov | bash -s -- -f coverage.out  -Z

setup_ci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	golangci-lint run
