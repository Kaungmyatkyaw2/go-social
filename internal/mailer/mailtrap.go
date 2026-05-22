package mailer

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	fromName  string
	host      string
	port      int
	username  string
	password  string
}

func NewMailTrap(host string, port int, username, password, fromEmail, fromName string) (*mailtrapClient, error) {

	if username == "" {
		return &mailtrapClient{}, fmt.Errorf("mailtrap username is required")
	}

	return &mailtrapClient{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
		fromName:  fromName,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSandbox bool) error {

	tmplData, err := buildTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", tmplData.subject.String())

	message.AddAlternative("text/html", tmplData.body.String())

	dialer := gomail.NewDialer(m.host, m.port, m.username, m.password)

	return retryAttempt(func() error {
		if err := dialer.DialAndSend(message); err != nil {
			return err
		}

		return nil

	}, maxRetries, email)

}
