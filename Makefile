BINARY_NAME=vacme

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test ./...

coverage:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out
	rm cover.out

all: build
