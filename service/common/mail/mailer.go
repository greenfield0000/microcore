package commonservice

import (
	"crypto/tls"
	"errors"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"os"
	"strconv"
)

var host, login, password string
var port int

type MailSender interface {
	Send(subject string, to string, message string) error
}

type PlainMailSender struct {
	SmtpAuth smtp.Auth
}

func NewPlainMailSender() MailSender {
	host = os.Getenv("MAIL_HOST")
	port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	login = os.Getenv("MAIL_LOGIN")
	password = os.Getenv("MAIL_PASS")
	return PlainMailSender{
		SmtpAuth: smtp.PlainAuth("", login, password, host),
	}
}

func (c PlainMailSender) Send(subject string, to string, message string) error {
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
