package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tryCmd)
}

var tryCmd = &cobra.Command{
	Use:   "server",
	Short: "run bredis server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := server(); err != nil {
			return err
		}
		return nil
	},
}

func server() error {
	s := http.Server{
		Addr: ":8080",
	}
	bRedis := NewBRedis()
	http.HandleFunc(
		"/",
		func(reply http.ResponseWriter, resp *http.Request) {
			uri := resp.RequestURI
			fmt.Println(resp.Method, uri)
			paths := strings.Split(strings.TrimLeft(uri, "/"), "/")
			if len(paths) > 1 {
				switch paths[0] {
				case "get":
					if len(paths) == 2 {
						result, err := bRedis.Get(paths[1])
						if err != nil {
							if err.Error() == "not found" {
								reply.WriteHeader(404)
								reply.Write([]byte(err.Error()))
							} else {
								reply.WriteHeader(500)
								reply.Write([]byte(err.Error()))
							}
						} else {
							reply.WriteHeader(200)
							reply.Write([]byte(result + "\n"))
						}
					} else {
						reply.WriteHeader(400)
						reply.Write([]byte("invalid params"))
					}
				case "set":
					if len(paths) == 3 {
						err := bRedis.Set(paths[1], paths[2])
						if err != nil {
							reply.WriteHeader(500)
							reply.Write([]byte(err.Error()))
						} else {
							reply.WriteHeader(200)
							reply.Write([]byte("ok\n"))
						}
					} else {
						reply.WriteHeader(400)
						reply.Write([]byte("invalid params"))
					}
				}
			}
		},
	)

	http.HandleFunc(
		"/ping",
		func(reply http.ResponseWriter, resp *http.Request) {
			fmt.Println(resp.Method, resp.RequestURI)
			reply.WriteHeader(200)
			reply.Write([]byte("pong\n"))
		},
	)

	log.Fatal(s.ListenAndServe())
	return nil
}

type BRedis struct {
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

func NewBRedis() *BRedis {
	br := &BRedis{
		keys: make(map[string]string),
		in:   make(chan inType),
		out:  make(chan outType),
	}
	go br.do()
	return br
}

func (r *BRedis) do() {
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

func (r *BRedis) Get(key string) (string, error) {
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

func (r *BRedis) Set(key string, val string) error {
	i := inType{
		cmd:    "set",
		key:    "key",
		params: []string{val},
	}
	r.in <- i
	o := <-r.out
	return o.err
}

func (r *BRedis) get(key string) (string, error) {
	if v, ok := r.keys[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("not found")
}

func (r *BRedis) set(key string, val string) error {
	if v, ok := r.keys[key]; ok {
		if v != val {
			r.keys[key] = val
		}
	} else {
		r.keys[key] = val
	}
	return nil
}
