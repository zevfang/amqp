package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
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
		"logs_topic",
		"topic", //声明主题模式
		true,    //持久化
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal("声明交换机错误", err)
	}

	body := bodyFrom(os.Args)

	if err := ch.Publish(
		"logs_topic",
		severityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	); err != nil {
		log.Fatal("发送消息错误", err)
	}

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

//获取路由键
func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
