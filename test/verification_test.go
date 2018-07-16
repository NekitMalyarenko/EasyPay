package test

import (
	"testing"
	"log"
)

func TestCreateVerification(t *testing.T) {
	initHandlers()

	requestData := requestData{
		methodName  : "startVerification",
		requestData : map[string]interface{} {
			"phone_number" : "380967519035",
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}


func TestVerify(t *testing.T) {
	initHandlers()

	requestData := requestData{
		methodName  : "verifyPhone",
		requestData : map[string]interface{} {
			"phone_number" : "380967519035",
			"verification_code" : 61807,
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}