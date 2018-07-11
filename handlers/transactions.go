package handlers

import (
	"db"
	"encoding/json"
	"types"
	"log"
	"time"
	"errors"
	"math/rand"
)


const(
	CustomerAccount = 0
	SellerAccount   = 1
)


func AddTransaction(inputData map[string]interface{}) string {
	if inputData["phone_number"] != nil && inputData["password"] != nil &&
		inputData["shop_id"] != nil && inputData["transaction"] != nil {

		user, err := db.GetInstance().Customers.GetCustomer(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return myErrors[authenticationError]
		}

		if user != nil && user.Password == inputData["password"].(string) {
			rowTransaction := inputData["transaction"].(map[string]interface{})
			shopId := int64(inputData["shop_id"].(float64))

			marshaledProducts, err := json.Marshal(rowTransaction["products"].([]interface{}))
			if err != nil {
				log.Println(err)
				return myErrors[jsonMarshalError]
			}

			transaction := types.Transaction{
				UserId:     user.Id,
				ShopId:     shopId,
				Products:   string(marshaledProducts),
				TotalPrice: countTotalPrice(rowTransaction["products"].([]interface{})),
				Date:       time.Now().Format("Mon Jan _2 15:04:05 MST 2006"),
				VerificationCode : rand.Int63n(1000000),
			}

			err = db.GetInstance().Transactions.AddTransaction(transaction)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			return successfullyOperation()
		} else {
			return myErrors[authenticationError]
		}
	} else {
		return myErrors[argumentsError]
	}
}


func GetOldestTransactions(inputData map[string]interface{}) string {
	if inputData["phone_number"] != nil && inputData["password"] != nil &&
		inputData["account_type"] != nil && inputData["transaction_id"] != nil &&
		inputData["number_of_transactions"] != nil{

		userId, ok, err := getUserId(inputData["phone_number"].(string), inputData["password"].(string),
			int(inputData["account_type"].(float64)))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if ok {
			transactionId := int64(inputData["transaction_id"].(float64))
			numberOfTransactions := int(inputData["number_of_transactions"].(float64))

			transactions, err := db.GetInstance().Transactions.GetOldestTransactions(userId, transactionId,
				numberOfTransactions)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			response, _ := json.Marshal(map[string]interface{}{
				"status" : "ok",
				"transactions" : transactions,
			})
			return string(response)
		} else {
			return myErrors[authenticationError]
		}
	} else {
		return myErrors[argumentsError]
	}
}


func GetNewestTransactions(inputData map[string]interface{}) string {

	if inputData["phone_number"] != nil && inputData["password"] != nil &&
		inputData["account_type"] != nil && inputData["last_transaction_id"] != nil {

		userId, ok, err := getUserId(inputData["phone_number"].(string), inputData["password"].(string),
			int(inputData["account_type"].(float64)))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if ok {
			lastTransactionId := int64(inputData["last_transaction_id"].(float64))

			transactions, err := db.GetInstance().Transactions.GetNewestTransactions(lastTransactionId, userId)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			response, _ := json.Marshal(map[string]interface{}{
				"status" : "ok",
				"transactions" : transactions,
			})
			return string(response)
		} else {
			return myErrors[authenticationError]
		}
	} else {
		return myErrors[argumentsError]
	}
}


func GetTransactionById(inputData map[string]interface{}) string {
	if inputData["phone_number"] != nil && inputData["password"] != nil &&
		inputData["account_type"] != nil && inputData["transaction_id"] != nil {

		userId, ok, err := getUserId(inputData["phone_number"].(string), inputData["password"].(string),
			int(inputData["account_type"].(float64)))
		if err != nil {
			log.Println(err)
			return myErrors[dbError]
		}

		if ok {
			transactionId := int64(inputData["transaction_id"].(float64))

			transaction, err := db.GetInstance().Transactions.GetTransactionById(userId, transactionId)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			response, _ := json.Marshal(map[string]interface{}{
				"status" : "ok",
				"transaction" : transaction,
			})
			return string(response)
		} else {
			return myErrors[authenticationError]
		}
	} else {
		return myErrors[argumentsError]
	}
}


func getUserId(phoneNumber, password string, rowAccountType int) (int64, bool, error)   {
	switch rowAccountType {
		case CustomerAccount:
			customer, err := db.GetInstance().Customers.GetCustomer(phoneNumber)
			if err != nil {
				return -1, false, err
			} else if customer.Password == password {
				return customer.Id, true, nil
			} else {
				return -1, false, nil
			}
		break

		case SellerAccount:
			seller, err := db.GetInstance().Sellers.GetSeller(phoneNumber)
			if err != nil {
				return -1, false, err
			} else if seller.Password == password {
				return seller.ShopId, true, nil
			} else {
				return -1, false, nil
			}
		break
	}

	return -1, false, errors.New("unknown account type")
}


func countTotalPrice(products []interface{}) int {
	/*var totalPrice int
	for id, number := range products {
		currentItem,_ := db.GetInstance().Products.GetPrice(id)
		totalPrice +=currentItem*number
	}*/
	return 100
}