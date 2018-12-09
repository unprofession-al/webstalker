package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/template"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var notifiers map[string]func(string) (Notifier, error)

func init() {
	notifiers = make(map[string]func(string) (Notifier, error))
	notifiers["stdout"] = NewStdOutNotifier
	notifiers["sendgrid"] = NewSendGridNotifier
}

func PrepareNotifiers() ([]Notifier, error) {
	n := []Notifier{}
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		key := strings.ToLower(kv[0])
		value := kv[1]

		if strings.HasPrefix(key, "webstalker_notifier") {
			for name, fn := range notifiers {
				if strings.Contains(key, name) {
					notifier, err := fn(value)
					if err != nil {
						return n, err
					}
					log.Printf("Notifier config for '%s' found\n", name)
					n = append(n, notifier)
				}
			}
		}
	}

	if len(n) < 1 {
		log.Println("No notifier configured, using StdOut Notifier")
		n = append(n, StdOutNotifier{})
	}

	return n, nil
}

type Notifier interface {
	Notify(recipient string, message string, diff string) error
}

func renderTemplate(m, d string) (string, error) {
	var data = struct {
		Diff string
	}{
		Diff: d,
	}

	tmpl, err := template.New("tmpl").Parse(m)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

type StdOutNotifier struct{}

func NewStdOutNotifier(c string) (Notifier, error) {
	return StdOutNotifier{}, nil
}

func (son StdOutNotifier) Notify(r, m, d string) error {
	msg, err := renderTemplate(m, d)
	if err != nil {
		return err
	}
	fmt.Printf("\tTo: %s\n\t%s\n", r, msg)
	return nil
}

type SendGridNotifier struct {
	APIKey string
	Sender string
}

func NewSendGridNotifier(c string) (Notifier, error) {
	n := SendGridNotifier{}

	tokens := strings.Fields(c)
	if len(tokens) != 2 {
		return n, fmt.Errorf("Malformend config for SendGrid notifier: %s", c)
	}
	n.Sender = tokens[0]
	n.APIKey = tokens[1]

	return n, nil
}

func (sgn SendGridNotifier) Notify(r, m, d string) error {
	msg, err := renderTemplate(m, d)
	if err != nil {
		return err
	}
	from := mail.NewEmail(sgn.Sender, sgn.Sender)
	subject := "Updates from webstalker"
	to := mail.NewEmail(r, r)
	plainTextContent := msg
	htmlContent := msg
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sgn.APIKey)
	_, err = client.Send(message)
	return err
}
