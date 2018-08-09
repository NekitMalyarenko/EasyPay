package handlers

import (
	"encoding/json"
	"types"
	"my_errors"
)


func UserLogin(inputData map[string]interface{}) (string, error) {
	user := inputData["user"].(*types.User)
	var response []byte

	if user.Customer != nil && user.Seller == nil {
		response, _ = json.Marshal(map[string]interface{}{
			"status" : "ok",
			"first_name" : user.Customer.FirstName,
			"last_name" : user.Customer.LastName,
			"image" : user.Customer.Image,
			"account_type" : 0,
		})

	} else if user.Seller != nil && user.Customer == nil {
		response, _ = json.Marshal(map[string]interface{}{
			"status" : "ok",
			"first_name" : user.Seller.FirstName,
			"last_name" : user.Seller.LastName,
			"image" : user.Seller.Image,
			"description" : user.Seller.Description,
			"shop_id" : user.Seller.ShopId,
			"account_type" : 1,
		})
	}

	if response == nil {
		return my_errors.GetError(my_errors.UserNotFound)
	}else {
		return string(response), nil
	}

	/*if inputData["phone_number"] != nil && inputData["password"] != nil {
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
	}*/
}