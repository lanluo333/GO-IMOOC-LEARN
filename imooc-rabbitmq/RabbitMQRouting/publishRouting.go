package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main()  {
	imoocOne := RabbitMQ.NewRabbitMQRouting("route_imooc", "imooc_one")
	imoocTwo := RabbitMQ.NewRabbitMQRouting("route_imooc", "imooc_two")

	for i := 0; i<= 100; i++ {
		imoocOne.PublishRouting("Hello imooc one " + strconv.Itoa(i))
		imoocTwo.PublishRouting("Hello imooc two " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}



