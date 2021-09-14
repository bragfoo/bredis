package main

import (
	"github.com/bragfoo/bredis/cmd"

	_ "net/http/pprof"
)

func main() {
	cmd.Execute()
}
