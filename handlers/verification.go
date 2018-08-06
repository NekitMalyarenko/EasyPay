package handlers

import (
	"db"
	"time"
	"my_errors"
	"encoding/json"
	"services"
	"log"
	"strconv"
)


func StartVerification(inputData map[string]interface{}) (string, error) {

	if inputData["user"] == nil {
		ver := db.GetInstance().Verification.GetVerification(inputData["phone_number"].(string))

		if ver != nil {
			return my_errors.GetError(my_errors.VerificationExists)
		} else {
			verificationCode := services.GetRandom(10000, 99999)

			err := services.SendSMS(strconv.FormatInt(verificationCode, 10), inputData["phone_number"].(string))
			if err != nil {
				log.Println(err)
			}

			err = db.GetInstance().Verification.CreateVerification(inputData["phone_number"].(string),
				verificationCode, time.Now().Format("Mon Jan _2 15:04:05 MST 2006"))
			if err != nil {
				return my_errors.GetError(my_errors.DBError)
			} else {
				return my_errors.SuccessfullyOperation()
			}
		}
	} else {
		return my_errors.GetError(my_errors.UserAlreadyExists)
	}
}


func VerifyPhone(inputData map[string]interface{}) (string, error) {

	if inputData["user"] == nil {
		data := db.GetInstance().Verification.GetVerification(inputData["phone_number"].(string))
		if data != nil && data.IsVerified == false {
			isVerified := false

			parsedTime, err := time.Parse("Mon Jan _2 15:04:05 MST 2006", data.StartTime)
			if err != nil {
				return my_errors.GetError(my_errors.DateParsingError)
			}

			elapsedTime := time.Now().Sub(parsedTime)
			if elapsedTime.Minutes() > 5 {
				db.GetInstance().DeleteVerification(inputData["phone_number"].(string))
				return my_errors.GetError(my_errors.DateParsingError)
			} else {
				parsedVerificationCode := int64(inputData["verification_code"].(float64))
				if data.VerificationCode == parsedVerificationCode {
					err = db.GetInstance().Verify(inputData["phone_number"].(string))
					if err != nil {
						return my_errors.GetError(my_errors.DBError)
					}
					isVerified = true
				}

				response, _ := json.Marshal(map[string]interface{}{
					"is_verified": isVerified,
					"status":      "ok",
				})
				return string(response), nil
			}

		} else {
			return my_errors.GetError(my_errors.VerificationNotExists)
		}
	} else {
		return my_errors.GetError(my_errors.UserAlreadyExists)
	}
}