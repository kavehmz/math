# math

[![Go Lang](http://kavehmz.github.io/static/gopher/gopher-front.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/kavehmz/math.svg?branch=master)](https://travis-ci.org/kavehmz/math)
[![Coverage Status](https://coveralls.io/repos/github/kavehmz/math/badge.svg?branch=master)](https://coveralls.io/github/kavehmz/math?branch=master)

Math is a service in Go which is expandable in a modular way. I started it mainly to serve some math formulas with standard caching, monitoring and other details.

Each new endpoint, like `/math/add` is a Go strcuture that impletemnts the following interface:

```go
type operation interface {
	Immutable() bool
	Calc(param []byte, end chan bool) ([]byte, error)
}
```

Adding a new endpoint will be one line code.

```go
m.formulas["/math/add"] = formulate("add", &add.Add{}, cacheTTL)
```

It was mainly an exercise for 
- Having a complete service that handles monitoring, timeout, caching, shutdown.
- instrumenting different calls automatically.
- dockerizing the whole setup so user __goes not need__ anything beside docker to run the service.
- having a modular design that adding/removing a new call will only take one line change. 
- caching the results automatically based on endpoint requirments.

# Build

This setup is fully dockerized to all you need is to do:

```
docker build -t imagename:version .
# or use make file examples
make build
```

# Example of how run locally

Makefile has some sample command to help in testing the service.


```bash
$ make run
2020/03/14 17:20:06 Listening at 8080

# Now open a new terminal and try some calls

$ make example-factorial 
curl -X POST -d '{"x":9}' 'http://127.0.0.1:8080/math/factorial'
362880

# Timeout is 1s by default. To test that we can try this
$ curl -X POST -d '{"x":900000000}' 'http://127.0.0.1:8080/math/factorial'
Operation timedout

$ make stop

# To see the metrics published by this service we can do
$ make check-metrics
curl 'http://127.0.0.1:8080/metrics'
factorial_accesses 1.0
factorial_cache_misses 1.0
factorial_cache_size 0.0
factorial_calculation_time_bucket{le="0.0005"} 1.0
...
```
