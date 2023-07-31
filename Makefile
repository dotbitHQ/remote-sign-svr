# build file
GO_BUILD=go build -ldflags -s -v

svr: BIN_BINARY_NAME=remote-sign-svr
svr:
	GO111MODULE=on $(GO_BUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run ./$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

cli: BIN_BINARY_NAME=remote-sign-cli
cli:
	GO111MODULE=on $(GO_BUILD) -o $(BIN_BINARY_NAME) cmd/cli/main.go
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run ./$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

update:
	go mod tidy -compat=1.17


