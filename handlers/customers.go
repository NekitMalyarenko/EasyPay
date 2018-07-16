package handlers

import (
	"db"
	"my_errors"
	"types"
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


//email, password, image[]byte
func CustomerAddImage(inputData map[string]interface{}, image []byte) string {

	/*if inputData["phone_number"] != nil && inputData["password"] != nil && len(image) > 0{
		customer, err := db.GetInstance().Customers.GetCustomer(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if customer != nil && customer.Password == inputData["password"].(string) {

			//services.UploadImage()

		} else {
			return myErrors[authenticationError]
		}

	} else {
		return myErrors[argumentsError]
	}*/

	return ""
}