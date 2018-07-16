package handlers

import (
	"db"
	"encoding/json"
	"types"
	"time"
	"my_errors"
	"services"
)


const(
	CustomerAccount = 0
	SellerAccount   = 1
)


type jsonTransaction struct {
	Id               int64         `json:"id"`
	UserId           int64         `json:"user_id"`
	ShopId           int64         `json:"shop_id"`
	Products         []interface{} `json:"products"`
	TotalPrice       int           `json:"total_price"`
	Date             string        `json:"date"`
	VerificationCode int64         `json:"verification_code"`
}


func AddTransaction(inputData map[string]interface{}) (string, error) {
	customer := inputData["user"].(*types.User).Customer

	if customer != nil {
		rowTransaction := inputData["transaction"].(map[string]interface{})
		shopId := int64(inputData["shop_id"].(float64))

		marshaledProducts, err := json.Marshal(rowTransaction["products"].([]interface{}))
		if err != nil {
			return my_errors.GetError(my_errors.JsonMarshalError)
		}

		transaction := types.Transaction{
			UserId:     customer.Id,
			ShopId:     shopId,
			Products:   string(marshaledProducts),
			TotalPrice: countTotalPrice(rowTransaction["products"].([]interface{})),
			Date:       time.Now().Format("02.01.2006 15:04:05 -0700"),
			VerificationCode : services.GetRandom(10000, 99999),
		}

		err = db.GetInstance().Transactions.AddTransaction(transaction)
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}

		return my_errors.SuccessfullyOperation()

	} else {
		return my_errors.GetError(my_errors.WrongAccountType)
	}
}


/*
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

			dbTransactions, err := db.GetInstance().Transactions.GetOldestTransactions(userId, transactionId,
				numberOfTransactions)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			transactions, err := parseDBTransactions(dbTransactions)
			if err != nil {
				log.Println(err)
				return myErrors[jsonMarshalError]
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
			dbTransactions, err := db.GetInstance().Transactions.GetNewestTransactions(lastTransactionId, userId)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			transactions, err := parseDBTransactions(dbTransactions)
			if err != nil {
				log.Println(err)
				return myErrors[jsonMarshalError]
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

			dbTransaction, err := db.GetInstance().Transactions.GetTransactionById(userId, transactionId)
			if err != nil {
				log.Println(err)
				return myErrors[dbError]
			}

			var transaction = &jsonTransaction{}
			if err = transaction.fromDBTransaction(dbTransaction); err != nil {
				log.Println(err)
				return myErrors[jsonMarshalError]
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
}*/


func countTotalPrice(products []interface{}) int {
	/*var totalPrice int
	for id, number := range products {
		currentItem,_ := db.GetInstance().Products.GetPrice(id)
		totalPrice +=currentItem*number
	}*/
	return 100
}


func parseDBTransactions(dbTransactions []*types.Transaction) ([]*jsonTransaction, error) {
	var transactions = make([]*jsonTransaction, len(dbTransactions))
	for position, val := range dbTransactions {
		var jsonTransaction  = &jsonTransaction{}
		if err := jsonTransaction.fromDBTransaction(val); err != nil {
			return nil, err
		}
		transactions[position] = jsonTransaction
	}

	return transactions, nil
}


func (jsonTransaction *jsonTransaction) fromDBTransaction(dbTransaction *types.Transaction) error {
	jsonTransaction.Id = dbTransaction.Id
	jsonTransaction.UserId = dbTransaction.UserId
	jsonTransaction.ShopId = dbTransaction.ShopId
	jsonTransaction.TotalPrice = dbTransaction.TotalPrice
	jsonTransaction.Date = dbTransaction.Date
	jsonTransaction.VerificationCode = dbTransaction.VerificationCode
	return json.Unmarshal([]byte(dbTransaction.Products), &jsonTransaction.Products)
}