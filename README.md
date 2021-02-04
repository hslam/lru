# lru
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/lru)](https://pkg.go.dev/github.com/hslam/lru)
[![Build Status](https://github.com/hslam/lru/workflows/build/badge.svg)](https://github.com/hslam/lru/actions)
[![codecov](https://codecov.io/gh/hslam/lru/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/lru)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/lru)](https://goreportcard.com/report/github.com/hslam/lru)
[![LICENSE](https://img.shields.io/github/license/hslam/lru.svg?style=flat-square)](https://github.com/hslam/lru/blob/master/LICENSE)

Package lru implements an LRU cache.

## Get started

### Install
```
go get github.com/hslam/lru
```
### Import
```
import "github.com/hslam/lru"
```
### Usage
#### Example
```go
package main

import (
	"fmt"
	"github.com/hslam/lru"
)

func main() {
	var capacity = 10
	var free lru.Free = func(key, value interface{}) {}
	l := lru.New(capacity, free)
	key := 1
	l.Set(key, "Hello world")
	l.Done(key)
	if v, ok := l.Get(key); ok {
		fmt.Println(v)
		l.Done(key)
	}
	l.Remove(key)
	l.Reset()
}
```

### Output
```
Hello world
```

### License
This package is licensed under a MIT license (Copyright (c) 2021 Meng Huang)


### Author
lru was written by Meng Huang.


