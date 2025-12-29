all: goose-up start


goose-up:
	goose up
goose-down:
	goose down

tests-cover:
	go test -cover ./...

tests-v:
	go test -v ./...

tests-all:
	go test -v -cover ./...


clean:
	rm -rf build/*

assembly: clean
	mkdir -p build
	go build -o build/main ./cmd/main

start: assembly
	./build/main

run:
	go run ./cmd/main/main.go