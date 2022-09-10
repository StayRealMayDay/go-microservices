package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	//start listen message
	log.Println("Listening for and comsuming RabbitMQ message...")

	//create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	//watch the queue and consume event
	err = consumer.Listen([]string{"log.INFO", "log.WORNING", "log.ERROR"})
	if err != nil {
		log.Panicln(err)
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//do not continue untill rabbitmq is reday

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
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
