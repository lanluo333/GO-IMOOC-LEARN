package main

import "imooc-rabbitmq/RabbitMQ"

func main()  {
	imoocOne := RabbitMQ.NewRabbitMQRouting("route_imooc", "imooc_one")
	imoocOne.ReceiveRouting()
}

