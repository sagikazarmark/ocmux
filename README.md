# Gorilla Mux integration for Open Census

[![CircleCI](https://circleci.com/gh/sagikazarmark/ocmux.svg?style=svg)](https://circleci.com/gh/sagikazarmark/ocmux)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagikazarmark/ocmux?style=flat-square)](https://goreportcard.com/report/github.com/sagikazarmark/ocmux)
[![GolangCI](https://golangci.com/badges/github.com/sagikazarmark/ocmux.svg)](https://golangci.com/r/github.com/sagikazarmark/ocmux)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.12-61CFDD.svg?style=flat-square)](https://github.com/logur/logur)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/sagikazarmark/ocmux)


## Installation

```bash
$ go get github.com/sagikazarmark/ocmux
```


## Usage

```go
package main

import (
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/ocmux"
)

func main() {
	router := mux.NewRouter()
	router.Use(ocmux.Middleware())
}
```


## Attribution

Based on the work of [@basvanbeek](https://github.com/basvanbeek): https://github.com/basvanbeek/opencensus-gorilla_mux-example

Removed some Zipkin specific code, added a few features.


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
