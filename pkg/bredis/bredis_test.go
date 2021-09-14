package bredis

import (
	"errors"
	"strconv"
	"testing"
)

func BRedisGetTest(t *testing.T, r BRedis) {
	if err := PrepareSampleDataSet(r); err != nil {
		t.Errorf("set data error: %s", err)
	}

	tests := []struct {
		name string
		key  string
		want string
		err  error
	}{
		{
			name: "get_empty_key",
			key:  "",
			want: "",
			err:  ErrEmptyKey,
		},
		{
			name: "get_exist_key",
			key:  "a",
			want: "redisA",
			err:  nil,
		},
		{
			name: "get_no_exist_key",
			key:  "no_exist",
			want: "",
			err:  ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Get(tt.key)
			if !errors.Is(err, tt.err) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BRedisSetTest(t *testing.T, r BRedis) {
	if err := PrepareSampleDataSet(r); err != nil {
		t.Errorf("set data error: %s", err)
	}

	type args struct {
		key string
		val string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "set_empty_key_and_val",
			args: args{
				key: "",
				val: "",
			},
			err: ErrEmptyKey,
		},
		{
			name: "set_empty_key",
			args: args{
				key: "",
				val: "value",
			},
			err: ErrEmptyKey,
		},
		{
			name: "set_empty_val",
			args: args{
				key: "key",
				val: "",
			},
			err: nil,
		},
		{
			name: "set_new_key",
			args: args{
				key: "new_a",
				val: "newA",
			},
			err: nil,
		},
		{
			name: "set_exist_key",
			args: args{
				key: "b",
				val: "new val",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Set(tt.args.key, tt.args.val); !errors.Is(err, tt.err) {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.err)
			}
			if tt.err == nil {
				val, err := r.Get(tt.args.key)
				if err != nil {
					t.Errorf("Set Success But Get() error = %v", err)
				}
				if val != tt.args.val {
					t.Errorf("Set Success But Get wrong val, val = %s, wantVal = %s", val, tt.args.val)
				}
			}
		})
	}
}

func BRedisDeleteTest(t *testing.T, r BRedis) {
	if err := PrepareSampleDataSet(r); err != nil {
		t.Errorf("set data error: %s", err)
	}

	tests := []struct {
		name string
		key  string
		err  error
	}{
		{
			name: "delete_empty_key",
			key:  "",
			err:  ErrEmptyKey,
		},
		{
			name: "delete_empty_val",
			key:  "key",
			err:  nil,
		},
		{
			name: "delete_new_key",
			key:  "new_a",
			err:  nil,
		},
		{
			name: "delete_exist_key",
			key:  "b",
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Delete(tt.key); !errors.Is(err, tt.err) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.err)
			}
			if tt.err == nil {
				_, err := r.Get(tt.key)
				if errors.Is(err, ErrEmptyKey) {
					t.Errorf("Delete Success But Get() error = %v", err)
				}
			}
		})
	}
}

func PrepareSampleDataSet(r BRedis) error {
	dataset := map[string]string{
		"a": "redisA",
		"b": "redisB",
	}
	for k, v := range dataset {
		err := r.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrepareNKDataSet(r BRedis, size int) error {
	for i := 0; i < size*1000; i++ {
		key := strconv.Itoa(i)
		val := strconv.Itoa(i)
		err := r.Set(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBenchmark(b *testing.B, r BRedis) {
	err := PrepareNKDataSet(r, 10)
	if err != nil {
		b.Fatal("prepare dataset error", err)
	}
	b.StartTimer()
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		r.Get(strconv.Itoa(n))
	}
	b.StopTimer()
}

func SetBenchmark(b *testing.B, r BRedis) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		key := strconv.Itoa(n)
		val := strconv.Itoa(n)
		r.Set(key, val)
	}
}

func ParallelBenchmark(b *testing.B, r BRedis) {
	// run the Fib function b.N times
	b.SetParallelism(2)
	b.RunParallel(func(pb *testing.PB) {
		n := 0
		for pb.Next() {
			key := strconv.Itoa(n)
			val := strconv.Itoa(n)
			r.Set(key, val)
			r.Get(key)
			r.Delete(key)
			n++
		}
	})
}
