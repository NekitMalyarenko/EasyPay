package handlers

import (
	"db"
	"encoding/json"
	"types"
	"time"
	"my_errors"
	"services"
	"errors"
)


type jsonTransaction struct {
	Id               int64            `json:"id"`
	UserId           int64            `json:"user_id"`
	ShopId           int64            `json:"shop_id"`
	Products         []*jsonProduct `json:"products"`
	TotalPrice       int              `json:"total_price"`
	Date             string           `json:"date"`
	VerificationCode int64            `json:"verification_code"`
}


type jsonProduct struct {
	types.Product
	Quantity     int `json:"quantity"`
}


func AddTransaction(inputData map[string]interface{}) (string, error) {
	customer := inputData["user"].(*types.User).Customer

	if customer != nil {
		shop, err := db.GetInstance().GetShopById(int64(inputData["shop_id"].(float64)))
		if err != nil {
			return "", err
		}

		var rawProducts []interface{}
		err = json.Unmarshal([]byte(inputData["products"].(string)), &rawProducts)
		if err != nil {
			return "", err
		}

		products := productsFromString(rawProducts)

		transaction := types.Transaction{
			UserId:     customer.Id,
			ShopId:     shop.Id,
			Products:   productsToString(products),
			TotalPrice: countTotalPrice(shop, products),
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


func GetOldestTransactions(inputData map[string]interface{}) (string, error) {

	user := inputData["user"].(*types.User)
	var(
		dbTransactions []*types.Transaction
		err error
	)

	if user.Customer != nil && user.Seller == nil {
		dbTransactions, err = db.GetInstance().Transactions.GetOldestTransactions(user.Customer.Id,
			int64(inputData["transaction_id"].(float64)),
			int(inputData["number_of_transactions"].(float64)),false)
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}

	} else if user.Customer == nil && user.Seller != nil{
		dbTransactions, err = db.GetInstance().Transactions.GetOldestTransactions(user.Seller.ShopId,
			int64(inputData["transaction_id"].(float64)),
			int(inputData["number_of_transactions"].(float64)),true)
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}

	}

	transactions, err := parseDBTransactions(dbTransactions)
	if err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status":       "ok",
		"transactions": transactions,
	})

	return string(response), nil
}


func GetNewestTransactions(inputData map[string]interface{}) (string, error) {
	user := inputData["user"].(*types.User)
	var(
		dbTransactions []*types.Transaction
		err error
	)

	if user.Customer != nil && user.Seller == nil {
		dbTransactions, err = db.GetInstance().Transactions.GetNewestTransactions(int64(inputData["last_transaction_id"].(float64)),
			user.Customer.Id, false)
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}

	} else if user.Customer == nil && user.Seller != nil {
		dbTransactions, err = db.GetInstance().Transactions.GetNewestTransactions(int64(inputData["last_transaction_id"].(float64)),
			user.Seller.ShopId, true)
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}

	}

	transactions, err := parseDBTransactions(dbTransactions)
	if err != nil {
		return "", err
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status":       "ok",
		"transactions": transactions,
	})

	return string(response), nil
}


func GetTransactionById(inputData map[string]interface{}) (string, error) {

	transactionId := int64(inputData["transaction_id"].(float64))

	dbTransaction, err := db.GetInstance().Transactions.GetTransactionById(transactionId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	var transaction= &jsonTransaction{}
	if err = transaction.fromDBTransaction(dbTransaction); err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status":      "ok",
		"transaction": transaction,
	})

	return string(response), nil
}


func countTotalPrice(shop *types.Shop, products []jsonProduct) int {
	parsedShop, _ := parseDBShop(shop)

	var totalPrice int
	for _, rawProduct := range products {
		product := parsedShop.getProduct(rawProduct.Id)
		totalPrice += product.Price * rawProduct.Quantity
	}
	return totalPrice
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


func productsToString(products []jsonProduct) string {
	var rawRes = make([]map[string]interface{}, len(products))

	for index, val := range products {
		rawRes[index] = map[string]interface{}{
			"id" : val.Id,
			"quantity" : val.Quantity,
		}
	}

	res, _ := json.Marshal(rawRes)
	return string(res)
}


func productsFromString(rawProducts []interface{}) (products []jsonProduct) {
	products = make([]jsonProduct, len(rawProducts))

	for i, val := range rawProducts {
		parsedVal := val.(map[string]interface{})

		 jsonProduct := jsonProduct{
		 	Quantity : int(parsedVal["quantity"].(float64)),
		 }
		 jsonProduct.Id = int64(parsedVal["id"].(float64))

		products[i] = jsonProduct
	}

	return products
}


func (jsonTransaction *jsonTransaction) toString() string {
	res, _ := 	json.Marshal(jsonTransaction)
	return string(res)
}


func (jsonTransaction *jsonTransaction) fromDBTransaction(dbTransaction *types.Transaction) error {
	jsonTransaction.Id = dbTransaction.Id
	jsonTransaction.UserId = dbTransaction.UserId
	jsonTransaction.ShopId = dbTransaction.ShopId
	jsonTransaction.TotalPrice = dbTransaction.TotalPrice
	jsonTransaction.Date = dbTransaction.Date
	jsonTransaction.VerificationCode = dbTransaction.VerificationCode
	return jsonTransaction.getExtendedProducts(dbTransaction)
}


func (jsonTransaction *jsonTransaction) getExtendedProducts(dbTransaction *types.Transaction) error {
	err := json.Unmarshal([]byte(dbTransaction.Products), &jsonTransaction.Products)
	if err != nil {
		return err
	}

	dbShop, err := db.GetInstance().Shops.GetShopById(dbTransaction.ShopId)
	if err != nil {
		return err
	}

	shop := jsonShop{}
	shop.fromDBShop(dbShop)
	var getProduct = func(id int64) *types.Product {
		for _, val := range shop.RowProducts {
			if val.Id == id {
				return &val
			}
		}

		return nil
	}

	for _, val := range jsonTransaction.Products {
		fullProduct := getProduct(val.Id)
		if fullProduct != nil {
			val.Price = fullProduct.Price
			val.Name = fullProduct.Name
		} else {
			errors.New("can't find product")
		}
	}

	return nil
}