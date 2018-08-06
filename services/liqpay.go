package services

import (
	"encoding/json"
	"crypto/sha1"
	b64 "encoding/base64"
	"html/template"
	"bytes"
	"strconv"
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
)

const (
 	publicKey  = ""
 	privateKey = ""
	LiqpayURL  = "https://www.liqpay.ua/api/"
)


type DataSignature struct {
	Data      string
	Signature string
}


func GenerateP2P(data map[string]interface{}) (string, string) {
	orderId := data["orderId"].(int64)
	receiverCard := data["receiverCard"].(int64)
	amount := data["amount"].(int)

	parsedRecieverCard := strconv.FormatInt(receiverCard, 10)
	parsedOrderId := strconv.FormatInt(orderId, 10)

	return form(map[string]interface{}{
		"action": "p2p",
		"version": 3,
		"public_key": publicKey,
		"amount": amount,
		"currency": "UAH",
		"description": "EasyPay payment",
		"order_id": parsedOrderId,
		"receiver_card": parsedRecieverCard,
		"server_url" : "http://95.158.39.205:8080/callback",
		"sandbox" : "1",
		"result_url" : "exit",
	})
}


func GetPaymentStatus(orderId int64) (string, error) {
	rawResult := api(map[string]interface{}{
		"action": "status",
		"version": 3,
		"public_key": publicKey,
		"order_id": orderId,
	})

	var result map[string]interface{}
	err := json.Unmarshal(rawResult, &result)
	if err != nil {
		return "", err
	}

	return result["status"].(string), nil
}


func form(Data map[string]interface{}) (string, string) {
	DataBytes, _ := json.Marshal(Data)

	DataBase64 := makeData(string(DataBytes))
	SignBase64 := makeSignature(DataBase64)
	//fmt.Println("Signature:",SignBase64)

	buf := new(bytes.Buffer)

	t, _ := template.ParseFiles("static/liqpay_form.html")
	t.ExecuteTemplate(buf, "liqpay_form.html", DataSignature{
		Data: DataBase64,
		Signature: SignBase64,
	})

	return buf.String(), SignBase64
}


func api(Data map[string]interface{}) []byte {
	DataBytes, _ := json.Marshal(Data)

	DataBase64 := makeData(string(DataBytes))
	SignBase64 := makeSignature(DataBase64)

	form := url.Values{
		"data":      {DataBase64},
		"signature": {SignBase64},
	}

	body := bytes.NewBufferString(form.Encode())
	rsp, err := http.Post(LiqpayURL + "request", "application/x-www-form-urlencoded", body)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Liqpay response:",string(bodyByte))
	return bodyByte
}


func makeData (Data string) (string) {
	return b64.StdEncoding.EncodeToString([]byte(Data))
}


func makeSignature (DataBase64 string) (string) {
	hasher := sha1.New()
	hasher.Write([]byte(privateKey))
	hasher.Write([]byte(DataBase64))
	hasher.Write([]byte(privateKey))
	SignBase64 := b64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return SignBase64
}