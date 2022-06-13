OUT=./bin/kindli

.PHONY: build run

run:
	go run $(OUT)

build:
	go build -o $(OUT) $(CFLAGS) .
