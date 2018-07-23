package my_errors

import (
	"encoding/json"
	"errors"
)

type myErrorsType map[int]error

const(
	ArgumentsError           = 1
	DBError                  = 2
	JsonMarshalError         = 3
	PhoneNumberError         = 4
	AuthenticationError      = 5
	VerificationTimeOutError = 6
	DateParsingError         = 7
	VerificationExists       = 8
	VerificationNotExists    = 9
	UnverifiedPhone          = 10
	ParsingError             = 11
	UserNotFound             = 12
	CantReadBody             = 13
	JsonUnmarshalError       = 14
	MethodNameIsEmpty        = 15
	NoSuchMethodError        = 16
	UserAlreadyExists        = 17
	WrongAccountType         = 18
	SellerAlreadyAdded       = 19
)

var Errors myErrorsType


func init() {
	Errors = map[int]error {
		1 : errorToJson(1, "not enough arguments"),
		2 : errorToJson(2, "db error"),
		3 : errorToJson(3, "json marshaling error"),
		4 : errorToJson(4, "phoneNumber is already registered"),
		5 : errorToJson(5, "incorrect password or phone number"),
		6 : errorToJson(6, "time of verification has ended"),
		7 : errorToJson(7, "date parsing error"),
		8 : errorToJson(8, "verification is already exists"),
		9 : errorToJson(9, "verification is not exists"),
		10 : errorToJson(10, "phone is not verified"),
		11 : errorToJson(11, "parsing error"),
		12 : errorToJson(12, "user not found"),
		13 : errorToJson(13, "can't read request body"),
		14 : errorToJson(14, "json unmarshal error"),
		15 : errorToJson(15, "method name is empty"),
		16 : errorToJson(16, "no such method"),
		17 : errorToJson(17, "user already exists"),
		18 : errorToJson(18, "wrong account type"),
		19 : errorToJson(19,"there's already seller with such id"),
	}
}



/*var Errors = map[int]error {
	1 : errorToJson(1, "not enough arguments"),
	2 : errorToJson(2, "db error"),
	3 : errorToJson(3, "json marshaling error"),
	4 : errorToJson(4, "phoneNumber is already registered"),
	5 : errorToJson(5, "incorrect password or phone number"),
	6 : errorToJson(6, "time of verification has ended"),
	7 : errorToJson(7, "date parsing error"),
	8 : errorToJson(8, "verification is already exists"),
	9 : errorToJson(9, "verification is not exists"),
	10 : errorToJson(10, "phone is not verified"),
	11 : errorToJson(11, "parsing error"),
	12 : errorToJson(12, "user not found"),
	13 : errorToJson(13, "can't read request body"),
	14 : errorToJson(14, "json unmarshal error"),
	15 : errorToJson(15, "method name is empty"),
	16 : errorToJson(16, "no such method"),
	17 : errorToJson(17, "user already exists"),
}*/


func errorToJson(errorCode int, errorDescription string) error {
	result, _ := json.Marshal(map[string]interface{}{
		"status" : "error",
		"error" : map[string]interface{}{
			"errorCode" : errorCode,
			"errorDescription" : errorDescription,
		},
	})

	return errors.New(string(result))
}


func SuccessfullyOperation() (string, error) {
	result, _ := json.Marshal(map[string]interface{}{
		"status" : "ok",
	})

	return string(result), nil
}


func GetError(id int) (string, error) {
	return "", Errors[id]
}