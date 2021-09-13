package bredis

import "errors"

var ErrorNotFound = errors.New("key not exists")
var ErrorEmptyKey = errors.New("empty key")

type BRedis interface {
	Get(key string) (string, error)
	Set(key string, val string) error
}
