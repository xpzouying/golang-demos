package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	mqUser     string
	mqPassword string
	mqHost     string
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	flag.StringVar(&mqUser, "user", "rabbitmq", "user for rabbitmq")
	flag.StringVar(&mqPassword, "password", "rabbitmq", "password for rabbitmq")
	flag.StringVar(&mqHost, "host", "localhost", "host of rabbitmq")

	flag.Parse()

	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:5672/", mqUser, mqPassword, mqHost),
	)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	const queue = "zyqueue"
	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
