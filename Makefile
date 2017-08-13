fmt:
	go fmt ./...

test:
	go test ./...

deps:
	dep ensure

deps-update:
	dep ensure --update

fast-build:
	go build -o ./bin/highloadcup ./cmd/main.go

build-alpine: deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/highloadcup ./cmd/main.go
