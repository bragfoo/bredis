package bredis

import (
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
	if key == "" {
		return "", ErrEmptyKey
	}
	if v, ok := r.keys.Load(key); ok {
		return v.(string), nil
	}
	return "", ErrNotFound
}

func (r *syncMapBRedis) Set(key string, val string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if v, ok := r.keys.Load(key); ok {
		if v.(string) != val {
			r.keys.Store(key, val)
		}
	} else {
		r.keys.Store(key, val)
	}
	return nil
}

func (r *syncMapBRedis) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	r.keys.Delete(key)
	return nil
}
