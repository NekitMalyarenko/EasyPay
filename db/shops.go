package db

import (
	"upper.io/db.v3"
	"types"
)

type Shops string


const (
	shopsTable = "shops"
)


func (*Shops) GetShopById(shopId int64)(*types.Shop ,error){
	var shop *types.Shop
	res := currentInstance.instance.Collection(shopsTable).Find(db.Cond{
		"id =":shopId,
	})
	err := res.One(&shop)
	if err != nil {
		return nil, err
	}
	return shop, nil
}


func (*Shops) GetShopCardNumber(shopId int64) (int64,error){
	var shop *types.Shop
	res := currentInstance.instance.Collection(shopsTable).Find(db.Cond{
		"id =":shopId,
	})

	err := res.One(&shop)
	if err != nil {
		return -1, err
	}

	return shop.CardNumber, nil
}


func (*Shops) AddShop(shop types.Shop) error{
	_, err := currentInstance.instance.
		InsertInto(shopsTable).Values(shop).Exec()
	return err
}


func (*Shops) UpdateProducts(products string,shopId int64) error{
	var shop *types.Shop
	res := currentInstance.instance.Collection(shopsTable).Find("id",shopId)

	err := res.One(&shop)
	if err != nil {
		return err
	}

	shop.RowProducts = products
	err = res.Update(shop)
	if err != nil {
		return err
	}

	return nil
}


func (*Shops) UpdateSellers(sellers string,shopId int64) error {
	var shop *types.Shop
	res := currentInstance.instance.Collection(shopsTable).Find("id",shopId)

	err := res.One(&shop)
	if err != nil {
		return err
	}

	shop.RowSellers = sellers
	err = res.Update(shop)
	if err != nil {
		return err
	}

	return nil
}


func (*Shops) AddDislikes(dislikesNumber int, shopId int64) error {
	var shop *types.Shop

	res := currentInstance.instance.Collection(shopsTable).Find(db.Cond{"id": shopId})
	err := res.One(&shop)
	if err != nil {
		return err
	}

	shop.Dislikes += dislikesNumber
	return res.Update(shop)
}


func (*Shops) AddLikes(likesNumber int, shopId int64) error {
	var shop *types.Shop

	res := currentInstance.instance.Collection(shopsTable).Find(db.Cond{"id": shopId})
	err := res.One(&shop)
	if err != nil {
		return err
	}

	shop.Likes += likesNumber
	return res.Update(shop)
}