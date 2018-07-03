package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"errors"
	"handlers"
	"services"
)

var methodHandlers map[string]func(input map[string]interface{})string


func main() {
	/*initHandlers()

	db.GetInstance()

	http.HandleFunc("/", Handle)
	http.ListenAndServe(":8080", nil)*/

	services.SendSMS("test", "380967519036")
}


func Handle(w http.ResponseWriter, r *http.Request) {
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("error"))
		return
	}

	var parsedInput map[string]interface{}
	err = json.Unmarshal(input, &parsedInput)
	if err != nil {
		log.Println(err)
		w.Write([]byte("json parsing error"))
		return
	}

	log.Println("imput data:", parsedInput)

	if parsedInput["method_name"] != nil {
		methodName := parsedInput["method_name"].(string)
		var methodData = parsedInput["method_data"].(map[string]interface{})

		if methodHandlers[methodName] != nil {
			response := methodHandlers[methodName](methodData)
			log.Println("response:", response)
			w.Write([]byte(response))
		} else {
			log.Println(errors.New("no such method error"))
			w.Write([]byte("no such method error"))
			return
		}
	} else {
		log.Println(errors.New("method_name is empty error"))
		w.Write([]byte("method_name is empty error"))
		return
	}

	//w.WriteHeader(200)
}


func initHandlers() {
	methodHandlers = make(map[string]func(input map[string]interface{})string)
	methodHandlers["healthTest"] = handlers.TempHandler
	methodHandlers["userLogin"] = handlers.UserLogin
	methodHandlers["userRegister"] = handlers.UserRegister
	methodHandlers["hasUserWithEmail"] = handlers.HasUserWithEmail
}