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

		user := types.User{
			FirstName : inputData["first_name"].(string),
			LastName : inputData["last_name"].(string),
			Email: inputData["email"].(string),
			Password : inputData["password"].(string),
		}

		err :=  db.GetInstance().Users.AddUser(&user)
		if err != nil {
			log.Println(err)
			return errors[2]
		} else {
			return successfullyOperation()
		}

	} else {
		return errors[1]
	}
}


func UserLogin(inputData map[string]interface{}) string {

	if inputData["email"] != nil && inputData["password"] != nil {
		email := inputData["email"].(string)
		password := inputData["password"].(string)

		user := db.GetInstance().GetUser(email)
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
				return errors[5]
			} else {
				return string(response)
			}
		} else {
			return errors[4]
		}
	} else {
		return errors[3]
	}

}


//email, password, image[]byte
func UserAddImage(inputData map[string]interface{}) string {



	return ""
}


//input:
// 	email
//output:
//	has bool
func HasUserWithEmail(inputData map[string]interface{}) string {
	return ""
}
