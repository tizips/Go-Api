package amqp

import "github.com/streadway/amqp"

type Connection struct {
	//连接
	Conn *amqp.Connection
	//通道
	Ch *amqp.Channel
	//连接异常结束
	ConnNotifyClose chan *amqp.Error
	//通道异常接收
	ChNotifyClose chan *amqp.Error
	//用于关闭进程
	CloseProcess chan bool
	//消费者列表
	Consumers []Consumer
	//生产者信息
	RabbitProducerMap map[string]string
	//自定义消费者处理函数
	ConsumeHandle func(amqp.Delivery)
}

type Consumer struct {
	//交换机
	Exchange string
	//交换机类型
	Type string
	//队列
	Queue string
	//名称
	Name string
	//处理数量
	Numbers int
}

type Producer struct {
	Queue   string `json:"queue"`
	Body    string `json:"body"`
	Retries int    `json:"retries"`
}
