package mailer

import (
	"bytes"
	"html/template"
	"log"
	"time"
)

type templateRes struct {
	subject bytes.Buffer
	body    bytes.Buffer
}

func buildTemplate(templateFile string, data any) (*templateRes, error) {

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)

	if err != nil {
		return nil, err
	}

	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, err
	}

	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil, err
	}

	return &templateRes{
		subject: *subject,
		body:    *body,
	}, nil
}

type emailSendFn func() error

func retryAttempt(send emailSendFn, maxAttempt int, toEmail string) {
	for i := range maxAttempt {
		err := send()

		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", toEmail, i+1, maxAttempt)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email Successfully send to %v", toEmail)

		return

	}
}
