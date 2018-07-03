package handlers

import (
	"db"
	"log"
	"math/rand"
	"encoding/json"
)


func StartVerification(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil {
		err := db.GetInstance().
			Verification.CreateVerification(inputData["phone_number"].(string),
				rand.Int63n(100000))
		if err != nil {
			log.Println(err)
			return errors[dbError]
		}

		return successfullyOperation()
	} else {
		return errors[argumentsError]
	}

}


func VerifyPhone(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil && inputData["verification_code"] != nil {

		data, err := db.GetInstance().Verification.GetVerification(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return errors[dbError]
		}

		isVerified := false

		if data.VerificationCode == inputData["verification_code"].(int64) {
			err = db.GetInstance().Verify(inputData["phone_number"].(string))
			if err != nil {
				log.Println(err)
				return errors[dbError]
			}
			isVerified = true
		}

		response, _ := json.Marshal(map[string]interface{}{
			"is_verified" : isVerified,
			"status"      : "ok",
		})
		return string(response)

	} else {
		return errors[argumentsError]
	}

}