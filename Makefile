.PHONY: build run

run:
	go run .

build:
	go build $(CFLAGS) .
