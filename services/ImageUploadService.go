package services

import (
	"golang.org/x/net/context"
	"github.com/kyokomi/cloudinary"
	"io/ioutil"
	"bytes"
	"log"
)

func UploadImage(image []byte, fileName string) string {
	ctx := context.Background()
	ctx = cloudinary.NewContext(ctx, "cloudinary://245738261838881:lSLutX6LmWZKc4hfYPENoMUgCGg@dbogdiydy")

	cloudinary.UploadStaticImage(ctx, fileName, bytes.NewBuffer(image))
	return ""
}


func TestUpload() {
	data, err := ioutil.ReadFile("test.png")
	if err != nil {
		log.Println(err)
	}

	UploadImage(data, "test")
}
