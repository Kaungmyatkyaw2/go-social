package mailer

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	fromName  string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail, fromName string) (*SendGridMailer, error) {

	if apiKey == "" {
		return &SendGridMailer{}, fmt.Errorf("sendgrid api key is required")
	}

	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		fromName:  fromName,
		apiKey:    apiKey,
		client:    client,
	}, nil
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	from := mail.NewEmail(m.fromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmplData, err := buildTemplate(templateFile, data)

	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, tmplData.subject.String(), to, "", tmplData.body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	retryAttempt(func() error {
		response, err := m.client.Send(message)

		if err != nil {
			return nil
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			log.Printf("Email sent with status code %v", response.StatusCode)

			return nil
		}

		return fmt.Errorf("failed to sent email with status %v", response.StatusCode)

	}, maxRetries, email)

	return fmt.Errorf("failed to send email after %d attemtps", maxRetries)
}
