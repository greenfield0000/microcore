package service

import "github.com/sirupsen/logrus"

type CustomService struct {
	logger logrus.Logger
}

func NewCustomService(logger logrus.Logger) CustomService {
	return CustomService{
		logger: logger,
	}
}

func (c CustomService) DoLog(message string) {
	c.logger.Println(message)
}
