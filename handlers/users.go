package handlers

import (
	"db"
	"log"
	"encoding/json"
)


func UserLogin(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil && inputData["password"] != nil {
		customer, err := db.GetInstance().Customers.GetCustomer(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if customer != nil {

			if customer.Password == inputData["password"].(string) {
				response, _ := json.Marshal(map[string]interface{}{
					"status" : "ok",
					"first_name" : customer.FirstName,
					"last_name" : customer.LastName,
					"image" : customer.Image,
					"account_type" : 0,
				})
				return string(response)
			} else {
				return myErrors[authenticationError]
			}

		} else {
			seller, err := db.GetInstance().Sellers.GetSeller(inputData["phone_number"].(string))
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			if seller != nil {
				if seller.Password == inputData["password"].(string) {
					response, _ := json.Marshal(map[string]interface{}{
						"status" : "ok",
						"first_name" : seller.FirstName,
						"last_name" : seller.LastName,
						"image" : seller.Image,
						"description" : seller.Description,
						"account_type" : 1,
					})
					return string(response)
				} else {
					return myErrors[authenticationError]
				}
			} else {
				return myErrors[userNotFound]
			}
		}

	} else {
		return myErrors[argumentsError]
	}
}

