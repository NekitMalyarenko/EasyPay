package handlers

import (
	"types"
	"db"
	"log"
	"encoding/json"
)

//Without image
func CustomerRegister(inputData map[string]interface{}) string {

	if inputData["first_name"] != nil && inputData["last_name"] != nil &&
		inputData["phone_number"] != nil && inputData["password"] != nil {

		hasUser, err := db.GetInstance().Customers.HasCustomer(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if !hasUser {
			isVerified, err := db.GetInstance().Verification.IsVerified(inputData["phone_number"].(string))
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
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
					log.Println(err)
					return myErrors[dbError]
				} else {
					db.GetInstance().Verification.DeleteVerification(inputData["phone_number"].(string))
					return successfullyOperation()
				}
			} else {
				return myErrors[unverifiedPhone]
			}
		} else {
			return myErrors[phoneNumberError]
		}

	} else {
		return myErrors[argumentsError]
	}
}


func CustomerLogin(inputData map[string]interface{}) string {

	if inputData["email"] != nil && inputData["password"] != nil {
		email := inputData["email"].(string)
		password := inputData["password"].(string)

		user, err := db.GetInstance().Customers.GetCustomer(email)
		if err != nil {
			log.Println("error:", err)
			return myErrors[dbError]
		}

		if user != nil && user.Password == password {
			rowResponse := make(map[string]interface{})
			rowResponse["first_name"] = user.FirstName
			rowResponse["last_name"] = user.LastName
			//rowResponse["email"] = user.Email
			rowResponse["image"] = user.Image
			rowResponse["status"] = "ok"

			response, err := json.Marshal(rowResponse)
			if err != nil {
				log.Println(err)
				return myErrors[jsonMarshalError]
			} else {
				return string(response)
			}
		} else {
			return myErrors[authenticationError]
		}
	} else {
		return myErrors[argumentsError]
	}

}


//email, password, image[]byte
func CustomerAddImage(inputData map[string]interface{}) string {



	return ""
}


func HasCustomerWithPhoneNumber(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil && inputData["phone_number"] != "" {
		hasUser, err := db.GetInstance().HasCustomer(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		response, err := json.Marshal(map[string]interface{}{
			"status" : "ok",
			"has" : hasUser,
		})
		if err != nil {
			return myErrors[jsonMarshalError]
		}

		return string(response)
	} else {
		return myErrors[argumentsError]
	}
}
