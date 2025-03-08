package main

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wneessen/go-mail"
	"os"
)

func sendMail(info TrackData, id string) {
	var recipient, name string
	err := dbpool.QueryRow(context.Background(), "SELECT email, name FROM mail_receipts WHERE id = $1", id).Scan(&recipient, &name)
	if err != nil {
		log.Errorf("failed to get email for id %s: %v", id, err)
		return
	}

	if recipient == "" {
		return
	}

	message := mail.NewMsg()
	if err := message.From(os.Getenv("EMAIL_USERNAME")); err != nil {
		log.Errorf("failed to set From address: %s", err)
		return
	}

	if err := message.To(recipient); err != nil {
		log.Errorf("failed to set To address: %s", err)
		return
	}

	if info.Url != "" {
		message.Subject("Somebody clicked your link!")
		message.SetBodyString(mail.TypeTextPlain,
			"Hi 👋,\n\nSomeone ("+info.Ip+") was brought to \""+info.Url+"\" after clicking the link in the email you call \""+name+
				"\"\nThis happened at: "+info.Timestamp+
				"\n\nTheir browser identified as: \""+info.UserAgent+
				"\"\n\nHave a great day,\nmailReceipt")
	} else {
		message.Subject("Somebody opened your mail!")
		message.SetBodyString(mail.TypeTextPlain,
			"Hi 👋,\n\nSomeone ("+info.Ip+") opened the email you call \""+name+
				"\"\nThis happened at: "+info.Timestamp+
				"\n\nTheir browser identified as: \""+info.UserAgent+
				"\"\n\nHave a great day,\nmailReceipt")
	}

	if err := mailClient.DialAndSend(message); err != nil {
		log.Errorf("failed to send mail: %s", err)
		return
	}
}
