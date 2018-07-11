package handlers

import (
	"db"
	"log"
	"types"
	"encoding/json"
)


func SellerRegister(inputData map[string]interface{}) string {

	if inputData["first_name"] != nil && inputData["last_name"] != nil &&
		inputData["phone_number"] != nil && inputData["password"] != nil {

		hasUser, err := db.GetInstance().Sellers.HasSeller(inputData["phone_number"].(string))
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
				user := types.Seller{
					FirstName:   inputData["first_name"].(string),
					LastName:    inputData["last_name"].(string),
					PhoneNumber: inputData["phone_number"].(string),
					Description: inputData["description"].(string),
					Password:    inputData["password"].(string),
				}

				err := db.GetInstance().Sellers.AddSeller(&user)
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


func SellerLogin(inputData map[string]interface{}) string {

	if inputData["email"] != nil && inputData["password"] != nil {
		phoneNumber := inputData["phone_number"].(string)
		password := inputData["password"].(string)

		user, err := db.GetInstance().Sellers.GetSeller(phoneNumber)
		if err != nil {
			log.Println("error:", err)
			return myErrors[dbError]
		}

		if user != nil && user.Password == password {
			rowResponse := make(map[string]interface{})
			rowResponse["first_name"] = user.FirstName
			rowResponse["last_name"] = user.LastName
			rowResponse["description"] = user.Description
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