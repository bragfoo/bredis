package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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
	keys  map[string]string
	mutex *sync.RWMutex
}

func NewBRedis() *BRedis {
	return &BRedis{
		keys:  make(map[string]string),
		mutex: &sync.RWMutex{},
	}
}

func (r *BRedis) Get(key string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if v, ok := r.keys[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("not found")
}

func (r *BRedis) Set(key string, val string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if v, ok := r.keys[key]; ok {
		if v != val {
			r.keys[key] = val
		}
	} else {
		r.keys[key] = val
	}
	return nil
}
