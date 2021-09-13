package bredis

import (
	"sync"
)

type lockBRedis struct {
	keys  map[string]string
	mutex *sync.RWMutex
}

func NewLockBRedis() BRedis {
	return &lockBRedis{
		keys:  make(map[string]string),
		mutex: &sync.RWMutex{},
	}
}

func (r *lockBRedis) Get(key string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if v, ok := r.keys[key]; ok {
		return v, nil
	}
	return "", ErrorNotFound
}

func (r *lockBRedis) Set(key string, val string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if v, ok := r.keys[key]; ok {
		if v != val {
			r.keys[key] = val
		}
	} else {
		r.keys[key] = val
	}
	return nil
}
