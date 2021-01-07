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

const (
	exchange = "zyexchange"
	queue    = "zyqueue"
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

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// binding
	err = ch.QueueBind(
		q.Name,   // queue name
		q.Name,   // routing key
		exchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue to exchange")

	body := "Hello World!"
	err = ch.Publish(
		exchange, // exchange
		q.Name,   // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
