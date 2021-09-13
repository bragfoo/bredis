package bredis

import (
	"errors"
	"testing"
)

func BRedisGetTest(t *testing.T, r BRedis) {
	if err := r.Set("b", "redis"); err != nil {
		t.Errorf("set data error: %s", err)
	}

	tests := []struct {
		name string
		key  string
		want string
		err  error
	}{
		{
			name: "get empty key",
			key:  "",
			want: "",
			err:  ErrorEmptyKey,
		},
		{
			name: "get no exist key",
			key:  "no_exist",
			want: "",
			err:  ErrorNotFound,
		},
		{
			name: "get exist key",
			key:  "b",
			want: "redis",
			err:  nil,
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
	if err := r.Set("b", "redis"); err != nil {
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
			name: "set empty key and val",
			args: args{
				key: "",
				val: "",
			},
			err: ErrorEmptyKey,
		},
		{
			name: "set empty key",
			args: args{
				key: "",
				val: "value",
			},
			err: ErrorEmptyKey,
		},
		{
			name: "set empty val",
			args: args{
				key: "key",
				val: "",
			},
			err: nil,
		},
		{
			name: "set exist key",
			args: args{
				key: "b",
				val: "new val",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Set(tt.args.key, tt.args.val); errors.Is(err, tt.err) {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.err)
			}
			val, err := r.Get(tt.args.key)
			if err != nil {
				t.Errorf("Set Success But Get() error = %v", err)
			}
			if val != tt.args.val {
				t.Errorf("Set Success But Get wrong val, val = %s, wantVal = %s", val, tt.args.val)
			}
		})
	}
}
