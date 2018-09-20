package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/telegram-bot-api.v4"
)

var notifiers map[string]func(string) (Notifier, error)

func init() {
	notifiers = make(map[string]func(string) (Notifier, error))
	notifiers["stdout"] = NewStdOutNotifier
	notifiers["sendgrid"] = NewSendGridNotifier
	notifiers["telegram"] = NewTelegramNotifier
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

	return n, nil
}

func (sgn SendGridNotifier) Notify(r, m string) error {
	from := mail.NewEmail(sgn.Sender, sgn.Sender)
	subject := "Updates from webstalker"
	to := mail.NewEmail(r, r)
	plainTextContent := m
	htmlContent := m
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sgn.APIKey)
	_, err := client.Send(message)
	return err
}

type TelegramNotifier struct {
	bot      *tgbotapi.BotAPI
	chatName string
}

func NewTelegramNotifier(c string) (Notifier, error) {
	n := TelegramNotifier{}

	tokens := strings.Fields(c)
	if len(tokens) != 2 {
		return n, fmt.Errorf("Malformend config for Telegram notifier: %s", c)
	}

	n.chatName = tokens[0]
	bot, err := tgbotapi.NewBotAPI(tokens[1])
	if err != nil {
		return n, fmt.Errorf("Error while preparing Telegram notifier: %s", err.Error())
	}
	n.bot = bot

	return n, nil
}

func (tn TelegramNotifier) Notify(r, m string) error {
	msg := tgbotapi.NewMessageToChannel(tn.chatName, m)
	_, err := tn.bot.Send(msg)
	return err
}
