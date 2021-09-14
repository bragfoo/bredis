package bredis

import (
	"testing"
)

func Test_singleCoreBRedis_Get(t *testing.T) {
	r := NewSingleCoreBRedis()
	BRedisGetTest(t, r)
}

func Test_singleCoreBRedis_Set(t *testing.T) {
	r := NewSingleCoreBRedis()
	BRedisSetTest(t, r)
}

func Test_singleCoreBRedis_Delete(t *testing.T) {
	r := NewSingleCoreBRedis()
	BRedisDeleteTest(t, r)
}

func BenchmarkSingleCoreBRedis_Get(b *testing.B) {
	r := NewSingleCoreBRedis()
	GetBenchmark(b, r)
}

func BenchmarkSingleCoreBRedis_Set(b *testing.B) {
	r := NewSingleCoreBRedis()
	SetBenchmark(b, r)
}

func BenchmarkSingleCoreBRedis_Parallel(b *testing.B) {
	r := NewSingleCoreBRedis()
	ParallelBenchmark(b, r)
}
