package test

import "testing"


func TestAddLike(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "addLike",
		requestData : map[string]interface{}{
			"phone_number" : "380967519036",
			"password" : "1234",
			"shop_id" : 23,
		},
		path : "/",
	}
	sendRequest(requestData, t)
}


func TestAddDislike(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "addDislike",
		requestData : map[string]interface{}{
			"phone_number" : "380967519036",
			"password" : "1234",
			"shop_id" : 23,
		},
		path : "/",
	}
	sendRequest(requestData, t)
}