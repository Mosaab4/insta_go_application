package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"instabug/go/api"
	"instabug/go/util"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config", err)
	}

	connection, err := connect(config.RabbitMqURL)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	server := api.NewServer(connection)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Can not start server", err)
	}
}

func connect(url string) (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial(url)
		if err != nil {
			log.Println(url)
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off ...")
		time.Sleep(backOff)
		continue
	}

	log.Println("Connected to RabbitMQ")

	return connection, nil
}
