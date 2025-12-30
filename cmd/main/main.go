package main

import (
	"log"

	"github.com/Piccadilly98/subscription_service/internal/server"
)

func main() {
	server, err := server.InitServer()
	if err != nil {
		log.Fatal(err)
	}
	ch, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = <-ch
	log.Fatal(err)

}
