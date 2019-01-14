package main

import (
	"fmt"
	"os"
	"voting/rest/handler"
	"voting/rest/router"

	flags "github.com/jessevdk/go-flags"
)

// Opts represents options for running the processor
type Opts struct {
	Port string `short:"P" long:"port" description:"Port for hosting the client" default:"9009"`
	Rest string `long:"rest" description:"REST-API endpoint to connect to" default:"http://localhost:8008"`
}

func main() {
	var opts Opts

	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println("Failed to parse args: ", err)
			os.Exit(2)
		}
	}

	h := &handler.Handler{RestURL: opts.Rest}
	r := router.Init(h)
	r.Run(":" + opts.Port)
}
