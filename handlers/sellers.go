package handlers

import (
	"db"
	"my_errors"
	"types"
	"encoding/json"
	"log"
)


func SellerRegister(inputData map[string]interface{}) (string, error) {
	isVerified, err := db.GetInstance().Verification.IsVerified(inputData["phone_number"].(string))
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	if isVerified {
		seller := types.Seller{
			FirstName:   inputData["first_name"].(string),
			LastName:    inputData["last_name"].(string),
			PhoneNumber: inputData["phone_number"].(string),
			Description: inputData["description"].(string),
			Password:    inputData["password"].(string),
		}

		err := db.GetInstance().Sellers.AddSeller(&seller)
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		} else {
			db.GetInstance().Verification.DeleteVerification(inputData["phone_number"].(string))
			return my_errors.SuccessfullyOperation()
		}
	} else {
		return my_errors.GetError(my_errors.UnverifiedPhone)
	}
}


func GetSeller(inputData map[string]interface{}) (string, error) {
	sellerId := int64(inputData["seller_id"].(float64))

	seller, err := db.GetInstance().Sellers.GetSellerById(sellerId)
	if err != nil {
		log.Println(err)
		return my_errors.GetError(my_errors.DBError)
	}

	if seller != nil {
		response, _ := json.Marshal(map[string]interface{}{
			"status" : "ok",
			"first_name" : seller.FirstName,
			"last_name" : seller.LastName,
			"description" : seller.Description,
			"image" : seller.Image,
		})

		return string(response), nil
	} else {
		return my_errors.GetError(my_errors.UserNotFound)
	}
}