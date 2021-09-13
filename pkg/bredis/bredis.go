package bredis

type BRedis interface {
	Get(key string) (string, error)
	Set(key string, val string) error
}
