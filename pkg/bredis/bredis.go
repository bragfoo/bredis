package bredis

import "errors"

var (
	ErrNotFound = errors.New("key not exists")
	ErrEmptyKey = errors.New("empty key")
)

type BRedis interface {
	Get(key string) (string, error)
	Set(key string, val string) error
	Delete(key string) error
}
