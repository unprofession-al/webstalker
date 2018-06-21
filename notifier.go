package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

		if strings.HasPrefix(key, "sitewatch_notifier") {
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
	Notify(recipient string, message string) error
}

type StdOutNotifier struct{}

func NewStdOutNotifier(c string) (Notifier, error) {
	return StdOutNotifier{}, nil
}

func (son StdOutNotifier) Notify(r, m string) error {
	fmt.Printf("\tTo: %s\n\t%s\n", r, m)
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
	log.Printf("Sending as '%s' with API KEY '%s'", n.Sender, n.APIKey)

	return n, nil
}

func (sgn SendGridNotifier) Notify(r, m string) error {
	from := mail.NewEmail(sgn.Sender, sgn.Sender)
	subject := "Updates from sitewatch"
	to := mail.NewEmail(r, r)
	plainTextContent := m
	htmlContent := m
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sgn.APIKey)
	_, err := client.Send(message)
	return err
}
