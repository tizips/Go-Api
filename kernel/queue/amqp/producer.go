package amqp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"saas/kernel/config"
	"saas/kernel/logger"
	"time"
)

func (connect *Connection) InitProducer(isClose bool) {

	if isClose {
		connect.CloseProcess <- true
	}

	conn, err := dial()
	if err != nil {
		logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 生产者连接异常: %v", err))
		time.Sleep(time.Duration(config.Values.Amqp.Reconnect) * time.Second)
		connect.InitProducer(false)
		return
	}

	defer conn.Close()

	logger.Logger.Amqp.Info("Amqp 生产者连接成功")

	// 打开一个并发服务器通道来处理消息
	ch, err := conn.Channel()
	if err != nil {
		logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 打开通道异常:%s", err.Error()))
		return
	}

	defer ch.Close()

	connect.Conn = conn
	connect.Ch = ch
	connect.CloseProcess = make(chan bool, 1)
	connect.ProducerReConnect()
}

func (connect *Connection) Publish(body []byte, queue string, retries *int) error {

	if connect.RabbitProducerMap == nil {
		return errors.New("未初始化生产者信息")
	}

	if queue == "" {
		return errors.New("队列名称不能为空")
	}

	exchange := connect.RabbitProducerMap[queue]
	if exchange == "" {
		return errors.New("交换机名称不能为空")
	}

	m := Producer{
		Queue:   queue,
		Body:    string(body),
		Retries: 0,
	}

	if retries != nil {
		m.Retries = *retries + 1
	}

	body, _ = json.Marshal(m)

	// 发布
	err := connect.Ch.Publish(
		exchange, // exchange 默认模式，exchange为空
		queue,    // routing key 默认模式路由到同名队列，即是task_queue
		false,    // mandatory
		false,
		amqp.Publishing{
			// 持久性的发布，因为队列被声明为持久的，发布消息必须加上这个（可能不用），但消息还是可能会丢，如消息到缓存但MQ挂了来不及持久化。
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

	if err != nil {
		logger.Logger.Amqp.Error("Amqp 发送消息失败:" + err.Error())
		return err
	}

	return nil
}

//ProducerReConnect 生产者重连
func (connect *Connection) ProducerReConnect() {
closeTag:
	for {
		connect.ConnNotifyClose = connect.Conn.NotifyClose(make(chan *amqp.Error))
		connect.ChNotifyClose = connect.Ch.NotifyClose(make(chan *amqp.Error))
		select {
		case connErr, _ := <-connect.ConnNotifyClose:
			if connErr != nil {
				logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 连接异常:%s", connErr.Error()))
			}
			// 判断连接是否关闭
			if !connect.Conn.IsClosed() {
				if err := connect.Conn.Close(); err != nil {
					logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 连接关闭异常:%s", err.Error()))
				}
			}
			//重新连接
			if conn, err := dial(); err != nil {
				logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 重连失败:%s", err.Error()))
				_, isConnChannelOpen := <-connect.ConnNotifyClose
				if isConnChannelOpen {
					close(connect.ConnNotifyClose)
				}
				//connection关闭时会自动关闭channel
				connect.InitProducer(false)
				//结束子进程
				break closeTag
			} else { //连接成功
				connect.Ch, _ = conn.Channel()
				connect.Conn = conn
				//logger.Logger.Amqp.Info("rabbitMQ重连成功")
			}
			// IMPORTANT: 必须清空 Notify，否则死连接不会释放
			for err := range connect.ConnNotifyClose {
				println(err)
			}
		case chErr, _ := <-connect.ChNotifyClose:
			if chErr != nil {
				logger.Logger.Amqp.Error(fmt.Sprintf("Amqp 通道连接关闭:%s", chErr.Error()))
			}
			// 重新打开一个并发服务器通道来处理消息
			if !connect.Conn.IsClosed() {
				ch, err := connect.Conn.Channel()
				if err != nil {
					logger.Logger.Amqp.Error(fmt.Sprintf("Amqp channel重连失败:%s", err.Error()))
					connect.ChNotifyClose <- chErr
				} else {
					logger.Logger.Amqp.Info("Amqp 通道重新创建成功")
					connect.Ch = ch
				}
			} else {
				_, isConnChannelOpen := <-connect.ConnNotifyClose
				if isConnChannelOpen {
					close(connect.ConnNotifyClose)
				}
				connect.InitProducer(false)
				break closeTag
			}
			for err := range connect.ChNotifyClose {
				println(err)
			}
		case <-connect.CloseProcess:
			break closeTag
		}
	}
}
