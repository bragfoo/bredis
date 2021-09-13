package bredis

import (
	"testing"
)

func Test_lockBRedis_Get(t *testing.T) {
	r := NewLockBRedis()
	BRedisGetTest(t, r)
}

func Test_lockBRedis_Set(t *testing.T) {
	r := NewLockBRedis()
	BRedisSetTest(t, r)
}
