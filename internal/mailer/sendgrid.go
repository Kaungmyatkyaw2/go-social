package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	fromName string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail,fromName string) (*SendGridMailer,error) {

	if apiKey == "" {
		return &SendGridMailer{}, fmt.Errorf("sendgrid api key is required")
	}


	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		fromName: fromName,
		apiKey:    apiKey,
		client:    client,
	},nil 
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {


	from := mail.NewEmail(m.fromName,m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl,err := template.ParseFS(FS,"templates/"+templateFile)

	if err !=nil {
		return err 
	}

	
	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)


	err = tmpl.ExecuteTemplate(subject,"subject",data)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(body,"body",data)
	if err != nil {
		return err
	}


	message := mail.NewSingleEmail(from,subject.String(),to,"",body.String())


	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})


	for i := range(maxRetries) {
		response, err := m.client.Send(message)

		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d",email,i + 1,maxRetries)
			log.Printf("Error: %v",err.Error())


			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			log.Printf("Email sent with status code %v",response.StatusCode)
		
			return nil
		}

		// return fmt.Errorf("something went wrong during the mailing attempts with status code %v",response.StatusCode)
	}

	return fmt.Errorf("failed to send email after %d attemtps",maxRetries) 
}