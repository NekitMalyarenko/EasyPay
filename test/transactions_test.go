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
			"phone_number" : "380671779808",
			"password" : "hottabych98",
			"transaction" : map[string] interface{}{
				"products" : []map[string]interface{}{
					{
						"product_id" : 2,
						"quantity" : 1,
					},
					{
						"product_id" : 1,
						"quantity" : 2,
					},
					{
						"product_id" : 4,
						"quantity" : 1,
					},
				},
			},
			"shop_id" : 1,
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
			"phone_number": "380671779808",
			"password": "hottabych98",
			"account_type" : 0,
			"number_of_transactions": 2,
			"transaction_id" : 9,
		},
		path: "/",
	}
	sendRequest(requestData, t)
}


func TestGetNewestTransactions(t *testing.T){
	initHandlers()
	requestData:= requestData{
		methodName:"getNewestTransactions",
		requestData:map[string]interface{}{
			"phone_number": "380671779808",
			"password": "hottabych98",
			"account_type" : 0,
			"last_transaction_id" : 7,
		},
		path: "/",
	}
	sendRequest(requestData, t)
}


func TestGetTransactionById(t *testing.T){
	initHandlers()
	requestData:= requestData{
		methodName:"getTransactionById",
		requestData:map[string]interface{}{
			"phone_number": "380671779808",
			"password": "hottabych98",
			"transaction_id" : 7,
			"account_type" : 0,
		},
		path: "/",
	}
	res:= sendRequest(requestData, t)
	log.Println(res)
}