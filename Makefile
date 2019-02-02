.PHONY: $(shell ls -d *)

default:
	@echo "Usage: make server/examples/stop"

test:
	go test -v -cover -coverprofile=cover.out ./...

lint:
	command -v gometalinter.v2 || (go get -u gopkg.in/alecthomas/gometalinter.v2 && ${GOPATH}/bin/gometalinter.v2 --install)
	${GOPATH}/bin/gometalinter.v2 --vendor --deadline=180s --cyclo-over=15 --disable errcheck ./...

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
