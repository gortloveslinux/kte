BINARY_NAME=vacme

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test

all: build
