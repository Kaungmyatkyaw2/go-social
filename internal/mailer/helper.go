package mailer

import (
	"bytes"
	"fmt"
	"html/template"
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

func retryAttempt(send emailSendFn, maxAttempt int, toEmail string) error {

	var respErr error

	for i := range maxAttempt {
		respErr := send()

		if respErr != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}


		return nil 

	}

	return fmt.Errorf("failed to send email after %d attemtps, error: %v", maxRetries, respErr)
}
