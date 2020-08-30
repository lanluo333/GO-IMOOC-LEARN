package main

import "imooc-rabbitmq/RabbitMQ"

func main()  {
	imoocTwo := RabbitMQ.NewRabbitMQRouting("route_imooc", "imooc_two")
	imoocTwo.ReceiveRouting()
}

