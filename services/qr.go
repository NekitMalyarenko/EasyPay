package services
/*
import (
	"github.com/skip2/go-qrcode"
	"log"
	"strconv"
)


func CreateQR(shopId int64) ([]byte, error) {
	content := strconv.FormatInt(shopId,10)
	image, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		log.Println("Something wrong with qr code creation")
		return nil, err
	}

	return image, nil
}*/