package bredis

import (
	"fmt"
	"sync"
)

type syncMapBRedis struct {
	keys *sync.Map
}

func NewSyncMapBRedis() BRedis {
	return &syncMapBRedis{
		keys: &sync.Map{},
	}
}

func (r *syncMapBRedis) Get(key string) (string, error) {
	if v, ok := r.keys.Load(key); ok {
		return v.(string), nil
	}
	return "", fmt.Errorf("not found")
}

func (r *syncMapBRedis) Set(key string, val string) error {
	if v, ok := r.keys.Load(key); ok {
		if v.(string) != val {
			r.keys.Store(key, val)
		}
	} else {
		r.keys.Store(key, val)
	}
	return nil
}
