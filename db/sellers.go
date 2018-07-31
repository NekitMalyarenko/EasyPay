package db

import (
	"types"
	"log"
	"upper.io/db.v3"
)

type Sellers string


const (
	sellersTable = "sellers"
)


func (*Sellers) GetSeller(phoneNumber string) (*types.Seller, error) {
	var user *types.Seller

	res := currentInstance.instance.Collection(sellersTable).Find(db.Cond{"phone_number": phoneNumber})
	has, err := res.Exists()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if has {
		err := res.One(&user)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return user, nil
	} else {
		return nil, nil
	}
}


func (*Sellers) GetSellerById(id int64) (*types.Seller, error) {
	var user *types.Seller

	res := currentInstance.instance.Collection(sellersTable).Find(db.Cond{"id": id})
	has, err := res.Exists()
	if err != nil {
		return nil, err
	}

	if has {
		err = res.One(&user)
		if err != nil {
			return nil, err
		}

		return user, nil
	} else {
		return nil, nil
	}
}


func (*Sellers) AddSeller(seller *types.Seller) error {
	_, err := currentInstance.instance.
		InsertInto(sellersTable).Values(seller).Exec()
	return err
}


func (db *Sellers) HasSeller(phoneNumber string) (bool, error) {
	user, err := db.GetSeller(phoneNumber)
	if err != nil {
		return false, err
	} else if user != nil {
		return true, nil
	} else {
		return false, nil
	}
}


func (*Sellers) UpdateShopId(sellerId int64, shopId int64) error {
	var seller *types.Seller
	res := currentInstance.instance.Collection(sellersTable).Find("id",sellerId)

	err := res.One(&seller)
	if err != nil {
		log.Println("cant find user with such id")
		return err
	}

	seller.ShopId = shopId
	err = res.Update(seller)
	if err != nil {
		return err
	}

	return nil
}
