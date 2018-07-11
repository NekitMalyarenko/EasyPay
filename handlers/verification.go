package handlers

import (
	"db"
	"log"
	"math/rand"
	"encoding/json"
	"time"
	"strconv"
)


func StartVerification(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil {
		ver := db.GetInstance().Verification.GetVerification(inputData["phone_number"].(string))
		hasUser, err := db.GetInstance().Customers.HasCustomer(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if ver != nil || hasUser {
			return myErrors[verificationExists]
		} else {
			err := db.GetInstance().Verification.CreateVerification(inputData["phone_number"].(string),
				rand.Int63n(100000),time.Now().Format("Mon Jan _2 15:04:05 MST 2006"))
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			} else {
				return successfullyOperation()
			}
		}
	} else {
		return myErrors[argumentsError]
	}

}


func VerifyPhone(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil && inputData["verification_code"] != nil {
		data := db.GetInstance().Verification.GetVerification(inputData["phone_number"].(string))
		if data != nil && data.IsVerified == false {
			isVerified := false

			parsedTime, err := time.Parse("Mon Jan _2 15:04:05 MST 2006", data.StartTime)
			if err != nil {
				log.Println(err)
				return myErrors[dateParsingError]
			}

			//log.Println("parsed time:", parsedTime.String(), "current time:", time.Now().String())

			elapsedTime := time.Now().Sub(parsedTime)
			log.Println("elapsedTime:", elapsedTime.String())
			log.Println("1:", elapsedTime.Hours())
			if elapsedTime.Minutes() > 5 || (elapsedTime.Minutes() == 5 && elapsedTime.Seconds() > 0) {
				db.GetInstance().DeleteVerification(inputData["phone_number"].(string))
				return myErrors[verificationTimeOutError]
			} else {
				parsedVerificationCode, err := strconv.ParseInt(inputData["verification_code"].(string), 10, 64)
				if err != nil {
					log.Println(err)
					return myErrors[jsonMarshalError]
				}

				if data.VerificationCode == parsedVerificationCode {
					err = db.GetInstance().Verify(inputData["phone_number"].(string))
					if err != nil {
						log.Println(err)
						return myErrors[dbError]
					}
					isVerified = true
				}

				response, _ := json.Marshal(map[string]interface{}{
					"is_verified": isVerified,
					"status":      "ok",
				})
				return string(response)
			}
		} else {
			return myErrors[verificationNotExists]
		}
	} else {
		return myErrors[argumentsError]
	}
}