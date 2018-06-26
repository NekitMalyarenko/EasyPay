package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"errors"
	"handlers"
)

var methodHandlers map[string]func(input map[string]interface{})string


func main() {
	initHandlers()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}


func handler(w http.ResponseWriter, r *http.Request) {
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("error"))
	}

	var parsedInput map[string]interface{}
	err = json.Unmarshal(input, parsedInput)
	if err != nil {
		log.Println(err)
		w.Write([]byte("json parsing error"))
	}

	if parsedInput["method_name"] != nil {
		methodName := parsedInput["method_name"].(string)
		var methodData map[string]interface{}
		err = json.Unmarshal([]byte(parsedInput["method_data"].(string)), methodData)
		if err != nil {
			log.Println(err)
			w.Write([]byte("error parsing method data"))
		}

		if methodHandlers[methodName] != nil {
			w.Write([]byte(methodHandlers[methodName](methodData)))
		} else {
			log.Println(errors.New("no such method error"))
			w.Write([]byte("no such method error"))
		}
	} else {
		log.Println(errors.New("method_name is empty error"))
		w.Write([]byte("method_name is empty error"))
	}

}


func initHandlers() {
	methodHandlers = make(map[string]func(input map[string]interface{})string)
	methodHandlers["test"] = handlers.TempHandler
}
