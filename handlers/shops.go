package handlers

import (
	"db"
	"types"
	"encoding/json"
	"my_errors"
)


type jsonShop struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string 		    `json:"image"`
	Likes       int             `json:"likes"`
	Dislikes    int             `json:"dislikes"`
	CardNumber  int64           `json:"card_number"`
	Email       string          `json:"email"`
	RowSellers  string          `json:"sellers"`
	RowProducts []types.Product `json:"products"`
}


func ShopRegister(inputData map[string]interface{}) (string, error) {
	seller := inputData["user"].(*types.User).Seller
	inputShop := inputData["shop"].(map[string]interface{})

	products, err := json.Marshal(inputShop["products"].([]interface{}))
	if err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	parsedSeller, err := json.Marshal([]int64{seller.Id})
	if err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	shop := types.Shop{
		Name:        inputShop["name"].(string),
		Email:       inputShop["email"].(string),
		Description: inputShop["description"].(string),
		CardNumber:  int64(inputShop["cardNumber"].(float64)),
		RowProducts: string(products),
		RowSellers:  string(parsedSeller),
	}
	err = db.GetInstance().Shops.AddShop(shop)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	return my_errors.SuccessfullyOperation()
}


func ShopAddProducts(inputData map[string]interface{})(string, error){
	seller := inputData["user"].(*types.User).Seller

	dbShop, err := db.GetInstance().Shops.GetShopById(seller.ShopId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	shop, err := parseDBShop(dbShop)
	if err != nil {
		my_errors.GetError(my_errors.JsonUnmarshalError)
	}

	extendedProducts, err:= getExtendedProducts(inputData["products"].([]types.Product), shop.RowProducts)
	if err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	err = db.GetInstance().Shops.UpdateProducts(extendedProducts, shop.Id)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	return my_errors.SuccessfullyOperation()
}


func ShopAddSeller(inputData map[string]interface{})(string, error){
	seller := inputData["user"].(*types.User).Seller
	var sellers []int64

	shop, err := db.GetInstance().Shops.GetShopById(seller.ShopId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	err = json.Unmarshal([]byte(shop.RowSellers), &sellers)
	if err != nil {
		return my_errors.GetError(my_errors.JsonUnmarshalError)
	}

	newSellerId := int64(inputData["seller"].(float64))
	if !contains(sellers, newSellerId) {
		sellers = append(sellers, newSellerId)
		rowSellers, err := json.Marshal(sellers)
		if err != nil {
			return my_errors.GetError(my_errors.JsonMarshalError)
		}

		err = db.GetInstance().Shops.UpdateSellers(string(rowSellers), shop.Id)
		if err != nil {
			my_errors.GetError(my_errors.DBError)
		}
	} else {
		return my_errors.GetError(my_errors.SellerAlreadyAdded)
	}

	err = db.GetInstance().Sellers.UpdateShopId(int64(inputData["seller"].(float64)),shop.Id)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	return my_errors.SuccessfullyOperation()
}


func GetShop(inputData map[string]interface{}) (string, error) {
	shopId := int64(inputData["shop_id"].(float64))
	dbShop, err := db.GetInstance().Shops.GetShopById(shopId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	shop, err := parseDBShop(dbShop)
	if err != nil {
		my_errors.GetError(my_errors.JsonUnmarshalError)
	}

	response, err := json.Marshal(map[string]interface{}{
		"id":          shop.Id,
		"image":       shop.Image,
		"name":        shop.Name,
		"description": shop.Description,
		"likes":       shop.Likes,
		"dislikes":    shop.Dislikes,
		"email":       shop.Email,
		"sellers":     shop.RowSellers,
		"products":    shop.RowProducts,
		"status":      "ok",
	})
	if err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	return string(response),nil
}


func GetShopProducts(inputData map[string]interface{}) (string, error) {
	shopId := int64(inputData["shop_id"].(float64))
	dbShop, err := db.GetInstance().Shops.GetShopById(shopId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	shop, err := parseDBShop(dbShop)
	if err != nil {
		my_errors.GetError(my_errors.JsonUnmarshalError)
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status" : "ok",
		"products" : shop.RowProducts,
	})

	return string(response), nil
}


func GetShopCardNumber(inputData map[string]interface{}) (string,error) {
	shopId := int64(inputData["shop_id"].(float64))
	cardNumber, err := db.GetInstance().Shops.GetShopCardNumber(shopId)
	if err != nil {
		return my_errors.GetError(my_errors.DBError)
	}

	rowResponse, err := json.Marshal(map[string]interface{}{
		"status" : "ok",
		"card_number" : cardNumber,
	})
	if err != nil {
		return my_errors.GetError(my_errors.JsonMarshalError)
	}

	return string(rowResponse) ,nil
}


func getExtendedProducts(newProducts []types.Product, products []types.Product) (string, error) {
	for i:=0;i< len(newProducts);i++  {
		products = append(products, newProducts[i])
	}

	rowProducts, err := json.Marshal(products)
	if err != nil {
		return "", err
	}

	return string(rowProducts),nil
}


func parseDBShop(dbShop *types.Shop) (*jsonShop, error) {
	var jsonShop = &jsonShop{}
	err := jsonShop.fromDBShop(dbShop)
	if err != nil {
		return nil, err
	}
	return jsonShop, nil
}


func contains(sellers []int64, newSeller int64) bool {
	for i := 0;i< len(sellers); i++{
		if sellers[i] == newSeller {
			return true
		}
	}

	return false
}


func (jsonShop *jsonShop) getProduct(id int64) *types.Product {
	for _, val := range jsonShop.RowProducts {
		if val.Id == id {
			return &val
		}
	}

	return nil
}


func (jsonShop *jsonShop) fromDBShop(dbShop *types.Shop) error {
	jsonShop.Id = dbShop.Id
	jsonShop.Image = dbShop.Image
	jsonShop.Name = dbShop.Name
	jsonShop.CardNumber = dbShop.CardNumber
	jsonShop.Likes = dbShop.Likes
	jsonShop.Dislikes = dbShop.Dislikes
	jsonShop.Email = dbShop.Email
	jsonShop.Description = dbShop.Description
	jsonShop.RowSellers = dbShop.RowSellers
	return json.Unmarshal([]byte(dbShop.RowProducts), &jsonShop.RowProducts)
}