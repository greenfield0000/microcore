package serviceimpl

import (
	"crypto/tls"
	"errors"
	"github.com/greenfield0000/microcore/bussines/repository"
	"github.com/greenfield0000/microcore/bussines/service"
	"gopkg.in/gomail.v2"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

var host, login, password string
var port int

type plainMailSender struct {
	SmtpAuth smtp.Auth
}

func newPlainMailSender() *plainMailSender {
	host = os.Getenv("MAIL_HOST")
	port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	login = os.Getenv("MAIL_LOGIN")
	password = os.Getenv("MAIL_PASS")
	return &plainMailSender{
		SmtpAuth: smtp.PlainAuth("", login, password, host),
	}
}

func (c plainMailSender) sendMail(subject string, to string, message string) error {
	d := gomail.NewDialer(host, port, login, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// message
	m := gomail.NewMessage()
	m.SetHeader("From", login)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)
	// try send
	if err := d.DialAndSend(m); err != nil {
		return errors.New("Не удалось отправить почту")
	}

	return nil
}

type PlainMailService struct {
	Repo   *repository.Repository
	sender *plainMailSender
}

func NewPlainMailService(repo *repository.Repository) service.MailService {
	return PlainMailService{Repo: repo, sender: newPlainMailSender()}
}

func (ms PlainMailService) Send(subject string, to string, message string) error {
	if err := ms.sender.sendMail(subject, to, message); err != nil {
		log.Printf("Was error: %s", err.Error())
		return errors.New("Не удалось отправить почту!")
	}
	return nil
}
