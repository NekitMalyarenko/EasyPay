package handlers

import "encoding/json"

var errors = map[int]string{
	1 : errorToJson(1),
	2 : errorToJson(2),
	3 : errorToJson(3),
	4 : errorToJson(4),
	5 : errorToJson(5),
}


func errorToJson(errorCode int) string {
	result, _ := json.Marshal(map[string]interface{}{
		"status" : "error",
		"errorCode" : errorCode,
	})

	return string(result)
}


func successfullyOperation() string {
	result, _ := json.Marshal(map[string]interface{}{
		"status" : "ok",
	})

	return string(result)
}
