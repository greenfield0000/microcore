package serviceimpl

import (
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	commonservice "github.com/greenfield0000/microcore/service/common/mail"
	"log"
)

type MailerService struct {
	Repo   *repository.Repository
	sender commonservice.MailSender
}

func NewMailService(repo *repository.Repository) MailerService {
	return MailerService{Repo: repo, sender: commonservice.NewPlainMailSender()}
}

func (ms MailerService) SendMail(subject string, to string, message string) error {
	if err := ms.sender.Send(subject, to, message); err != nil {
		log.Printf("Was error: %s", err.Error())
		return errors.New("Не удалось отправить почту!")
	}
	return nil
}
