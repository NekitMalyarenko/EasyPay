package test

import (
	_"io/ioutil"
	"log"
	"encoding/json"
	"testing"
	"net/http"
	"bytes"
	"net/http/httptest"
)

type requestData struct {
	methodName  string
	requestData map[string]interface{}
	path        string
}



func TestHealth(t *testing.T) {
	initHandlers()

	requestData := requestData{
		methodName  : "healthTest",
		requestData : map[string]interface{} {
			"temp" : "pong",
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}


func TestUserRegister(t *testing.T) {
	requestData := requestData{
		methodName  : "userRegister",
		requestData : map[string]interface{} {
			"first_name" : "Nikita2",
			"last_name" : "Maliarenko2",
			"email" : "n.a.m.62608@gmail.com",
			"password" : "12345",
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}


func TestUserLogin(t *testing.T) {
	requestData := requestData{
		methodName  : "userLogin",
		requestData : map[string]interface{} {
			"email" : "n.a.m.62608@gmail.com",
			"password" : "1234",
		},
		path : "/",
	}
	res := sendRequest(requestData, t)
	log.Println(res)
}


func sendRequest(requestData requestData, t *testing.T) string {
	req, err := http.NewRequest("POST", requestData.path, bytes.NewBuffer([]byte(requestData.toString())))
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handle)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	return rr.Body.String()
}


func (data *requestData) toString() string{
	result, err := json.Marshal(map[string]interface{}{
		"method_name" : data.methodName,
		"method_data" : data.requestData,
	})
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}
