package handlers

import (
	"db"
	"my_errors"
	"types"
	"services"
	"encoding/json"
)

//Without image
func CustomerRegister(inputData map[string]interface{}) (string, error) {
	isVerified, err := db.GetInstance().Verification.IsVerified(inputData["phone_number"].(string))
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	if isVerified {
		user := types.Customer{
			FirstName:   inputData["first_name"].(string),
			LastName:    inputData["last_name"].(string),
			PhoneNumber: inputData["phone_number"].(string),
			Password:    inputData["password"].(string),
		}

		err := db.GetInstance().Customers.AddCustomer(&user)
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


func GetCustomer(inputData map[string]interface{}) (string, error) {
	userId := int64(inputData["customer_id"].(float64))

	customer, err := db.GetInstance().Customers.GetCustomerById(userId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status" : "ok",
		"id" : customer.Id,
		"first_name" : customer.FirstName,
		"last_name" : customer.LastName,
		"image" : customer.Image,
	})

	return string(response), nil
}


//email, password, image[]byte
func AddCustomerImage(inputData map[string]interface{}, image []byte) (string, error) {
	user := inputData["user"].(*types.User)
	var (
		id       int64
		fileName string
		err      error
	)

	if user.Customer != nil {
		id = user.Customer.Id
		fileName = user.Customer.FirstName + "_" + user.Customer.LastName + "_" + user.Customer.PhoneNumber
	} else if user.Seller != nil {
		id = user.Seller.Id
		fileName = user.Seller.FirstName + "_" + user.Seller.LastName + "_" + user.Seller.PhoneNumber
	}

	imageLink := services.UploadImage(image, fileName)
	if err = db.GetInstance().Customers.AddCustomerImage(id, imageLink); err != nil {
		return "", err
	} else {
		return my_errors.SuccessfullyOperation()
	}
}