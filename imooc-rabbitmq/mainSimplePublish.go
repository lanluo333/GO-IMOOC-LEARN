package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
)


func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.PublishSimple("Hello world!")
	fmt.Println("发送成功")
}



