.PHONY: build
.SILENT:

build:
	go mod tidy && CGO_ENABLED=0 go build -o ./bin/server ./cmd/remote-server-control.go && ./bin/server

test:
	go test internal/server/* -v