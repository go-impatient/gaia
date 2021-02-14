# request
[![Build Status](https://img.shields.io/travis/mushroomsir/request.svg?style=flat-square)](https://travis-ci.org/mushroomsir/request)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/request.svg?style=flat-square)](https://coveralls.io/github/mushroomsir/request?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/mushroomsir/request/blob/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/mushroomsir/request)

A concise HTTP request client for Go.

## Installation

```sh
go get github.com/mushroomsir/request
```
## Feature

- Easy to use.
- Support gzip、deflate decompress.


## Usage

```go
package main

import (
	"github.com/mushroomsir/request"
)

func main() {
	url := "https://github.com"

	body := struct {
		Any interface{} `json:"Any"`
	}{
		Any: "",
	}

	result := struct {
		Any interface{} `json:"Any"`
	}{}

	request.Post(url).Body(&body).Result(&result).Do()
}

```

### Header

```go
request.Get("url").Set("Content-Type", "application/json").Result(&result).Do()
```
or

```go
header := http.Header{}
header.Set("Content-Type", "application/json")

request.Get("url").Header(header).Result(&result).Do()
```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/request/blob/master/LICENSE).
