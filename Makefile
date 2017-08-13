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

docker-build:
	docker build . -t stor.highloadcup.ru/travels/tapir_winner

docker-run: docker-build
	docker run -p 80:8000 -v $(shell pwd)/data:/tmp/data stor.highloadcup.ru/travels/tapir_winner

docker-upload: docker-build
	docker push stor.highloadcup.ru/travels/tapir_winner

local-run: fast-build
	./bin/highloadcup --config=config.yaml
