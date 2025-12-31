package main

import (
	"log"
	"os"

	_ "github.com/Piccadilly98/subscription_service/docs"
	"github.com/Piccadilly98/subscription_service/internal/server"
)

// @title           Subscriptions service
// @description     REST API для управления данными о подписках пользователей и расчёта их стоимости. Реализован в качестве тестового задания на вакансию junior Goolang для Effective Mobile.
// @termsOfService  https://effective-mobile.ru

// @contact.name    MyResumeHH
// @contact.url     https://hh.ru/resume/875f5df3ff0fbfea5f0039ed1f39506e574431
// @contact.email   senior.filippov2017@yandex.ru

// @host            		 localhost:8080
// @BasePath        		 /
// @schemes         		 http
// @externalDocs.description Test Task By Effective Mobile
// @externalDocs.url 		 https://disk.360.yandex.ru/i/iT7aYzpZH_-lVg

func main() {
	server, err := server.InitServer()
	if err != nil {
		os.Exit(1)
	}
	ch, err := server.Start()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = <-ch
	log.Fatal(err)

}
