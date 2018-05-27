package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"task_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	body := bodyFrom(os.Args)

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [x] Sent %s", body)

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [x] Sent %s", body)

}

func bodyFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = fmt.Sprintf("%s-%v", "hello", time.Now())
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
