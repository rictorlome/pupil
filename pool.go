package main

import (
	"sync"
)

var siPool *sync.Pool

func initPool() {
	siPool = &sync.Pool{
		New: func() interface{} {
			return new(StateInfo)
		},
	}
}
