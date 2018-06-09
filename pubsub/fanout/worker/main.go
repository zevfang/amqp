package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open a channel")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Failed to declare an exchange")
	}

	q, err := ch.QueueDeclare(
		"",
		true,
		false,
		true, //独占
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Failed to declare an queue")
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, //自动确认
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Failed to register a consumer")
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
