package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main()  {
	rabbitMq := RabbitMQ.NewRabbitMQPubSub("newProduct")
	for i := 0; i<= 100; i++  {
		message := "订阅模式下生产的第 " + strconv.Itoa(i) + " 条消息"
		rabbitMq.PublishPub(message)
		fmt.Println(message)
		time.Sleep(1 * time.Second)
	}
}

