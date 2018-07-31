package test

import (
	"testing"
	"log"
)


func TestShopAdd(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "shopRegister",
		requestData : map[string]interface{}{
			"phone_number" : "380971779808",
			"password" : "dzhamshut4",
			"shop" : map[string] interface{}{
				"name" : "shaurma4you",
				"email" : "test@emial.com",
				"description" : "Totally new kind of shaurma shop...",
				"cardNumber" : 12345,
				"products" : []map[string]interface{}{
					{
						"id" : 1,
						"name" : "item 1",
						"price" : 101,
						"image" : "",
					},
					{
						"id" : 3,
						"name" : "item 3",
						"price" : 103,
						"image" : "",
					},
				},
			},
		},
		path : "/",
	}
	sendRequest(requestData, t)
}


func TestShopAddProducts(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "shopAddProducts",
		requestData : map[string]interface{}{
			"phone_number" : "380971779808",
			"password" : "dzhamshut4",
			"products" : []map[string]interface{}{
					{
						"id" : 4,
						"name" : "item 4",
						"price" : 104,
						"image" : nil,
					},
					{
						"id" : 5,
						"name" : "item 5",
						"price" : 105,
						"image" : nil,
					},
					{
						"id" : 6,
						"name" : "item 6",
						"price" : 106,
						"image" : nil,
					},
				},
		},
		path : "/",
	}
	sendRequest(requestData, t)
}


func TestShopAddSellers(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "shopAddSeller",
		requestData : map[string]interface{}{
			"phone_number" : "380971779808",
			"password" : "dzhamshut4",
			"seller": 3 ,
		},
		path : "/",
	}
	sendRequest(requestData, t)
}


func TestGetShop(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "getShop",
		requestData : map[string]interface{}{
			"shop_id": 23,
		},
		path : "/",
	}

	log.Println(sendRequest(requestData, t))
}


func TestGetShopCardNumber(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "getShopCardNumber",
		requestData : map[string]interface{}{
			"phone_number" : "380671779808",
			"password" : "hottabych98",
			"shop_id" : 22,
		},
		path : "/",
	}
	sendRequest(requestData, t)
}


func TestGetShopProducts(t *testing.T){
	initHandlers()
	requestData := requestData{
		methodName : "getShopProducts",
		requestData : map[string]interface{}{
			"shop_id" : 23,
		},
		path : "/",
	}
	sendRequest(requestData, t)
}