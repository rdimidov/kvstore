CLI_APP_NAME=kvstore-cli

build-cli:
	go build -o ${CLI_APP_NAME} cmd/cli/main.go

run-cli: build-cli
	./${CLI_APP_NAME} $(ARGS)

run-unit-test:
	go test ./internal/...

run-test-coverage:
	go test ./... -coverprofile=coverage.out