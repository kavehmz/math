.PHONY: $(shell ls -d *)

default:
	@echo "Usage: make server/examples/stop"

test:
	go test -v -cover -coverprofile=cover.out ./...

lint:
	docker run --rm -v $${PWD}:/app -w /app golangci/golangci-lint golangci-lint run -v

build:
	docker build -t kaveh-math .

help: build
	docker run --rm kaveh-math --help

run: build
	docker run -ti --rm --name kaveh-math -p 8080:8080 kaveh-math

serve: build
	docker run -d --rm --name kaveh-math -p 8080:8080 kaveh-math

stop:
	docker stop kaveh-math

example-factorial:
	curl -X POST -d '{"x":9}' 'http://127.0.0.1:8080/math/factorial'

example-fibonacci:
	curl -X POST -d '{"x":9000}' 'http://127.0.0.1:8080/math/fibonacci'

check-metrics:
	curl 'http://127.0.0.1:8080/metrics'
