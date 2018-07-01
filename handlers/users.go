package handlers

import (
	"types"
	"db"
	"log"
	"encoding/json"
)

//Without image
func UserRegister(inputData map[string]interface{}) string {

	if inputData["first_name"] != nil && inputData["last_name"] != nil &&
		inputData["email"] != nil && inputData["password"] != nil {

		hasUser, err := db.GetInstance().HasUser(inputData["email"].(string))
		if err != nil {
			log.Println(err)
			return errors[dbError]
		}

		if !hasUser {
			user := types.User{
				FirstName : inputData["first_name"].(string),
				LastName : inputData["last_name"].(string),
				Email: inputData["email"].(string),
				Password : inputData["password"].(string),
			}

			err :=  db.GetInstance().Users.AddUser(&user)
			if err != nil {
				log.Println(err)
				return errors[dbError]
			} else {
				return successfullyOperation()
			}
		} else {
			return errors[emailError]
		}

	} else {
		return errors[argumentsError]
	}
}


func UserLogin(inputData map[string]interface{}) string {

	if inputData["email"] != nil && inputData["password"] != nil {
		email := inputData["email"].(string)
		password := inputData["password"].(string)

		user, err := db.GetInstance().GetUser(email)
		if err != nil {
			log.Println("error:", err)
			return errors[dbError]
		}

		if user != nil && user.Password == password {
			rowResponse := make(map[string]interface{})
			rowResponse["first_name"] = user.FirstName
			rowResponse["last_name"] = user.LastName
			rowResponse["email"] = user.Email
			rowResponse["image"] = user.Image
			rowResponse["status"] = "ok"

			response, err := json.Marshal(rowResponse)
			if err != nil {
				log.Println(err)
				return errors[jsonMarshalError]
			} else {
				return string(response)
			}
		} else {
			return errors[authenticationError]
		}
	} else {
		return errors[argumentsError]
	}

}


//email, password, image[]byte
func UserAddImage(inputData map[string]interface{}) string {



	return ""
}


func HasUserWithEmail(inputData map[string]interface{}) string {

	if inputData["email"] != nil && inputData["email"] != "" {
		hasUser, err := db.GetInstance().HasUser(inputData["email"].(string))
		if err != nil {
			log.Println(err)
			return errors[dbError]
		}

		response, err := json.Marshal(map[string]interface{}{
			"status" : "ok",
			"has" : hasUser,
		})
		if err != nil {
			return errors[jsonMarshalError]
		}

		return string(response)
	} else {
		return errors[argumentsError]
	}
}
