BINARY  := teral
BIN_DIR := bin

# On Windows, Go executables need the .exe suffix.
ifeq ($(OS),Windows_NT)
EXT := .exe
endif

.DEFAULT_GOAL := build
.PHONY: build run test fmt fmt-check vet clean

build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY)$(EXT) ./cmd/teral

run: build
	@./$(BIN_DIR)/$(BINARY)$(EXT)

fmt:
	@gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "files need formatting:"; gofmt -l .; exit 1)

vet:
	@go vet ./...

clean:
	@rm -rf $(BIN_DIR)

test:
	@go test ./...