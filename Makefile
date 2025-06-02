# Makefile for the project

BINARY_MASTER=master-app
BINARY_SLAVE=slave-app
BINARY_CLI=cli-app

CMD_MASTER=./cmd/master
CMD_SLAVE=./cmd/slave
CMD_CLI=./cmd/cli

BIN_DIR=bin

.PHONY: all build run-master run-slave run-cli test fmt lint clean

all: build

build:
	go build -o $(BIN_DIR)/$(BINARY_MASTER) $(CMD_MASTER)
	go build -o $(BIN_DIR)/$(BINARY_SLAVE) $(CMD_SLAVE)
	go build -o $(BIN_DIR)/$(BINARY_CLI) $(CMD_CLI)

run-master:
	go run $(CMD_MASTER)

run-slave:
	go run $(CMD_SLAVE)

run-cli:
	go run $(CMD_CLI)

test:
	go test ./...

fmt:
	gofmt -w .

fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

clean:
	rm -rf $(BIN_DIR)/
