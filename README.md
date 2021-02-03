# lru
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
	l := lru.New(10)
	key := 1
	l.Set(key, "Hello world")
	if v, ok := l.Get(key); ok {
		fmt.Println(v)
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


