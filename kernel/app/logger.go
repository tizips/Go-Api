package app

import "github.com/sirupsen/logrus"

var Logger struct {
	Api       *logrus.Logger
	SQL       *logrus.Logger
	Exception *logrus.Logger
	Amqp      *logrus.Logger
}
