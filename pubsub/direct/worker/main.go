package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("打开链接错误", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("打开通道错误", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal("声明交换机错误", err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("声明队列错误", err)
	}

	if len(os.Args) < 2 {
		log.Fatalln("请输入： [info] [warning] [error] ", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", s)
		if err = ch.QueueBind(
			q.Name,
			s,
			"logs_direct",
			false,
			nil,
		); err != nil {
			log.Fatal("绑定队列错误", err)
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("注册消费者错误", err)
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
