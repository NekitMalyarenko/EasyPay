package services

import (
	"golang.org/x/net/context"
	"github.com/kyokomi/cloudinary"
	"io/ioutil"
	"bytes"
	"log"
)

const path = "users_images"


func UploadImage(image []byte, fileName string) string {
	ctx := context.Background()
	ctx = cloudinary.NewContext(ctx, "cloudinary://245738261838881:lSLutX6LmWZKc4hfYPENoMUgCGg@dbogdiydy")

	cloudinary.UploadStaticImage(ctx, fileName, bytes.NewBuffer(image))
	return "https://res.cloudinary.com/dbogdiydy/image/upload/v1530021087/" + fileName
}


func TestUpload() {
	data, err := ioutil.ReadFile("test.jpg")
	if err != nil {
		log.Println(err)
	}

	UploadImage(data, path + "/test1")
}
