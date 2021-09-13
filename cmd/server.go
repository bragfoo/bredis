package cmd

import (
	"errors"
	"fmt"
	"github.com/bragfoo/bredis/pkg/bredis"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

var (
	impl string
	port int

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "run bredis server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := server(); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&impl, "impl", "lock", `support three impl ["lock", "syncmap", "singlecore"]`)
	serverCmd.PersistentFlags().IntVar(&port, "port", 7379, `listen port`)
}

func server() error {
	s := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	var bRedis bredis.BRedis
	switch impl {
	case "lock":
		bRedis = bredis.NewLockBRedis()
	case "syncmap":
		bRedis = bredis.NewSyncMapBRedis()
	case "singlecore":
		bRedis = bredis.NewSingleCoreBRedis()
	default:
		bRedis = bredis.NewLockBRedis()
	}

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
							if errors.Is(err, bredis.ErrorNotFound) {
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
