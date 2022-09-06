package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Application struct{

}

func main(){

	app := Application{}

	log.Printf("start broker server at port: %s", webPort)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routers(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}