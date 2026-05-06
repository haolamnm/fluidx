.PHONY: build test run clean

BINARY = fluidx

build:
	go build -o $(BINARY) ./cmd/fluidx

test:
	go test ./...

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -f $(BINARY).exe
	go clean -cache
