package main

import (
	"encoding/json"
	"log"
)

func ParseMessageEvent(msg []byte) (MessageEvent, error) {
	var event MessageEvent
	err := json.Unmarshal(msg, &event)
	if err != nil {
		return MessageEvent{}, err
	}
	return event, nil
}

func SendMessage(msg MessageEvent) {
	switch msg.Channel {
	case "email":
		SendEmail(msg.Data[msg.Entity], msg.To)
	default:
		log.Printf("Unsupported channel: %s", msg.Channel)
	}
}

func SendEmail(msg interface{}, to string) {
	log.Printf("Sending email to %s with message: %s", to, msg)
	// add real email sending logic here
}
