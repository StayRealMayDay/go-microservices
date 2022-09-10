package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Application struct {
	Rabbit *amqp.Connection
}

func main() {

	conn, err := connect()
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
	defer conn.Close()

	app := Application{
		Rabbit: conn,
	}

	log.Printf("start broker server at port: %s", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routers(),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var count int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//do not continue untill rabbitmq is reday

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmp")
		if err != nil {
			fmt.Println("Rabbitmq is not ready")
			count++
		} else {
			log.Println("connected to rabbitMQ")
			connection = c
			break
		}

		if count > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		fmt.Println("backing off")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
