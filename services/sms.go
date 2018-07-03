package services

import (
	"gopkg.in/njern/gonexmo.v2"
)

const(
	apikey   = "7c4a13eb"
	apiSecret = "hJ11W6WwJRvvjvF8"
	userName = "n.a.m.62608@gmail.com"
	password = "Zikim62608"
	sender   = "EasyPay"
)


func SendSMS(message, recipient string) error {
	nexmoClient, _ := nexmo.NewClient(apikey, apiSecret)

	sms := &nexmo.SMSMessage{
		From            : sender,
		To              : recipient,
		Type            : nexmo.Text,
		Text            : message,
		Class           : nexmo.Standard,
	}

	_, err := nexmoClient.SMS.Send(sms)
	return err
}
