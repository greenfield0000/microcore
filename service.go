package main

import "github.com/sirupsen/logrus"

type CoreService interface {
	GetLogger() *logrus.Logger
}
