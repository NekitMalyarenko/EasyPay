package handlers

import (
	"services"
	"log"
	"os"
)

func TestPhoto(input map[string]interface{}, rawData []byte) string {
	photoName := input["photo_name"].(string)

	//log.Println("len:", string(rawData))

	file, err := os.Create("test.jpg")
	if err != nil {
		log.Println(err)
	}

	file.Write(rawData)
	file.Close()

	link := services.UploadImage(rawData, "user_images/" + photoName)
	log.Println("photo link:", link)

	return ""
}
