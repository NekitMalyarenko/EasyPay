package services

import (
	"github.com/plivo/plivo-go"
	"log"
)

const(
	token = "NGE2YThlNTVmODlkMDcwNGFlMTlkNzcxODkzOWE2"
	id    = "MANDZKMDBKYJAWMMY0YM"
)


func SendSMS(message, recipient string) error {
	client, err := plivo.NewClient(id, token, &plivo.ClientOptions{})
	if err != nil {
		return err
	}

	log.Println(recipient)

	_, err = client.Messages.Create(plivo.MessageCreateParams{
		Src: "380506260859",
		Dst: recipient,
		Text: message,
	})

	return err
}
