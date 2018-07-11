package handlers

import "encoding/json"


const(
	argumentsError           = 1
	dbError                  = 2
	jsonMarshalError         = 3
	phoneNumberError         = 4
	authenticationError      = 5
	verificationTimeOutError = 6
	dateParsingError         = 7
	verificationExists       = 8
	verificationNotExists    = 9
	unverifiedPhone          = 10
	parsingError             = 11
	userNotFound             = 12
)

var myErrors = map[int]string{
	1 : errorToJson(1, "not enough arguments"),
	2 : errorToJson(2, "db error"),
	3 : errorToJson(3, "json marshaling error"),
	4 : errorToJson(4, "phoneNumber is already registered"),
	5 : errorToJson(5, "incorrect password or email"),
	6 : errorToJson(6, "time of verification has ended"),
	7 : errorToJson(7, "date parsing error"),
	8 : errorToJson(8, "verification is already exists"),
	9 : errorToJson(9, "verification is not exists"),
	10 : errorToJson(10, "phone is not verified"),
	11 : errorToJson(11, "parsing error"),
	12 : errorToJson(12, "user not found"),
}


func errorToJson(errorCode int, errorDescription string) string {
	result, _ := json.Marshal(map[string]interface{}{
		"status" : "error",
		"error" : map[string]interface{}{
			"errorCode" : errorCode,
			"errorDescription" : errorDescription,
		},
	})

	return string(result)
}


func successfullyOperation() string {
	result, _ := json.Marshal(map[string]interface{}{
		"status" : "ok",
	})

	return string(result)
}
