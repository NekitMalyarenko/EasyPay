package test

import (
	"testing"
	"log"
)

func TestUserRegister(t *testing.T) {
	initHandlers()
	requestData := requestData{
		methodName  : "customerRegister",
		requestData : map[string]interface{} {
			"first_name" : "Nikita2",
			"last_name" : "Maliarenko2",
			"phone_number" : "380967519035",
			"password" : "123456",
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}


func TestUserLogin(t *testing.T) {
	initHandlers()
	requestData := requestData{
		methodName  : "userLogin",
		requestData : map[string]interface{} {
			"phone_number" : "380967519036",
			//"password" : "1234",
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}