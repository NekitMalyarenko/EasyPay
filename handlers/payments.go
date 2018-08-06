package handlers

import (
	"db"
	"services"
	"my_errors"
	"log"
	"local_storage"
	"types"
	"time"
)


func Pay(inputData map[string]interface{}) (string, error) {
	customer := inputData["user"].(*types.User).Customer

	shop, err := db.GetInstance().GetShopById(int64(inputData["shop_id"].(float64)))
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	lastId, err := db.GetInstance().Transactions.GetBiggestId()
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	products := productsFromString(inputData["products"].([]interface{}))
	totalPrice :=  countTotalPrice(shop, products)

	log.Println("orderId:", lastId + 1)

	webPage, _ := services.GenerateP2P(map[string]interface{}{
		"orderId" : lastId + 1,
		"receiverCard" : shop.CardNumber,
		"amount" : totalPrice,
	})

	local_storage.GetInstance().Payments.PutPayment(&types.PaymentData{
		Id : lastId + 1,
		Transaction : types.Transaction{
			UserId:     customer.Id,
			ShopId:     shop.Id,
			Products:   productsToString(products),
			TotalPrice: countTotalPrice(shop, products),
			Date:       time.Now().Format("02.01.2006 15:04:05 -0700"),
			VerificationCode : services.GetRandom(10000, 99999),
		},
	})
	if err != nil {
		return "", err
	}

	return webPage, nil
}


func CheckAllPayments() {
	payments := local_storage.GetInstance().Payments.GetAllPayments()
	for index, val := range payments {
		status, err := services.GetPaymentStatus(val.Id)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("order_id:", index, "status:", status)

		if status == "sandbox" || status == "success" {
			err = db.GetInstance().Transactions.AddTransaction(val.Transaction)
			if err != nil {
				log.Println(err)
			}

			local_storage.GetInstance().Payments.DeletePayment(index)
		}
	}
}