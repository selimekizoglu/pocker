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
| `consul`          | The location of the consul instance (IP/FQDN with port) to query. Defaults to localhost:8500
| `service*`         | Name of the service registered in consul.
| `endpoint`        | The endpoint of the service to be poked. Defaults to "/"

\* = Required parameter

### Examples
The CLI interface supports all of the options detailed above.

```shell
$ pocker -consul 127.0.0.1:8500   \
         -service healthy-service \
         -endpoint /healthcheck
```
