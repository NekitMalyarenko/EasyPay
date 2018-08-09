package services

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"net/url"
	"strings"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"errors"
	"log"
)

const(
	username = "380980517016"
	token = "VC8l7CJ9bh8LJcv"
	from = "Registracia"
	lifetime = "1"
	smsClubUrl = "https://gate.smsclub.mobi/token/"
)

func SendSMS(recipientNumber, text string) error {
	win1251Text, err := textToWin1251(text)
	if err != nil {
		return err
	}

	if strings.Contains(recipientNumber, "38") {
		recipientNumber = string([]byte(recipientNumber)[2:])
		log.Println(recipientNumber)
	}

	form := url.Values{
		"username" : {username},
		"token"    : {token},
		"from"     : {from},
		"text"     : {string(win1251Text)},
		"lifetime" : {lifetime},
		"to"       : {recipientNumber},
	}

	rsp, err := http.Post(smsClubUrl, "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(bodyByte),"=IDS START="){
		return errors.New("Error sending sms with error:" + string(bodyByte))
	}
	return nil
}


func textToWin1251(utf8 string) ([]byte,error) {
	return ioutil.ReadAll(transform.NewReader(strings.NewReader(utf8),charmap.Windows1251.NewEncoder()))
}
