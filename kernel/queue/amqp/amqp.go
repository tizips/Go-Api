package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"saas/kernel/config"
	"strings"
)

func dial() (*amqp.Connection, error) {

	vhost := config.Values.Amqp.Vhost

	if !strings.HasPrefix(vhost, "/") {
		vhost = "/" + vhost
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s", config.Values.Amqp.Username, config.Values.Amqp.Password, config.Values.Amqp.Host, config.Values.Amqp.Port, vhost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
