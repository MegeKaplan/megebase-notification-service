package main

import (
	"log"
	"net/smtp"
	"os"
)

type Mailer interface {
	Send(to []string, subject, body string) error
}

type GmailMailer struct {
	email    string
	password string
	host     string
	port     string
}

func NewGmailMailer() *GmailMailer {
	email := os.Getenv("GMAIL_EMAIL")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("GMAIL_EMAIL or GMAIL_APP_PASSWORD not set")
	}
	host := "smtp.gmail.com"
	port := "587"
	return &GmailMailer{
		email:    email,
		password: password,
		host:     host,
		port:     port,
	}
}

func (m *GmailMailer) Send(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", m.email, m.password, m.host)
	msg := buildMessage("Megebase", m.email, to, subject, body, true)
	return smtp.SendMail(m.host+":"+m.port, auth, m.email, to, msg)
}

func NewMailer(provider string) Mailer {
	switch provider {
	case "gmail":
		return NewGmailMailer()
	default:
		log.Fatalf("Unknown mail provider: %s", provider)
		return nil
	}
}

func buildMessage(fromName, fromEmail string, to []string, subject, body string, isHTML bool) []byte {
	headers := "From: " + fromName + " <" + fromEmail + ">\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n"

	if isHTML {
		headers += "MIME-version: 1.0;\r\n" + "Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	} else {
		headers += "Content-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"
	}

	return []byte(headers + body + "\r\n")
}
