package bredis

import (
	"testing"
)

func Test_syncMapBRedis_Get(t *testing.T) {
	r := NewSyncMapBRedis()
	BRedisGetTest(t, r)
}

func Test_syncMapBRedis_Set(t *testing.T) {
	r := NewSyncMapBRedis()
	BRedisSetTest(t, r)
}

func Test_syncMapBRedis_Delete(t *testing.T) {
	r := NewSyncMapBRedis()
	BRedisDeleteTest(t, r)
}

func BenchmarkSyncMapBRedis_Get(b *testing.B) {
	r := NewSyncMapBRedis()
	GetBenchmark(b, r)
}

func BenchmarkSyncMapBRedis_Set(b *testing.B) {
	r := NewSyncMapBRedis()
	SetBenchmark(b, r)
}
