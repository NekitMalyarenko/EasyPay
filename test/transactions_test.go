package test

import (
	"log"
	"testing"
)


func TestTransactionAdd(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "addTransaction",
		requestData : map[string]interface{}{
			"phone_number" : "380967519036",
			"password" : "1234",
			"transaction" : map[string] interface{}{
				"products" : []map[string]interface{}{
					{
						"id" : 1,
						"quantity" : 2,
					},
					{
						"id" : 3,
						"quantity" : 1,
					},
				},
			},
			"shop_id" : 23,
		},
		path : "/",
	}
	sendRequest(requestData, t)
}


func TestGetOldestTransactions(t *testing.T){
	initHandlers()
	requestData:= requestData{
		methodName : "getOldestTransactions",
		requestData : map[string]interface{}{
			"phone_number": "380967519036",
			"password": "1234",
			"number_of_transactions": 1,
			"transaction_id" : 20,
		},
		path: "/",
	}
	log.Println(sendRequest(requestData, t))
}


func TestGetNewestTransactions(t *testing.T){
	initHandlers()
	requestData:= requestData{
		methodName:"getNewestTransactions",
		requestData:map[string]interface{}{
			"phone_number": "380967519036",
			"password": "1234",
			"last_transaction_id" : 12,
		},
		path: "/",
	}
	res:=sendRequest(requestData, t)
	log.Println(res)
}


func TestGetTransactionById(t *testing.T){
	initHandlers()
	requestData:= requestData{
		methodName:"getTransactionById",
		requestData:map[string]interface{}{
			"phone_number": "380671779808",
			"password": "hottabych98",
			"transaction_id" : 17,
		},
		path: "/",
	}
	res:= sendRequest(requestData, t)
	log.Println(res)
}