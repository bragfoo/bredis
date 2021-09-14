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
	if key == "" {
		return "", ErrEmptyKey
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if v, ok := r.keys[key]; ok {
		return v, nil
	}
	return "", ErrNotFound
}

func (r *lockBRedis) Set(key string, val string) error {
	if key == "" {
		return ErrEmptyKey
	}
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

func (r *lockBRedis) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.keys, key)
	return nil
}
