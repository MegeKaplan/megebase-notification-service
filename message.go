package main

import (
	"encoding/json"
	"log"
)

var mailer Mailer

func init() {
	mailer = NewMailer("gmail")
}

func ParseMessageEvent(msg []byte) (MessageEvent, error) {
	var event MessageEvent
	err := json.Unmarshal(msg, &event)
	if err != nil {
		return MessageEvent{}, err
	}
	return event, nil
}

func buildEmail(entity string, data map[string]interface{}) (string, string) {
	switch entity {
	case "otp":
		otp, ok := data["otp"].(string)
		if !ok {
			log.Println("Invalid OTP data")
			return "Invalid OTP", "<p>Invalid OTP data received. Please contact support.</p>"
		}
		return "Your OTP Code", "<h1>Your OTP Code</h1><p>Code: <b>" + otp + "</b></p>"
	default:
		return "Megebase Notification", "<p>New notification from Megebase.</p>"
	}
}

func SendMessage(msg MessageEvent) {
	switch msg.Channel {
	case "email":
		subject, body := buildEmail(msg.Entity, msg.Data)
		sendEmail(subject, body, msg.To)
	default:
		log.Printf("Unsupported channel: %s", msg.Channel)
	}
}

func sendEmail(subject string, msg interface{}, to string) {
	err := mailer.Send([]string{to}, subject, msg.(string))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
