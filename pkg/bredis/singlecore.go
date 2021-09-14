package bredis

import (
	"time"
)

type singleCoreBRedis struct {
	keys map[string]string
	in   chan inType
	out  chan outType
	done chan int
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

func (r *singleCoreBRedis) Close() {
	r.done <- 0
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
				r.set(c.key, c.params[0])
				r.out <- outType{}
			default:
				time.Sleep(10 * time.Microsecond)
			}
		case <-r.done:
			return
		}
	}
}

func (r *singleCoreBRedis) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	i := inType{
		cmd: "get",
		key: key,
	}
	r.in <- i
	o := <-r.out
	if o.err != nil {
		return "", o.err
	}
	return o.result[0], nil
}

func (r *singleCoreBRedis) Set(key string, val string) error {
	if key == "" {
		return ErrEmptyKey
	}
	i := inType{
		cmd:    "set",
		key:    key,
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
	return "",
		ErrNotFound
}

func (r *singleCoreBRedis) set(key string, val string) {
	if v, ok := r.keys[key]; ok {
		if v != val {
			r.keys[key] = val
		}
	} else {
		r.keys[key] = val
	}
}
