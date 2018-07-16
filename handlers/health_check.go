package handlers

import "log"


func Health(inputData map[string]interface{}) string{
	result:= inputData["temp"].(string)
	if result != "" {
		return result
	} else {
		log.Println("Nothing in temp")
		return ""
	}
}
