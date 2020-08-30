package main

import "imooc-rabbitmq/RabbitMQ"

func main()  {
	rabbitMq := RabbitMQ.NewRabbitMQPubSub("newProduct")
	rabbitMq.ReceiveSub()
}

