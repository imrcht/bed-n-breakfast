package main

import (
	"time"

	"github.com/imrcht/bed-n-breakfast/internals/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMailChan() {
	// * This function will execute in the background and listen for any value that comes through the channel all the time
	go func() {
		for {
			// * Listening to mail channel
			msg := <-app.MailChan
			sendMessage(msg)
		}
	}()
}

func sendMessage(m models.MailData) {
	// * Sending mail
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Encryption = mail.EncryptionSTARTTLS
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = nil

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	app.InfoLog.Printf("Mail sent to %s", m.To)
}
