package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// 格式 amqp://账号:密码@rabbirmq服务器地址:端口号/vhost
const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMQ struct {
	conn	*amqp.Connection
	channel	*amqp.Channel
	// 队列名
	QueueName	string
	// 交换机
	Exchange	string
	// key
	Key			string
	// 连接信息
	Mqurl		string
}

// 创建MQ结构体实例
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ  {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}

	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "创建连接失败")
	rabbitmq.channel, err = rabbitmq.conn.Channel()

	rabbitmq.failOnError(err, "获取channel失败")

	return rabbitmq

}

// 断开
func (r *RabbitMQ) Destory()  {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnError(err error, message string)  {
	if err != nil {
		log.Fatalf("%s:%s", err, message)
		panic(fmt.Sprintf("%s:%s", err, message))
	}
}

// 简单模式step : 1.创建简单模式下rabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ  {

	return NewRabbitMQ(queueName, "", "")

}

// 简单模式step: 2.简单模式下生产者
func (r *RabbitMQ) PublishSimple(message string)  {
	// 1. 申请队列，如果队列不存在会自动创建，如果存在则跳过
	_, err := r.channel.QueueDeclare(
			r.QueueName,
			false, //是否持久化
			false, // 是否自动删除
			false,  // 是否排他性
			false,  // 是否阻塞
			nil,		// 额外属性
		)

	if err != nil {
		fmt.Println(err)
	}

	// 2.发送消息到队列中
	err = r.channel.Publish(
			r.Exchange,
			r.QueueName,
			// 如果为true，就会根据exchange 和 routkey规则，如果没有找到符合规则的，就会把消息发送回去发送者
			false,
			// 如果为true, 当exchange发送消息到队列以后发现队列上没绑定的消费者，就会把消息发送回给发送者
			false,
			amqp.Publishing{
				ContentType:"text/plan",
				Body:[]byte(message),
			},
		)
	r.failOnError(err, "发送消息失败")
}

func (r *RabbitMQ) ConsumeSimple()  {
	// 1. 申请队列，如果队列不存在会自动创建，如果存在则跳过
	_, err := r.channel.QueueDeclare(
			r.QueueName,
			false,
			false,
			false,
			false,
			nil,
		)

	if err != nil {
		fmt.Println(err)
	}

	// 接受消息
	msgs, err := r.channel.Consume(
			r.QueueName,
			"",  // 用来区分多个消费者
			true, // 是否自动应答
			false,  // 是否具有排他性
			false,  // 如果为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
			false,  // 消息是否阻塞
			nil,
		)

	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		for d := range msgs{
			// 实现我们要处理的函数
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}

// 订阅模式创建rabbitmq
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ  {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "创建连接失败")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "创建channel失败")

	return rabbitmq
}

// 订阅模式生产者
func (r *RabbitMQ) PublishPub(message string)  {
	err := r.channel.ExchangeDeclare(
			r.Exchange,
			"fanout", // 设置交换机的类型为广播类型
			true,
			false,
			false, // 如果为true，表示这个exchange不可以被client用来推送消息，仅用来进行exchange之间的绑定
			false,
			nil,
		)
	r.failOnError(err, "创建交换机失败")

	// 发送消息
	err = r.channel.Publish(
			r.Exchange,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType:     "text/plain",
				Body : []byte(message),
			},
		)
	r.failOnError(err, "发送消息失败")
}

// 订阅模式消费者
func (r *RabbitMQ) ReceiveSub()  {
	// 创建交换机
	err := r.channel.ExchangeDeclare(
			r.Exchange,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
	r.failOnError(err, "创建连接失败")

	// 创建队列
	q, err := r.channel.QueueDeclare(
			"",  // 随机生产队列名称
			false,
			false,
			true,
			false,
			nil,
		)
	r.failOnError(err, "创建队列失败")

	// 绑定队列和交换机
	err = r.channel.QueueBind(
			q.Name,
			"",  // Pub/Sub模式下这个值要为空
			r.Exchange,
			false,
			nil,
		)
	r.failOnError(err, "队列绑定交换机失败")

	// 消费消息
	messages, err := r.channel.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
	
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()

	fmt.Println("退出请按 Ctrl+C")
	<-forever
}

// 路由模式
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ  {
	rabbitMq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	rabbitMq.conn, err = amqp.Dial(rabbitMq.Mqurl)
	rabbitMq.failOnError(err, "创建连接失败")

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnError(err, "创建channel失败")

	return rabbitMq
}

func (r *RabbitMQ) PublishRouting(message string)  {
	// 尝试创建连接
	err := r.channel.ExchangeDeclare(
			r.Exchange,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
	r.failOnError(err, "创建交换机失败")

	// 发送消息
	err = r.channel.Publish(
			r.Exchange,
			r.Key,
			false,
			false,
			amqp.Publishing{
				ContentType:     "text/plain",
				Body : []byte(message),
			},
		)
	r.failOnError(err, "发送消息失败")
}

// 路由模式下接受消息
func (r *RabbitMQ) ReceiveRouting()  {
	// 尝试创建连接
	err := r.channel.ExchangeDeclare(
			r.Exchange,
			"direct",
			false,
			false,
			false,
			false,
			nil,
		)
	r.failOnError(err, "创建交换机失败")

	// 创建队列
	q, err := r.channel.QueueDeclare(
			"",
			false,
			false,
			false,
			false,
			nil,
		)
	r.failOnError(err, "创建队列失败")

	// 绑定
	err = r.channel.QueueBind(
			q.Name,
			r.Key,
			r.Exchange,
			false,
			nil,
		)
	r.failOnError(err, "队列绑定交换机失败失败")

	// 消费消息
	messages, err := r.channel.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()

	fmt.Println("退出请按 Ctrl+C")
	<-forever
}

