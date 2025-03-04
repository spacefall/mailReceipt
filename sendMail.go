package main

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wneessen/go-mail"
)

func sendMail(info TrackData, id string) {
	var recipient string
	err := dbpool.QueryRow(context.Background(), "SELECT email FROM mail_receipts WHERE id = $1", id).Scan(&recipient)
	if err != nil {
		log.Errorf("failed to get email for id %s: %v", id, err)
		return
	}

	if recipient == "" {
		return
	}

	message := mail.NewMsg()
	if err := message.From("mailreceiptbot@gmail.com"); err != nil {
		log.Errorf("failed to set From address: %s", err)
		return
	}

	if err := message.To(recipient); err != nil {
		log.Errorf("failed to set To address: %s", err)
		return
	}

	if info.Url != "" {
		message.Subject("Somebody clicked your link!")
		message.SetBodyString(mail.TypeTextPlain, "Hi 👋,\nSomeone ("+info.Ip+") clicked your link at "+info.Timestamp+"!\nTheir device identified as: "+
			info.UserAgent+" and they were brought to: "+info.Url+".\n\nHave a great day!")
	} else {
		message.Subject("Somebody opened your mail!")
		message.SetBodyString(mail.TypeTextPlain, "Hi 👋,\nSomeone ("+info.Ip+") opened your email at "+info.Timestamp+"!\nTheir device identified as: "+
			info.UserAgent+".\n\nHave a great day!")
	}

	if err := mailClient.DialAndSend(message); err != nil {
		log.Errorf("failed to send mail: %s", err)
		return
	}
}
