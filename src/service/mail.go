package service

import (
	"encoding/json"
	"log"
	"net/smtp"
	"os"
)

var sender Sender

// Sender struct
type Sender struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	SMTPServer string `json:"smtpserver"`
}

func init() {
	file, err := os.Open("./config/config.mail.json")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&sender); err != nil {
		log.Println(err.Error())
	}
}

func (sender Sender) sendMessage(dest, subject, message string) {
	msg := "From: " + sender.User + "\n" +
		"To: " + dest + "\n" +
		"Subject: " + subject + "\n" + message

	err := smtp.SendMail(sender.SMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, sender.SMTPServer),
		sender.User, []string{dest}, []byte(msg))

	if err != nil {
		log.Println("Mail not sent")
		return
	}

	log.Printf("Mail to %s with description of a problem has successfully sent\n\n", dest)
}
