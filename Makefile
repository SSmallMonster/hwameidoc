.PHONY: build
build:
    GOOS=linux GOARCH=amd64 go build -o hwameidoc cmd/main.go
