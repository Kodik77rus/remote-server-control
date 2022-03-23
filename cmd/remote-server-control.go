package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"remote-server-control/internal/server"
)

var certPath string

func init() {
	flag.StringVar(&certPath, "path", "configs/certificate", "path to sll certificate folder")
}

func main() {
	flag.Parse()

	//load base server config
	config := server.NewConfig(certPath)

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//Set server config
	s := server.New(config, ctx)

	//os signal handler
	go func() {
		oscall := <-c
		log.Printf("System call:%+v", oscall)
		s.Shutdown()
		cancel()
	}()

	setWorkDir()

	fatalError(s.Start())
}

//user work dir
func setWorkDir() {
	dir, err := os.UserHomeDir()
	fatalError(err)
	if err := os.Chdir(dir); err != nil {
		fatalError(err)
	}
}

func fatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
