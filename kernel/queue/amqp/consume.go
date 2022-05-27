package amqp

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"saas/kernel/config"
	"saas/kernel/logger"
	"time"
)

// InitConsumer 初始化消费者
func (connect *Connection) InitConsumer(isClose bool) {

	if isClose {
		connect.CloseProcess <- true
	}

	//连接 Amqp
	conn, err := dial()
	if err != nil {
		logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 消费者连接异常:%v", err))
		time.Sleep(time.Duration(config.Values.Amqp.Reconnect) * time.Second)
		connect.InitConsumer(false)
		return
	}

	defer conn.Close()

	logger.Logger.Amqp.Info("Amqp 消费者连接成功")

	connect.Conn = conn

	err = connect.CreateRabbitMQConsumer()
	if err != nil {
		logger.Logger.Amqp.Error(err.Error())
		return
	}

	connect.CloseProcess = make(chan bool, 1)
	connect.ConsumerReConnect()
}

func (connect *Connection) CreateRabbitMQConsumer() error {

	if len(connect.Consumers) == 0 {
		return errors.New("消费者信息不能为空")
	}

	var err error

	for _, value := range connect.Consumers {

		//创建一个通道
		if connect.Ch, err = connect.Conn.Channel(); err != nil {
			return fmt.Errorf("打开 Amqp 通道失败: %v", err)
		}

		typing := "topic"

		if value.Type != "" {
			typing = value.Type
		}

		if err = connect.Ch.ExchangeDeclare(value.Exchange, typing, true, false, false, false, nil); err != nil {
			return fmt.Errorf("交换机初始化失败,交换机名称:%s, 错误:%v", value.Exchange, err)
		}

		_, err = connect.Ch.QueueDeclare(value.Queue, true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("队列初始化失败,队列名称:%s,错误: %s", value.Queue, err.Error())
		}

		// 绑定队列
		err = connect.Ch.QueueBind(value.Queue, value.Queue, value.Exchange, false, nil)
		if err != nil {
			return fmt.Errorf("绑定队列失败: %v", err)
		}

		numbers := 1

		if value.Numbers > 0 {
			numbers = value.Numbers
		}

		for i := 1; i <= numbers; i += 1 {

			consumer := fmt.Sprintf("Consumer-%s.%d", value.Queue, i)
			//绑定消费者
			messages := make(<-chan amqp.Delivery)
			messages, err = connect.Ch.Consume(value.Queue, consumer, false, false, false, false, nil)

			if err != nil {
				return fmt.Errorf("创建消费者 %s 失败: %v", consumer, err)
			}

			logger.Logger.Amqp.Info(fmt.Sprintf("Process[%s] start", consumer))

			if connect.ConsumeHandle != nil {
				go func() {
					for message := range messages {
						connect.ConsumeHandle(message)
					}
				}()
			}
		}
	}

	return nil
}

//ConsumerReConnect 消费者重连
func (connect *Connection) ConsumerReConnect() {
closeTag:
	for {

		connect.ConnNotifyClose = connect.Conn.NotifyClose(make(chan *amqp.Error))
		connect.ChNotifyClose = connect.Ch.NotifyClose(make(chan *amqp.Error))

		var err *amqp.Error

		select {
		case err, _ = <-connect.ConnNotifyClose:
		case err, _ = <-connect.ChNotifyClose:

			if err != nil {
				logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 消费者连接异常: %v", err))
			}

			// 判断连接是否关闭
			if !connect.Conn.IsClosed() {
				if err := connect.Conn.Close(); err != nil {
					logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 连接关闭异常: %v", err))
				}
			}

			_, isConnChannelOpen := <-connect.ConnNotifyClose
			if isConnChannelOpen {
				close(connect.ConnNotifyClose)
			}

			connect.InitConsumer(false)

			break closeTag
		case <-connect.CloseProcess:
			break closeTag
		}
	}

	logger.Logger.Amqp.Info("结束消费者旧进程")
}
