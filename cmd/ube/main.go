package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gidor/ube/pkg/http/rest"
	"github.com/gidor/ube/pkg/infra"
)

func main() {
	config := flag.String("config", "", "api config files")
	// configdir := flag.String("config Dir", "", "api config directory")
	if *config == "" {
		args := flag.Args()
		if len(args) > 0 {
			infra.Setenv("ube_cfg", args[0])
		}
	} else {
		infra.Setenv("ube_cfg", *config)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run() error {
	// init logger
	// level 0 is a info level ,so debug level doesn't show
	// debug level can show in level -1
	// logger, err := infra.CreateLogger(0)
	_, err := infra.GetLogger()
	if err != nil {
		return err
	}

	// init storage

	// init repositories with given logger and storage

	// init services with given logger and repository

	// setup routes
	restHandler := rest.CreateHandler()
	restHandler.AddHealthCheckHandler()
	restHandler.AddApi()

	// listen and serve
	// webServer := server.CreateServer(restHandler.GetRouter(), ":"+os.Getenv("HTTP_PORT"))
	webServer := infra.CreateServer(restHandler.GetRouter(), ":3000")
	log.Println("starting server...")
	return webServer.ListenAndServe()
}
