BINARY=url_shorter

build:
	go build -o $(BINARY) .

run:
	go run .

test:
	go test ./...

fmt:
	go fmt ./...


