package mailer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailGunMailer struct {
	fromEmail   string
	apiKey      string
	sandboxMail string
	client      *mailgun.MailgunImpl
}

func NewMailGun(apiKey, fromEmail string) *MailGunMailer {
	mg := mailgun.NewMailgun(fromEmail, apiKey)

	return &MailGunMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    mg,
	}
}

func (m *MailGunMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	sender := m.fromEmail
	recipient := email

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mailgun.NewMessage(sender, subject.String(), body.String(), recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i := 0; i < maxRetires; i++ {
		resp, id, err := m.client.Send(ctx, message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetires)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		fmt.Printf("ID: %s Resp: %s\n", id, resp)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetires)
}
