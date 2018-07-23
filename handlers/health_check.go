package handlers

import "my_errors"


func HealthCheck(inputData map[string]interface{}) (string, error) {
	return my_errors.SuccessfullyOperation()
}
