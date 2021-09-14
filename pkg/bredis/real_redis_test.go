package bredis

import (
	"testing"
)

func Test_readRedisBRedis_Get(t *testing.T) {
	r := NewRealRedisBRedis()
	BRedisGetTest(t, r)
}

func Test_readRedisBRedis_Set(t *testing.T) {
	r := NewRealRedisBRedis()
	BRedisSetTest(t, r)
}

func BenchmarkRealRedisBRedis_Get(b *testing.B) {
	r := NewRealRedisBRedis()
	GetBenchmark(b, r)
}

func BenchmarkRealRedisBRedis_Set(b *testing.B) {
	r := NewRealRedisBRedis()
	SetBenchmark(b, r)
}
