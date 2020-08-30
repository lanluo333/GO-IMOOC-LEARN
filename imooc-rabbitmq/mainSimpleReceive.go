package main

import "imooc-rabbitmq/RabbitMQ"

func main()  {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}

