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

func Test_lockBRedis_Delete(t *testing.T) {
	r := NewLockBRedis()
	BRedisDeleteTest(t, r)
}

func BenchmarkLockBRedis_Get(b *testing.B) {
	r := NewLockBRedis()
	GetBenchmark(b, r)
}

func BenchmarkLockBRedis_Set(b *testing.B) {
	r := NewLockBRedis()
	SetBenchmark(b, r)
}

func BenchmarkLockBRedis_Parallel(b *testing.B) {
	r := NewLockBRedis()
	ParallelBenchmark(b, r)
}
