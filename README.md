Pocker [![Build Status](https://travis-ci.org/selimekizoglu/pocker.svg?branch=master)](https://travis-ci.org/selimekizoglu/pocker)
===============

Post-deployment testing consul services.


Installation
------------
```shell
$ go get github.com/selimekizoglu/pocker
```

Usage
-----
### Options
|       Option      | Description |
| ----------------- |------------ |
| `consul`          | The location of the consul instance (IP/FQDN with port) to query. Defaults to localhost:8500.
| `service`         | Name of the service registered in consul. Defaults to "".
| `expect`          | Number of expected service instances registered in consul. Defaults to 1.
| `endpoint`        | The endpoint of the service to be poked. Defaults to "/".
| `retry`           | Number of retries after a failing poke. Defaults to 0 (Poke once).
| `timeout`         | Timeout of each poke retry (in milliseconds). Defaults to 0.

### Examples
The CLI interface supports all of the options detailed above.

```shell
$ pocker -consul 127.0.0.1:8500   \
         -service healthy-service \
         -expect 5                \
         -endpoint /healthcheck   \
         -retry 3                 \
         -timeout 3000
```

```shell
$ docker run --rm selimekizoglu/pocker:latest -consul 127.0.0.1:8500   \
                                              -service healthy-service \
                                              -expect 5                \
                                              -endpoint /healthcheck   \
                                              -retry 3                 \
                                              -timeout 3000
```
