package bredis

import (
	"fmt"
	"time"
)

type singleCoreBRedis struct {
	keys map[string]string
	in   chan inType
	out  chan outType
}

type inType struct {
	cmd    string
	key    string
	params []string
}

type outType struct {
	result []string
	err    error
}

func NewSingleCoreBRedis() BRedis {
	br := &singleCoreBRedis{
		keys: make(map[string]string),
		in:   make(chan inType),
		out:  make(chan outType),
	}
	go br.do()
	return br
}

func (r *singleCoreBRedis) do() {
	for {
		select {
		case c := <-r.in:
			switch c.cmd {
			case "get":
				result, err := r.get(c.key)
				o := outType{result: []string{result}, err: err}
				r.out <- o
			case "set":
				err := r.set(c.key, c.params[0])
				o := outType{result: []string{}, err: err}
				r.out <- o
			default:
				time.Sleep(10 * time.Microsecond)
			}
		}
	}
}

func (r *singleCoreBRedis) Get(key string) (string, error) {
	i := inType{
		cmd: "get",
		key: "key",
	}
	r.in <- i
	o := <-r.out
	if o.err != nil {
		return "", o.err
	}
	return o.result[0], nil
}

func (r *singleCoreBRedis) Set(key string, val string) error {
	i := inType{
		cmd:    "set",
		key:    "key",
		params: []string{val},
	}
	r.in <- i
	o := <-r.out
	return o.err
}

func (r *singleCoreBRedis) get(key string) (string, error) {
	if v, ok := r.keys[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("not found")
}

func (r *singleCoreBRedis) set(key string, val string) error {
	if v, ok := r.keys[key]; ok {
		if v != val {
			r.keys[key] = val
		}
	} else {
		r.keys[key] = val
	}
	return nil
}
