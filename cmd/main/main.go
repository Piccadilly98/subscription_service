package main

import (
	"log"
	"time"

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

func GetTimeDate(t time.Time) *time.Time {
	return &t
}

func getIntPtr(i int) *int {
	return &i
}

func getBoolPtr(b bool) *bool {
	return &b
}

func getPtrStr(str string) *string {
	return &str
}
