.PHONY: lint

all: lint

lint:
	golangci-lint run -v
