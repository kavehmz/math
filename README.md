# math

[![Go Lang](http://kavehmz.github.io/static/gopher/gopher-front.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/kavehmz/math.svg?branch=master)](https://travis-ci.org/kavehmz/math)
[![Coverage Status](https://coveralls.io/repos/github/kavehmz/math/badge.svg?branch=master)](https://coveralls.io/github/kavehmz/math?branch=master)

math formula server.

This is a doodling for one way of serving apps in Go in a modular way.

It was mainly an exercise for 
- using docker for different actions.
- instrumenting different calls automatically.
- having a modular design that adding/removing a new call will only take one line change. 
```go
m.formulas["/math/add"] = formulate("add", &add.Add{}, ttl)
```
- caching the results from serve side.
- I intentionally did not use context.
- two stage docker build for a Go code.


# Example

Example of run

```bash
$ make help
docker run --rm kaveh-math --help
usage: math [<flags>]

Flags:
  --help          Show context-sensitive help (also try --help-long and
                  --help-man).
  --port=8080     port
  --timeout=1000  timeout in ms
  --ttl=3600      cache ttl in seconds

$ make serve 
docker run -d --rm --name kaveh-math -p 8080:8080 kaveh-math

$ make example-factorial 
curl -X POST -d '{"x":9}' 'http://127.0.0.1:8080/math/factorial'
362880

# REQ TIMEOUT is 1s by default and can be adjusted
$ curl -X POST -d '{"x":900000000}' 'http://127.0.0.1:8080/math/factorial'
Operation timedout

$ make stop
docker stop kaveh-math
kaveh-math

$ curl 'http://127.0.0.1:8080/metrics'

factorial_accesses 1.0
factorial_cache_misses 1.0
factorial_cache_size 0.0
factorial_calculation_time_bucket{le="0.0005"} 1.0
...
```
