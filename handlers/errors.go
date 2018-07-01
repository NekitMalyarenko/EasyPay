package handlers

import "encoding/json"


const(
	argumentsError      = 1
	dbError             = 2
	jsonMarshalError    = 3
	emailError          = 4
	authenticationError = 5
)

var errors = map[int]string{
	1 : errorToJson(1, "not enough arguments"),
	2 : errorToJson(2, "db error"),
	3 : errorToJson(3, "json marshaling error"),
	4 : errorToJson(4, "email is already registered"),
	5 : errorToJson(5, "incorrect password or email"),
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
