package main

import (
	"sync"
)

var si_pool *sync.Pool

func init_pool() {
	si_pool = &sync.Pool{
		New: func() interface{} {
			return new(StateInfo)
		},
	}
}
